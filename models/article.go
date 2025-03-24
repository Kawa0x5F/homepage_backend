package models

import "time"

// 記事の基本構造遺体
type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	ImageURL  *string   `json:"image_url,omitempty"`
	IsPublish bool      `json:"is_publish"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 記事一覧取得時用のタグあり構造体
type ArticleWithTag struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	ImageURL  *string   `json:"image_url,omitempty"`
	IsPublish bool      `json:"is_publish"`
	Tags      []Tag     `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 記事更新用のリクエストボディ用構造体
type UpdateArticleRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImageURL  string `json:"image_url"`
	IsPublish bool   `json:"is_publish"`
}
