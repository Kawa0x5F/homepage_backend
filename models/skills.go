package models

import "time"

// 技術などの基本構造遺体
type Skills struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	HasImage  bool      `json:"has_image"`
	CreatedAt time.Time `json:"created_at"`
}
