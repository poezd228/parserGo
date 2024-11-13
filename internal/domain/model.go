package domain

type Model struct {
	OriginalManufacturer string
	OriginalPartNumber   string
	PartManufacturer     string
	PartNumber           string
	PartDescription      string
	Price                string
	DeliveryTime         string
}
type Part struct {
	Oem        string
	PartNumber string
}
