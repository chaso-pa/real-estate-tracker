package models

import (
	"time"

	"github.com/chaso-pa/real-estate-tracker/internal/services"
	"github.com/chaso-pa/real-estate-tracker/internal/utils"
	"github.com/lucsky/cuid"
	"gorm.io/gorm/clause"
)

// EstateType represents the type of real estate property
type EstateType string

const (
	UsedHouse     EstateType = "used_house"
	NewHouse      EstateType = "new_house"
	Land          EstateType = "land"
	UsedApartment EstateType = "used_apartment"
	NewApartment  EstateType = "new_apartment"
)

// Estate represents a real estate property
type Estate struct {
	ID                 string      `json:"id" db:"id"`
	URL                string      `json:"url" db:"url"`
	Address            *string     `json:"address,omitempty" db:"address"`
	EstateType         *EstateType `json:"estate_type,omitempty" db:"estate_type"`
	Value              *int        `json:"value,omitempty" db:"value"`
	Railway            *string     `json:"railway,omitempty" db:"railway"`
	LandArea           *float64    `json:"land_area,omitempty" db:"land_area"`
	BuildingArea       *float64    `json:"building_area,omitempty" db:"building_area"`
	FloorPlan          *string     `json:"floor_plan,omitempty" db:"floor_plan"`
	YearOfConstruction *int        `json:"year_of_construction,omitempty" db:"year_of_construction"`
	FirstAppeared      *time.Time  `json:"first_appeared,omitempty" db:"first_appeared"`
	LastAppeared       *time.Time  `json:"last_appeared,omitempty" db:"last_appeared"`
	CreatedDate        *time.Time  `json:"created_date,omitempty" db:"created_date"`
	CreatedAt          time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at" db:"updated_at"`
}

type EstateResponse struct {
	Estates []*Estate
}

func EstatesSetValues(estates []*Estate) {
	curTime := time.Now()
	for _, estate := range estates {
		estate.ID = cuid.New()
		estate.FirstAppeared = &curTime
		estate.LastAppeared = &curTime
		estate.CreatedDate = &curTime
		estate.UpdatedAt = curTime
	}
}

func EstatesUpsert(estates []*Estate) {
	utils.GetDb().Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{`address`, `estate_type`, `value`, `railway`, `land_area`, `building_area`, `floor_plan`, `year_of_construction`, `last_appeared`, `updated_at`}),
	}).Create(estates)
}

func EstatesSchema() *services.JSONSchema {
	return &services.JSONSchema{
		Name:        "estate_extraction",
		Description: "Extract structured real estate information from text",
		Strict:      true,
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"estates": map[string]interface{}{
					"type":        "array",
					"description": "estate information",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"url": map[string]interface{}{
								"type":        "string",
								"description": "Property listing URL",
							},
							"address": map[string]interface{}{
								"type":        "string",
								"description": "Property address",
							},
							"estate_type": map[string]interface{}{
								"type":        "string",
								"enum":        []string{"used_house", "new_house", "land", "used_apartment", "new_apartment"},
								"description": "Type of real estate property",
							},
							"value": map[string]interface{}{
								"type":        "integer",
								"description": "Property price in yen",
							},
							"railway": map[string]interface{}{
								"type":        "string",
								"description": "Nearest railway station",
							},
							"land_area": map[string]interface{}{
								"type":        "number",
								"description": "Land area in square meters",
							},
							"building_area": map[string]interface{}{
								"type":        "number",
								"description": "Building area in square meters",
							},
							"floor_plan": map[string]interface{}{
								"type":        "string",
								"description": "Floor plan description (e.g., 3LDK, 2DK)",
							},
							"year_of_construction": map[string]interface{}{
								"type":        "integer",
								"description": "Year the property was constructed",
							},
						},
						"required":             []string{"url", "address", "estate_type", "building_area", "floor_plan", "land_area", "value", "railway", "year_of_construction"},
						"additionalProperties": false,
					},
				},
			},
			"required":             []string{"estates"},
			"additionalProperties": false,
		},
	}
}
