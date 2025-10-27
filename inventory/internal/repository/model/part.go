package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepoPartFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type RepoPart struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Uuid          string             `bson:"uuid"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description,omitempty"`
	Price         float64            `bson:"price"`
	StockQuantity int64              `bson:"stock_quantity"`
	Category      Category           `bson:"category"`
	Dimensions    *Dimensions        `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer      `bson:"manufacturer,omitempty"`
	Tags          []string           `bson:"tags,omitempty"`
	Metadata      map[string]*Value  `bson:"metadata,omitempty"`
	CreatedAt     *time.Time         `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at,omitempty"`
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
	Length float64 `bson:"length,omitempty"`
	Width  float64 `bson:"width,omitempty"`
	Height float64 `bson:"height,omitempty"`
	Weight float64 `bson:"weight,omitempty"`
}

type Manufacturer struct {
	Name    string `bson:"name,omitempty"`
	Country string `bson:"country,omitempty"`
	Website string `bson:"website,omitempty"`
}

type Value struct {
	StringValue *string  `bson:"string_value,omitempty"`
	Int64Value  *int64   `bson:"int64_value,omitempty"`
	DoubleValue *float64 `bson:"double_value,omitempty"`
	BoolValue   *bool    `bson:"bool_value,omitempty"`
}
