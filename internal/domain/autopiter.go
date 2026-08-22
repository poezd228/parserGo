package domain

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	autopiterSearchDetailsURL = "https://autopiter.ru/api/api/searchdetails"
	autopiterGetCostsURL      = "https://autopiter.ru/api/api/appraise/getcosts"
)

type AutopiterMeta struct {
	FrontendType int    `json:"frontendType"`
	RenderType   int    `json:"renderType"`
	RouteID      string `json:"routeId"`
}

type AutopiterSearchDetailsResponse struct {
	Data AutopiterSearchDetailsData `json:"data"`
	Code string                     `json:"code"`
}

type AutopiterSearchDetailsData struct {
	Catalogs     []AutopiterCatalog `json:"catalogs"`
	Total        int                `json:"total"`
	IsFullSearch bool               `json:"isFullSearch"`
}

type AutopiterCatalog struct {
	ID          int     `json:"id"`
	CatalogID   int     `json:"catalogId"`
	Number      string  `json:"number"`
	Name        *string `json:"name"`
	CatalogName string  `json:"catalogName"`
	CatalogURL  string  `json:"catalogUrl"`
}

type AutopiterGetCostsResponse struct {
	Data []AutopiterCost `json:"data"`
	Code string          `json:"code"`
}

type AutopiterCost struct {
	ID                       int     `json:"id"`
	AnalogPrice              float64 `json:"analogPrice"`
	AnalogPriceNet           float64 `json:"analogPriceNet"`
	OriginalPrice            float64 `json:"originalPrice"`
	OriginalPriceNet         float64 `json:"originalPriceNet"`
	DeliveryDays             *int    `json:"deliveryDays"`
	IsOurStore               bool    `json:"isOurStore"`
	IsExpressDelivery        bool    `json:"isExpressDelivery"`
	ExpressDeliveryHoursText *string `json:"expressDeliveryHoursText"`
}

func DefaultAutopiterMeta() AutopiterMeta {
	return AutopiterMeta{
		FrontendType: 2,
		RenderType:   1,
		RouteID:      "MAIN_PAGE",
	}
}

func AutopiterSearchDetailsURL(detailNumber string) string {
	return buildAutopiterURL(autopiterSearchDetailsURL, map[string]string{
		"detailNumber": detailNumber,
	})
}

func AutopiterGetCostsURL(articleIDs ...int) string {
	values := url.Values{}
	meta := DefaultAutopiterMeta()
	values.Set("meta[frontendType]", fmt.Sprintf("%d", meta.FrontendType))
	values.Set("meta[renderType]", fmt.Sprintf("%d", meta.RenderType))
	values.Set("meta[routeId]", meta.RouteID)
	for _, id := range articleIDs {
		values.Add("idArticles", fmt.Sprintf("%d", id))
	}

	return autopiterGetCostsURL + "?" + values.Encode()
}

func buildAutopiterURL(base string, params map[string]string) string {
	values := url.Values{}
	meta := DefaultAutopiterMeta()
	values.Set("meta[frontendType]", fmt.Sprintf("%d", meta.FrontendType))
	values.Set("meta[renderType]", fmt.Sprintf("%d", meta.RenderType))
	values.Set("meta[routeId]", meta.RouteID)
	for key, value := range params {
		values.Set(key, value)
	}

	return base + "?" + values.Encode()
}

func CatalogsToModels(part Part, catalogs []AutopiterCatalog, costs []AutopiterCost) []Model {
	if len(catalogs) == 0 {
		return nil
	}

	costByID := make(map[int]AutopiterCost, len(costs))
	for _, cost := range costs {
		costByID[cost.ID] = cost
	}

	parsedAt := time.Now().Format("02.01.2006 15:04")
	models := make([]Model, 0, len(catalogs))

	for _, catalog := range catalogs {
		cost, ok := costByID[catalog.ID]
		if !ok {
			continue
		}

		partNumber := strings.TrimSpace(catalog.Number)
		description := catalog.description()
		deliveryTime := formatDeliveryDays(cost.DeliveryDays)

		if cost.AnalogPrice > 0 {
			models = append(models, Model{
				OriginalManufacturer: part.Oem,
				OriginalPartNumber:   part.PartNumber,
				PartManufacturer:     catalog.CatalogName,
				PartNumber:           partNumber,
				PartDescription:      description,
				Price:                fmt.Sprintf("%.2f", cost.AnalogPrice),
				DeliveryTime:         deliveryTime,
				ParsedAt:             parsedAt,
			})
		}
		if cost.OriginalPrice > 0 {
			models = append(models, Model{
				OriginalManufacturer: part.Oem,
				OriginalPartNumber:   part.PartNumber,
				PartManufacturer:     catalog.CatalogName,
				PartNumber:           partNumber,
				PartDescription:      description,
				Price:                fmt.Sprintf("%.2f", cost.OriginalPrice),
				DeliveryTime:         deliveryTime,
				ParsedAt:             parsedAt,
			})
		}
	}

	return models
}

func (c AutopiterCatalog) description() string {
	if c.Name == nil {
		return ""
	}
	return strings.TrimSpace(*c.Name)
}

func formatDeliveryDays(days *int) string {
	if days == nil {
		return ""
	}
	return fmt.Sprintf("%d дн.", *days)
}
