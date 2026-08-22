package autopiter

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"parser/internal/config/autopiter"
	"parser/internal/domain"
	"parser/internal/errors"
	"parser/internal/utils"
	"time"

	"github.com/schollz/progressbar/v3"
)

const rateLimitPause = 10 * time.Second

type Service interface {
	ParseData()
}

type service struct {
	cfg   *autopiter.Config
	proxy []string
	parts []domain.Part
}

func NewService(cfg *autopiter.Config, proxy []string, parts []domain.Part) Service {
	return &service{
		cfg:   cfg,
		proxy: proxy,
		parts: parts,
	}
}

func (s *service) searchDetails(partNumber string, proxy string) (domain.AutopiterSearchDetailsResponse, errors.ServiceError) {
	var result domain.AutopiterSearchDetailsResponse

	link := domain.AutopiterSearchDetailsURL(partNumber)
	log.Printf("searchdetails request: partNumber=%q proxy=%q url=%s", partNumber, proxy, link)

	status, body, err := s.getWithRateLimit("searchdetails", link, proxy)
	if err != nil {
		log.Printf("searchdetails request error for %q: %v", partNumber, err)
		return domain.AutopiterSearchDetailsResponse{}, err
	}

	log.Printf("searchdetails response: partNumber=%q status=%d body=%s", partNumber, status, string(body))

	if unmarshalErr := json.Unmarshal(body, &result); unmarshalErr != nil {
		return domain.AutopiterSearchDetailsResponse{}, errors.UnableToUnmarshall(unmarshalErr)
	}

	log.Printf(
		"searchdetails parsed: partNumber=%q code=%q total=%d catalogs=%d",
		partNumber,
		result.Code,
		result.Data.Total,
		len(result.Data.Catalogs),
	)

	return result, nil
}

func (s *service) getCosts(proxy string, articleIDs ...int) (domain.AutopiterGetCostsResponse, errors.ServiceError) {
	var result domain.AutopiterGetCostsResponse
	if len(articleIDs) == 0 {
		return result, nil
	}

	link := domain.AutopiterGetCostsURL(articleIDs...)
	log.Printf("getcosts request: ids=%v proxy=%q url=%s", articleIDs, proxy, link)

	status, body, err := s.getWithRateLimit("getcosts", link, proxy)
	if err != nil {
		log.Printf("getcosts request error for ids=%v: %v", articleIDs, err)
		return domain.AutopiterGetCostsResponse{}, err
	}

	log.Printf("getcosts response: ids=%v status=%d body=%s", articleIDs, status, string(body))

	if unmarshalErr := json.Unmarshal(body, &result); unmarshalErr != nil {
		return domain.AutopiterGetCostsResponse{}, errors.UnableToUnmarshall(unmarshalErr)
	}

	log.Printf("getcosts parsed: ids=%v code=%q costs=%d", articleIDs, result.Code, len(result.Data))
	return result, nil
}

func (s *service) getWithRateLimit(requestName, link, proxy string) (int, []byte, errors.ServiceError) {
	for {
		ctx, cancel := s.requestContext()

		res, err := utils.MakeAutopiterGetRequest(ctx, link, proxy)
		cancel()
		if err != nil {
			return 0, nil, err
		}

		body, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return 0, nil, errors.BadRequest(readErr)
		}

		if res.StatusCode == http.StatusTooManyRequests {
			fmt.Fprintf(os.Stderr, "[%s] получен 429 Too Many Requests, пауза %s\n", requestName, rateLimitPause)
			time.Sleep(rateLimitPause)
			continue
		}

		return res.StatusCode, body, nil
	}
}

func (s *service) parsePartWithRetry(part domain.Part, proxy string) ([]domain.Model, errors.ServiceError) {
	details, err := s.searchDetails(part.PartNumber, proxy)
	if err != nil && errors2.Is(err.Error(), context.DeadlineExceeded) {
		details, err = s.searchDetails(part.PartNumber, proxy)
	}
	if err != nil {
		return nil, err
	}
	if len(details.Data.Catalogs) == 0 {
		return nil, nil
	}

	articleIDs := make([]int, 0, len(details.Data.Catalogs))
	for _, catalog := range details.Data.Catalogs {
		articleIDs = append(articleIDs, catalog.ID)
	}

	s.requestPause()

	costs, err := s.getCosts(proxy, articleIDs...)
	if err != nil && errors2.Is(err.Error(), context.DeadlineExceeded) {
		costs, err = s.getCosts(proxy, articleIDs...)
	}
	if err != nil {
		return nil, err
	}

	return domain.CatalogsToModels(part, details.Data.Catalogs, costs.Data), nil
}

func (s *service) requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(s.cfg.RequestTimeoutSec)*time.Second)
}

func (s *service) requestPause() {
	time.Sleep(utils.RandomizeMilliseconds(s.cfg.PauseMs))
}

func (s *service) ParseData() {
	logFile, err := os.OpenFile(s.cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicln(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if err = utils.WriteModelsToCSV(nil, s.cfg.OutputCSV, true); err != nil {
		log.Panicln("не удалось создать csv")
	}

	bar := progressbar.Default(int64(len(s.parts)))
	for _, detail := range s.parts {
		log.Printf("processing part: oem=%q partNumber=%q", detail.Oem, detail.PartNumber)

		var proxy string
		if proxyInterface := utils.ChooseRandom(s.proxy); proxyInterface != nil {
			proxy = proxyInterface.(string)
		}

		models, svcErr := s.parsePartWithRetry(detail, proxy)
		if svcErr != nil && errors2.Is(svcErr.Error(), context.DeadlineExceeded) {
			log.Println(svcErr)
			if writeErr := utils.WriteModelsToCSV([]domain.Model{
				{
					OriginalManufacturer: detail.Oem,
					OriginalPartNumber:   detail.PartNumber,
				},
			}, s.cfg.SkippedCSV, false); writeErr != nil {
				log.Print(writeErr)
			}
		} else if svcErr != nil {
			log.Print(svcErr)
			if writeErr := utils.WriteModelsToCSV([]domain.Model{
				{
					OriginalManufacturer: detail.Oem,
					OriginalPartNumber:   detail.PartNumber,
				},
			}, s.cfg.NotFoundCSV, false); writeErr != nil {
				log.Print(writeErr)
			}
		}

		if len(models) == 0 {
			log.Printf("нет результатов для %s (%s)", detail.PartNumber, detail.Oem)
			if writeErr := utils.WriteModelsToCSV([]domain.Model{
				{
					OriginalManufacturer: detail.Oem,
					OriginalPartNumber:   detail.PartNumber,
				},
			}, s.cfg.NotFoundCSV, false); writeErr != nil {
				log.Print(writeErr)
			}
		} else if writeErr := utils.WriteModelsToCSV(models, s.cfg.OutputCSV, false); writeErr != nil {
			log.Print(writeErr)
		}

		time.Sleep(utils.RandomizeMilliseconds(s.cfg.PauseMs))
		bar.Add(1)
	}
}
