package models

import "time"

type ArticleTags struct {
	ArticleID int       `json:"article_id"`
	TagID     int       `json:"tag_id"`
	CreatedAt time.Time `json:"date"`
}
