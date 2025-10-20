package model

import "time"

type PartGetServiceRequest struct {
	Uuid string
}
type PartGetServiceResponse struct {
	Part Part
}

type PartListServiceRequest struct {
	Filter *PartsFilterService
}
type PartListServiceResponse struct {
	Parts []Part
}

type PartsFilterService struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type PartCreateServiceRequest struct {
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	// Metadata      map[string]*Value
	CreatedAt time.Time
	UpdatedAt time.Time
}
type PartCreateServiceResponse struct {
	Uuid string
}

type Category int32

const (
	UNSPECIFIED Category = 0
	ENGINE      Category = 1
	FUEL        Category = 2
	PORTHOLE    Category = 3
	WING        Category = 4
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
