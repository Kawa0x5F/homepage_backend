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

// 全てのタグをデータベースから取得
func GetAllTags(db *sql.DB) ([]models.Tag, error) {
	query := `SELECT id, name FROM tags`
	rows, err := db.Query(query)
	if err != nil {
		log.Println("全てのタグの取得に失敗:", err)
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			log.Println("タグの取得に失敗:", err)
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
