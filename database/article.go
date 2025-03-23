package database

import (
	"database/sql"
	"errors"
	"fmt"
	"kawa_blog/models"
	"log"
)

// 記事の一覧をデータベースから取得する
func GetAllArticles(db *sql.DB) ([]models.Article, error) {
	rows, err := db.Query("SELECT id, title, slug, image_url, updated_at FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Slug, &article.ImageURL, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// 公開されている記事の一覧をデータベースから取得する
func GetPublicArticles(db *sql.DB) ([]models.Article, error) {
	rows, err := db.Query("SELECT id, title, slug, image_url, updated_at FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Slug, &article.ImageURL, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// 特定の記事をデータベースから取得する
func GetArticleBySlug(db *sql.DB, slug string) (*models.Article, error) {
	var article models.Article
	err := db.QueryRow("SELECT id, title, content, image_url, created_at, updated_at FROM articles WHERE slug = $1", slug).
		Scan(&article.ID, &article.Title, &article.Content, &article.ImageURL, &article.CreatedAt, &article.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

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

// 記事をデータベースから削除
func DeleteArticleBySlug(db *sql.DB, slug string) (bool, error) {
	result, err := db.Exec("DELETE FROM articles WHERE slug = $1", slug)
	if err != nil {
		return false, err
	}

	// 削除された行数を確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, err
	}

	return true, nil
}
