package database

import (
	"database/sql"
	"kawa_blog/models"
	"log"
)

// タグをデータベースに追加
func InsertTag(db *sql.DB, tag *models.Tag) error {
	query := `INSERT INTO tags (name) VALUES ($1) RETURNING id, created_at`

	// クエリ実行
	err := db.QueryRow(query, tag.Name).Scan(&tag.ID, &tag.CreatedAt)
	if err != nil {
		log.Println("タグの挿入に失敗:", err)
		return err
	}
	return nil
}
