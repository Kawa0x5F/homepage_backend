package models

import "time"

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"tag_name"`
	CreatedAt time.Time `json:"date"`
}
