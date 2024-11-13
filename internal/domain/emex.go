package domain

import "fmt"

type Result struct {
	SearchResult SearchResult `json:"searchResult"`
	ErrorMessage string       `json:"errorMessage"`
}

type SearchResult struct {
	Num           string    `json:"num"`
	IsPromoDetail bool      `json:"isPromoDetail"`
	Name          string    `json:"name"`
	Make          string    `json:"make"`
	MakeID        int       `json:"makeId"`
	Description   string    `json:"description"`
	Keywords      string    `json:"keywords"`
	Makes         Makes     `json:"makes"`
	Filters       []Filter  `json:"filters"`
	Location      Location  `json:"location"`
	Originals     []Product `json:"originals"`
	Analogs       []Product `json:"analogs"`
	Replacements  []Product `json:"replacements"`
	NoResults     bool      `json:"noResults"`
}

// Структура для Makes
type Makes struct {
	Header string     `json:"header"`
	List   []MakeItem `json:"list"`
}

type MakeItem struct {
	ID            int       `json:"id"`
	Make          string    `json:"make"`
	Num           string    `json:"num"`
	IsPromoDetail bool      `json:"isPromoDetail"`
	Name          string    `json:"name"`
	Tags          []string  `json:"tags"`
	URL           string    `json:"url"`
	BestPrice     BestPrice `json:"bestPrice"`
}

type BestPrice struct {
	Value      float64 `json:"value"`
	Symbol     string  `json:"symbol"`
	SymbolText string  `json:"symbolText"`
}

// Структура для Filters
type Filter struct {
	ID    string       `json:"id"`
	Type  string       `json:"type"`
	Title string       `json:"title"`
	Range *Range       `json:"range,omitempty"`
	Value interface{}  `json:"value"`
	Units string       `json:"units,omitempty"`
	Items []FilterItem `json:"items,omitempty"`
}

type Range struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type FilterItem struct {
	ID      string    `json:"id"`
	Enabled bool      `json:"enabled"`
	Title   string    `json:"title"`
	Value   BestPrice `json:"value"`
}

// Структура для Location
type Location struct {
	LocationID              int       `json:"locationId"`
	Name                    string    `json:"name"`
	City                    string    `json:"city"`
	Longitude               float64   `json:"longitude"`
	Latitude                float64   `json:"latitude"`
	Schedule                string    `json:"schedule"`
	Address                 string    `json:"address"`
	Description             string    `json:"description"`
	FreeDelivery            BestPrice `json:"freeDelivery"`
	Contact                 []string  `json:"contact"`
	Email                   string    `json:"email"`
	Phones                  []string  `json:"phones"`
	ImageUrls               []string  `json:"imageUrls"`
	URL                     string    `json:"url"`
	MultiBasketSupport      bool      `json:"multiBasketSupport"`
	SortingMode             string    `json:"sortingMode"`
	IsShop                  bool      `json:"isShop"`
	OptovikID               int       `json:"optovikId"`
	Tags                    []string  `json:"tags"`
	IsPostPayEnabled        bool      `json:"isPostPayEnabled"`
	IsOnlinePaymentsEnabled bool      `json:"isOnlinePaymentsEnabled"`
	CityLocationCount       int       `json:"cityLocationCount"`
	ShelfLifeDays           int       `json:"shelfLifeDays"`
	ShelfLifeDaysType       int       `json:"shelfLifeDaysType"`
	IsPublished             bool      `json:"isPublished"`
}

// Структура для продуктов, которые встречаются в Originals, Analogs и Replacements
type Product struct {
	DetailKey          string    `json:"detailKey"`
	ProductID          string    `json:"productId"`
	DetailNum          string    `json:"detailNum"`
	IsPromoDetail      bool      `json:"isPromoDetail"`
	FormattedDetailNum string    `json:"formattedDetailNum"`
	Make               string    `json:"make"`
	Name               string    `json:"name"`
	HasPhoto           bool      `json:"hasPhoto"`
	HasApplicability   bool      `json:"hasApplicability"`
	Offers             []Offer   `json:"offers"`
	Count              int       `json:"count"`
	MinDelivery        Delivery  `json:"minDelivery"`
	MinPrice           BestPrice `json:"minPrice"`
	OtherCount         int       `json:"otherCount"`
	OtherMinDelivery   Delivery  `json:"otherMinDelivery"`
	OtherMinPrice      BestPrice `json:"otherMinPrice"`
}

// Структура для Offers
type Offer struct {
	OfferKey     string    `json:"offerKey"`
	Quantity     int       `json:"quantity"`
	LotQuantity  int       `json:"lotQuantity"`
	Delivery     Delivery  `json:"delivery"`
	Price        BestPrice `json:"price"`
	DisplayPrice BestPrice `json:"displayPrice"`
	Data         OfferData `json:"data"`
	FilterID     int       `json:"filterId"`
	IsSeller     bool      `json:"isSeller"`
	IsCustom     bool      `json:"isCustom"`
	Rating       float64   `json:"rating"`
}

