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
		ON CONFLICT (name) DO NOTHING
		RETURNING id
		`
		err = tx.QueryRow(query, tag).Scan(&tagID)

		if err == sql.ErrNoRows || tagID == 0 {
			query = `SELECT id FROM tags WHERE name = $1`
			err = tx.QueryRow(query, tag).Scan(&tagID)
			if err != nil {
				tx.Rollback()
				return err
			}
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

	return tx.Commit()
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

	// 現在のタグを取得
	rows, err := tx.Query("SELECT t.name FROM tags t JOIN articles_tags at ON t.id = at.tag_id WHERE at.article_id = $1", articleID)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	existingTags := make(map[string]bool)
	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			tx.Rollback()
			return err
		}
		existingTags[tagName] = true
	}

	// 新しく追加するタグ
	for _, tag := range tags.Tags {
		if existingTags[tag] {
			continue // 既に存在するタグはスキップ
		}

		var tagID int
		query := `
		INSERT INTO tags (name) VALUES ($1)
		ON CONFLICT (name) DO NOTHING
		RETURNING id
		`
		err = tx.QueryRow(query, tag).Scan(&tagID)

		if err == sql.ErrNoRows || tagID == 0 {
			query = `SELECT id FROM tags WHERE name = $1`
			err = tx.QueryRow(query, tag).Scan(&tagID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		_, err = tx.Exec("INSERT INTO articles_tags (article_id, tag_id) VALUES ($1, $2)", articleID, tagID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 削除するタグ
	for existingTag := range existingTags {
		found := false
		for _, tag := range tags.Tags {
			if existingTag == tag {
				found = true
				break
			}
		}
		if !found {
			_, err = tx.Exec(`
				DELETE FROM articles_tags 
				WHERE article_id = $1 AND tag_id = (SELECT id FROM tags WHERE name = $2)
			`, articleID, existingTag)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}
