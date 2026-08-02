package domain

import (
	"fmt"
	"time"
)

const autopiterSearchURL = "https://autopiter.ru/api/api/showcase/elk/search-by-name"

type AutopiterSearchRequest struct {
	Meta        AutopiterSearchMeta `json:"meta"`
	SearchValue string              `json:"searchValue"`
	Top         int                 `json:"top"`
}

type AutopiterSearchMeta struct {
	FrontendType int    `json:"frontendType"`
	RenderType   int    `json:"renderType"`
	RouteID      string `json:"routeId"`
}

type AutopiterResponse struct {
	Data AutopiterData `json:"data"`
	Code string        `json:"code"`
}

type AutopiterData struct {
	Goods []AutopiterGood `json:"goods"`
	Total int             `json:"total"`
}

type AutopiterGood struct {
	CatalogName  string              `json:"catalogName"`
	ShortName    string              `json:"shortName"`
	Name         string              `json:"name"`
	Price        float64             `json:"price"`
	DeliveryDays int                 `json:"deliveryDays"`
	Properties   []AutopiterProperty `json:"properties"`
}

type AutopiterProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewAutopiterSearchRequest(partNumber string, top int) AutopiterSearchRequest {
	if top <= 0 {
		top = 12
	}
	return AutopiterSearchRequest{
		Meta: AutopiterSearchMeta{
			FrontendType: 2,
			RenderType:   1,
			RouteID:      "MAIN_PAGE",
		},
		SearchValue: partNumber,
		Top:         top,
	}
}

func AutopiterSearchURL() string {
	return autopiterSearchURL
}

func (r *AutopiterResponse) ToModel(part Part) []Model {
	if r == nil || len(r.Data.Goods) == 0 {
		return nil
	}

	parsedAt := time.Now().Format("02.01.2006 15:04")
	models := make([]Model, 0, len(r.Data.Goods))

	for _, good := range r.Data.Goods {
		models = append(models, Model{
			OriginalManufacturer: part.Oem,
			OriginalPartNumber:   part.PartNumber,
			PartManufacturer:     good.manufacturer(),
			PartNumber:           good.ShortName,
			PartDescription:      good.Name,
			Price:                fmt.Sprintf("%.2f", good.Price),
			DeliveryTime:         fmt.Sprintf("%d дн.", good.DeliveryDays),
			ParsedAt:             parsedAt,
		})
	}

	return models
}

func (g AutopiterGood) manufacturer() string {
	for _, property := range g.Properties {
		if property.Name == "Бренд" {
			return property.Value
		}
	}

	if g.CatalogName != "" {
		return g.CatalogName
	}

	return ""
}
