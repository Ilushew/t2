package models

type PlaceImage struct {
	ID        int    `json:"id"`
	PlaceID   int    `json:"place_id"`
	Filename  string `json:"filename"`
	SortOrder int    `json:"sort_order"`
}
