package domain

type Model struct {
	OriginalManufacturer string
	OriginalPartNumber   string
	PartManufacturer     string
	PartNumber           string
	PartDescription      string
	Price                string
	DeliveryTime         string
	ParsedAt             string
}
type Part struct {
	PartNumber string
	Oem        string
}