type Delivery struct {
	Value int    `json:"value"`
	Units string `json:"units"`
}

type OfferData struct {
	OfferKey         string    `json:"offerKey"`
	MakeName         string    `json:"makeName"`
	DetailNum        string    `json:"detailNum"`
	IsPromoDetail    bool      `json:"isPromoDetail"`
	DetailName       string    `json:"detailName"`
	DetailNumToShow  string    `json:"detailNumToShow"`
	Quantity         Quantity  `json:"quantity"`
	LotQuantity      Quantity  `json:"lotQuantity"`
	MaxQuantity      Quantity  `json:"maxQuantity"`
	Price            BestPrice `json:"price"`
	DisplayPrice     BestPrice `json:"displayPrice"`
	Delivery         Delivery  `json:"delivery"`
	NotAvailable     bool      `json:"notAvailable"`
	OldNotAvailable  bool      `json:"oldNotAvailable"`
	HasApplicability bool      `json:"hasApplicability"`
	HasImages        bool      `json:"hasImages"`
	LocalID          int       `json:"localId"`
	LocationID       int       `json:"locationId"`
	LocationName     string    `json:"locationName"`
}

type Quantity struct {
	Value int    `json:"value"`
	Units string `json:"units"`
}
type UrlModel struct {
	DetailNum  string
	Make       string
	LocationId string
}

type Settings struct {
	GeoIpEnabled              bool `json:"geoIpEnabled"`
	UseNewSuggestionsEndpoint bool `json:"useNewSuggestionsEndpoint"`
}

type ExperimentGroup struct {
	Group string `json:"group"`
}

type Experiments struct {
	VisitorId            string          `json:"visitorId"`
	MainPageRetainUsers2 ExperimentGroup `json:"mainPageRetainUsers2"`
	BasketExpertCheck    ExperimentGroup `json:"basketExpertCheck"`
	UseNewCheckoutPage   ExperimentGroup `json:"useNewCheckoutPage"`
	ShowExpertWidget     ExperimentGroup `json:"showExpertWidget"`
	ShowSintecPromo      ExperimentGroup `json:"showSintecPromo"`
}

type PosthogExperiments struct {
	NewCategoriesCatalogAutoProducts1 string `json:"newCategoriesCatalogAutoProducts1"`
}

type IpLocation struct {
	CityName    string  `json:"cityName"`
	RegionName  string  `json:"regionName"`
	CountryCode string  `json:"countryCode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type Config struct {
	Version             string             `json:"version"`
	CanSwitch           bool               `json:"canSwitch"`
	MultiBasketAllowed  bool               `json:"multiBasketAllowed"`
	Inmotion2Enabled    bool               `json:"inmotion2Enabled"`
	Settings            Settings           `json:"settings"`
	Experiments         Experiments        `json:"experiments"`
	PosthogExperiments  PosthogExperiments `json:"posthogExperiments"`
	IpLocation          IpLocation         `json:"ipLocation"`
	Maintenance2Enabled bool               `json:"maintenance2Enabled"`
}

func (r *Result) ToModel(part Part) []Model {
	var models []Model
	for _, offer := range r.SearchResult.Originals {
		for _, data := range offer.Offers {
			model := Model{
				OriginalManufacturer: part.Oem,
				OriginalPartNumber:   part.PartNumber,
				PartManufacturer:     data.Data.MakeName,
				PartNumber:           data.Data.DetailNum,
				PartDescription:      data.Data.DetailName,
				Price:                fmt.Sprintf("%.2f %s", data.Price.Value, data.Price.Symbol),
				DeliveryTime:         fmt.Sprintf("%d %s", data.Delivery.Value, data.Delivery.Units),
			}
			models = append(models, model)

		}
	}
	for _, offer := range r.SearchResult.Analogs {
		for _, data := range offer.Offers {
			model := Model{
				OriginalManufacturer: part.Oem,
				OriginalPartNumber:   part.PartNumber,
				PartManufacturer:     data.Data.MakeName,
				PartNumber:           data.Data.DetailNum,
				PartDescription:      data.Data.DetailName,
				Price:                fmt.Sprintf("%.2f %s", data.Price.Value, data.Price.Symbol),
				DeliveryTime:         fmt.Sprintf("%d %s", data.Delivery.Value, data.Delivery.Units),
			}
			models = append(models, model)

		}
	}
	for _, offer := range r.SearchResult.Replacements {
		for _, data := range offer.Offers {
			model := Model{
				OriginalManufacturer: part.Oem,
				OriginalPartNumber:   part.PartNumber,
				PartManufacturer:     data.Data.MakeName,
				PartNumber:           data.Data.DetailNum,
				PartDescription:      data.Data.DetailName,
				Price:                fmt.Sprintf("%.2f", data.Price.Value),
				DeliveryTime:         fmt.Sprintf("%d %s", data.Delivery.Value, data.Delivery.Units),
			}
			models = append(models, model)

		}
	}
	return models

}
