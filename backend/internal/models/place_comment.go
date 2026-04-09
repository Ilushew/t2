package models

import "time"

type PlaceComment struct {
	ID        int       `json:"id"`
	PlaceID   int       `json:"place_id"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
