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

// 新規記事のタグ情報を追加する関数
func InsertArticleTags(db *sql.DB, slug string, tags models.TagRequest) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var articleID int
	err = tx.QueryRow("SELECT id FROM articles WHERE slug = $1", slug).Scan(&articleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, tag := range tags.Tags {
		var tagID int
		query := `
		INSERT INTO tags (name) VALUES ($1)
		ON CONFLICT (name)
		DO UPDATE SET name = EXCLUDED.name
		RETURNING id
		`
		err = tx.QueryRow(query, tag).Scan(&tagID)
		if err != nil {
			tx.Rollback()
			return err
		}

		query = `
		INSERT INTO articles_tags (article_id, tag_id) VALUES ($1, $2)
		`
		_, err = tx.Exec(query, articleID, tagID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// 記事のタグ情報を更新する関数
func UpdateArticleTags(db *sql.DB, slug string, tags models.TagRequest) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var articleID int
	err = tx.QueryRow("SELECT id FROM articles WHERE slug = $1", slug).Scan(&articleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM articles_tags WHERE article_id = $1", articleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, tag := range tags.Tags {
		var tagID int
		query := `
		INSERT INTO tags (name) VALUES ($1)
		ON CONFLICT (name)
		DO UPDATE SET name = EXCLUDED.name
		RETURNING id
		`
		err = tx.QueryRow(query, tag).Scan(&tagID)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("INSERT INTO articles_tags (article_id, tag_id) VALUES ($1, $2)", articleID, tagID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
