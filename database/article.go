package database

import (
	"database/sql"
	"kawa_blog/models"
	"log"
)

// 記事をデータベースに追加
func InsertArticle(db *sql.DB, article *models.Article) error {
	query := `INSERT INTO articles (title, content, image_url) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	// クエリ実行
	err := db.QueryRow(query, article.Title, article.Content, article.ImageURL).Scan(&article.ID, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		log.Println("記事の挿入に失敗:", err)
		return err
	}
	return nil
}
