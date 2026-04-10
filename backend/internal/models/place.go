package models

type Place struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	Price           string       `json:"price"`
	Time            float64      `json:"time"`
	TypesOfMovement string       `json:"types_of_movement"`
	Category        string       `json:"category"`
	LatStart        float64      `json:"lat_start"`
	LonStart        float64      `json:"lon_start"`
	LatEnd          float64      `json:"lat_end"`
	LonEnd          float64      `json:"lon_end"`
	IsIndoor        bool         `json:"is_indoor"`
	WithChild       bool         `json:"with_child"`
	WithPets        bool         `json:"with_pets"`
	Description     string       `json:"description"`
	Images          []PlaceImage `json:"images"`
}
