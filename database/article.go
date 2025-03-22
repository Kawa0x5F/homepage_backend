package database

import (
	"database/sql"
	"fmt"
	"kawa_blog/models"
	"log"
)

// 記事をデータベースに追加
func InsertArticle(db *sql.DB, article *models.Article) error {
	query := `INSERT INTO articles (title, slug, content, image_url) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	// クエリ実行
	err := db.QueryRow(query, article.Title, article.Slug, article.Content, article.ImageURL).Scan(&article.ID, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		// `slug` の重複エラーかどうかを判定
		if err.Error() == `pq: duplicate key value violates unique constraint "articles_slug_key"` {
			log.Println("エラー: 同じ slug を持つ記事がすでに存在します:", article.Slug)
			return fmt.Errorf("slug '%s' はすでに存在します", article.Slug)
		}
		log.Println("記事の挿入に失敗:", err)
		return err
	}
	return nil
}
