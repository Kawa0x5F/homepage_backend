package models

import "time"

// プロフィールの基本構造遺体
type About struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Roma        string    `json:"roma"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	ImageURL    *string   `json:"image_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
