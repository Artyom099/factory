package model

import "time"

type PartGetRepoRequest struct {
	Uuid string
}
type PartGetRepoResponse struct {
	Part Part
}

type PartListRepoRequest struct {
	Filter *PartFilterRepo
}
type PartListRepoResponse struct {
	Parts []Part
}

type PartFilterRepo struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type PartCreateRepoRequest struct {
	Part
}
type PartCreateRepoResponse struct {
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
	// Metadata      map[string]*Value
	CreatedAt time.Time
	UpdatedAt time.Time
}
