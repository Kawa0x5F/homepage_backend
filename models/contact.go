package models

import "time"

// 連絡先などの基本構造遺体
type Contact struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Link      string    `json:"link"`
	HasImage  bool      `json:"has_image"`
	CreatedAt time.Time `json:"created_at"`
}
