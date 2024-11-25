package emex

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"os"
	"parser/internal/domain"
	"parser/internal/errors"
	"parser/internal/utils"
	"time"
)

type service struct {
	proxy      []string
	parts      []domain.Part
	locationId []string
}

type Service interface {
	ParseData()
}

func NewService(proxy []string, detailNumbers []domain.Part, locationId []string) Service {

	return &service{
		proxy:      proxy,
		parts:      detailNumbers,
		locationId: locationId,
	}

}

func (s *service) getVisitorId(proxy string) (domain.Config, errors.ServiceError) {
	var config domain.Config
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	link := "https://emex.ru/api/home/version"
	res, err := utils.MakeJsonRequestPost(ctx, link, proxy)
	if err != nil {
		return domain.Config{}, err
	}
	defer res.Body.Close()
	body, e := io.ReadAll(res.Body)
	if e != nil {
		return domain.Config{}, errors.BadRequest(e)
	}
	e = json.Unmarshal(body, &config)
	if e != nil {
		return domain.Config{}, errors.UnableToUnmarshall(e)
	}
	return config, nil

}

func (s *service) getCommonData(detailNum string, locationId string, proxy string) (domain.Result, errors.ServiceError) {
	var result domain.Result
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	link := fmt.Sprintf("https://emex.ru/api/search/search?detailNum=%s&locationId=%s", detailNum, locationId)

	res, err := utils.MakeJsonRequest(ctx, link, proxy)
	if err != nil {
		return domain.Result{}, err
	}
	defer res.Body.Close()
	body, e := io.ReadAll(res.Body)
	if e != nil {
		return domain.Result{}, errors.BadRequest(e)
	}
	e = json.Unmarshal(body, &result)
	if e != nil {
		return domain.Result{}, errors.UnableToUnmarshall(e)
	}
	return result, nil
}
func (s *service) getMainData(url string, proxy string, visitorId string) (domain.Result, errors.ServiceError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var result domain.Result
	urlModel, err := utils.ParseURL(url)
	if err != nil {

	}

	link := fmt.Sprintf("https://emex.ru/api/search/search?detailNum=%s&make=%s&locationId=%s", urlModel.DetailNum, urlModel.Make, urlModel.LocationId)
	res, e := utils.MakeJsonRequestVisitor(ctx, link, proxy, visitorId)
	if e != nil {
		return domain.Result{}, e

	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if e != nil {
		return domain.Result{}, errors.BadRequest(err)
	}
	er := json.Unmarshal(body, &result)
	if er != nil {
		return domain.Result{}, errors.UnableToUnmarshall(err)
	}
	return result, nil

}

func (s *service) parceTwice(detailNum string, locationId string, proxy string) (domain.Result, errors.ServiceError) {
	res, err := s.getCommonData(detailNum, locationId, proxy)
	if err != nil && errors2.Is(err.Error(), context.DeadlineExceeded) {
		res, err = s.getCommonData(detailNum, locationId, proxy)
		if err != nil {
			return domain.Result{}, err
		}

	} else if err != nil {
		return domain.Result{}, err

	}
	return res, nil

}

func (s *service) ParseData() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()
	log.SetOutput(logFile)
	header := utils.CheckHeader("emex.csv")
	err = utils.WriteModelsToCSV(nil, "emex.csv", header)
	if err != nil {
		log.Panicln("Ну не повезло ")
	}
	bar := progressbar.Default(int64(len(s.parts)))
	for _, detail := range s.parts {
		proxy := utils.ChooseRandom(s.proxy).(string)
		res, er := s.parceTwice(detail.PartNumber, utils.ChooseRandom(s.locationId).(string), proxy)
		if er != nil && errors2.Is(er.Error(), context.DeadlineExceeded) {
			_, _ = s.getVisitorId(proxy)
			log.Println(er)
			err = utils.WriteModelsToCSV([]domain.Model{
				{OriginalManufacturer: detail.Oem,
					OriginalPartNumber: detail.PartNumber},
			}, "skipped.csv", false)
			if err != nil {
				log.Print(err)
			}

		} else if er != nil {
			log.Print(er)
			err = utils.WriteModelsToCSV([]domain.Model{
				{OriginalManufacturer: detail.Oem,
					OriginalPartNumber: detail.PartNumber},
			}, "notfound.csv", false)
			if err != nil {
				log.Print(err)
			}

		}
		e := utils.WriteModelsToCSV(res.ToModel(detail), "emex.csv", false)
		if e != nil {
			log.Print(err)
		}

		time.Sleep(utils.RandomizeMilliseconds(200))

		bar.Add(1)

	}

}
