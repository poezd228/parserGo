package autopiller

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"io"
	"log"
	"os"
	"parser/internal/domain"
	"parser/internal/errors"
	"parser/internal/utils"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	outputCSV  = "autopiter.csv"
	skippedCSV = "autopiter_skipped.csv"
	notFoundCSV = "autopiter_notfound.csv"
)

type Service interface {
	ParseData()
}

type service struct {
	proxy []string
	parts []domain.Part
}

func NewService(proxy []string, parts []domain.Part) Service {
	return &service{
		proxy: proxy,
		parts: parts,
	}
}

func (s *service) search(partNumber string, proxy string) (domain.AutopiterResponse, errors.ServiceError) {
	var result domain.AutopiterResponse

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	requestBody := domain.NewAutopiterSearchRequest(partNumber)
	payload, _ := json.Marshal(requestBody)
	log.Printf("search request: partNumber=%q proxy=%q body=%s", partNumber, proxy, string(payload))

	res, err := utils.MakeAutopiterSearchRequest(ctx, domain.AutopiterSearchURL(), proxy, requestBody)
	if err != nil {
		log.Printf("search request error for %q: %v", partNumber, err)
		return domain.AutopiterResponse{}, err
	}
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Printf("search read body error for %q: %v", partNumber, readErr)
		return domain.AutopiterResponse{}, errors.BadRequest(readErr)
	}

	log.Printf("search response: partNumber=%q status=%d body=%s", partNumber, res.StatusCode, string(body))

	if unmarshalErr := json.Unmarshal(body, &result); unmarshalErr != nil {
		log.Printf("search unmarshal error for %q: %v", partNumber, unmarshalErr)
		return domain.AutopiterResponse{}, errors.UnableToUnmarshall(unmarshalErr)
	}

	log.Printf(
		"search parsed: partNumber=%q code=%q total=%d goods=%d",
		partNumber,
		result.Code,
		result.Data.Total,
		len(result.Data.Goods),
	)

	return result, nil
}

func (s *service) searchWithRetry(partNumber string, proxy string) (domain.AutopiterResponse, errors.ServiceError) {
	result, err := s.search(partNumber, proxy)
	if err != nil && errors2.Is(err.Error(), context.DeadlineExceeded) {
		return s.search(partNumber, proxy)
	}

	return result, err
}

func (s *service) ParseData() {
	logFile, err := os.OpenFile("autopiter.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicln(err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	if err = utils.WriteModelsToCSV(nil, outputCSV, true); err != nil {
		log.Panicln("не удалось создать csv")
	}

	bar := progressbar.Default(int64(len(s.parts)))
	for _, detail := range s.parts {
		log.Printf("processing part: oem=%q partNumber=%q", detail.Oem, detail.PartNumber)

		var proxy string
		if proxyInterface := utils.ChooseRandom(s.proxy); proxyInterface != nil {
			proxy = proxyInterface.(string)
		}

		result, svcErr := s.searchWithRetry(detail.PartNumber, proxy)
		if svcErr != nil && errors2.Is(svcErr.Error(), context.DeadlineExceeded) {
			log.Println(svcErr)
			if writeErr := utils.WriteModelsToCSV([]domain.Model{
				{
					OriginalManufacturer: detail.Oem,
					OriginalPartNumber:   detail.PartNumber,
				},
			}, skippedCSV, false); writeErr != nil {
				log.Print(writeErr)
			}
		} else if svcErr != nil {
			log.Print(svcErr)
			if writeErr := utils.WriteModelsToCSV([]domain.Model{
				{
					OriginalManufacturer: detail.Oem,
					OriginalPartNumber:   detail.PartNumber,
				},
			}, notFoundCSV, false); writeErr != nil {
				log.Print(writeErr)
			}
		}

		models := result.ToModel(detail)
		if len(models) == 0 {
			log.Printf("нет результатов для %s (%s)", detail.PartNumber, detail.Oem)
			if writeErr := utils.WriteModelsToCSV([]domain.Model{
				{
					OriginalManufacturer: detail.Oem,
					OriginalPartNumber:   detail.PartNumber,
				},
			}, notFoundCSV, false); writeErr != nil {
				log.Print(writeErr)
			}
		} else if writeErr := utils.WriteModelsToCSV(models, outputCSV, false); writeErr != nil {
			log.Print(writeErr)
		}

		time.Sleep(utils.RandomizeMilliseconds(200))
		bar.Add(1)
	}
}
