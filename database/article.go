package database

import (
	"database/sql"
	"errors"
	"fmt"
	"kawa_blog/models"
	"log"

	"github.com/lib/pq"
)

// 記事の一覧をデータベースから取得する
func GetAllArticles(db *sql.DB) ([]models.ArticleWithTag, error) {
	query := `
	SELECT id, title, slug, image_url, is_publish, created_at, updated_at
	FROM articles
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.ArticleWithTag
	articleMap := make(map[int]*models.ArticleWithTag)

	for rows.Next() {
		var article models.ArticleWithTag
		err := rows.Scan(&article.ID, &article.Title, &article.Slug, &article.ImageURL, &article.IsPublish, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
		articleMap[article.ID] = &articles[len(articles)-1]
	}

	tagQuery := `
	SELECT at.article_id, t.id, t.name
	FROM articles_tags at
	JOIN tags t ON at.tag_id = t.id
	WHERE at.article_id = ANY($1)
	`
	articleIDs := make([]int, len(articles))
	for i, article := range articles {
		articleIDs[i] = article.ID
	}

	tagRows, err := db.Query(tagQuery, pq.Array(articleIDs))
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var articleID int
		var tag models.Tag
		err := tagRows.Scan(&articleID, &tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}

		if article, exists := articleMap[articleID]; exists {
			article.Tags = append(article.Tags, tag)
		}
	}

	return articles, nil
}

// 公開されている記事の一覧をデータベースから取得する
func GetPublishArticles(db *sql.DB) ([]models.ArticleWithTag, error) {
	query := `
	SELECT id, title, slug, image_url, created_at, updated_at
	FROM articles
	WHERE is_publish = TRUE
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.ArticleWithTag
	articleMap := make(map[int]*models.ArticleWithTag)

	for rows.Next() {
		var article models.ArticleWithTag
		err := rows.Scan(&article.ID, &article.Title, &article.Slug, &article.ImageURL, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
		articleMap[article.ID] = &articles[len(articles)-1]
	}

	tagQuery := `
	SELECT at.article_id, t.id, t.name
	FROM articles_tags at
	JOIN tags t ON at.tag_id = t.id
	WHERE at.article_id = ANY($1)
	`
	articleIDs := make([]int, len(articles))
	for i, article := range articles {
		articleIDs[i] = article.ID
	}

	tagRows, err := db.Query(tagQuery, pq.Array(articleIDs))
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var articleID int
		var tag models.Tag
		err := tagRows.Scan(&articleID, &tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}

		if article, exists := articleMap[articleID]; exists {
			article.Tags = append(article.Tags, tag)
		}
	}

	return articles, nil
}

// 特定の記事をデータベースから取得する
func GetArticleBySlug(db *sql.DB, slug string) (*models.ArticleWithTag, error) {
	var article models.ArticleWithTag
	query := `
	SELECT id, title, content, image_url, is_publish, created_at, updated_at
	FROM articles
	WHERE slug = $1
	`
	err := db.QueryRow(query, slug).
		Scan(&article.ID, &article.Title, &article.Content, &article.ImageURL, &article.IsPublish, &article.CreatedAt, &article.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("記事が見つかりません:", slug)
			return nil, nil
		}
		log.Println("記事の取得に失敗しました:", err)
		return nil, err
	}

	tagQuery := `
	SELECT t.id, t.name
	FROM tags t
	JOIN articles_tags at ON t.id = at.tag_id
	WHERE at.article_id = $1
	`
	rows, err := db.Query(tagQuery, article.ID)
	if err != nil {
		log.Println("タグの取得に失敗しました:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			log.Println("タグデータのスキャンに失敗しました:", err)
			return nil, err
		}
		article.Tags = append(article.Tags, tag)
	}

	if err := rows.Err(); err != nil {
		log.Println("タグの取得中にエラーが発生しました:", err)
		return nil, err
	}

	return &article, nil
}

// 記事をデータベースに追加
func InsertArticle(db *sql.DB, article *models.Article) error {
	query := `INSERT INTO articles (title, slug, content, image_url, is_publish) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	// クエリ実行
	err := db.QueryRow(query, article.Title, article.Slug, article.Content, article.ImageURL, article.IsPublish).Scan(&article.ID, &article.CreatedAt, &article.UpdatedAt)
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

// 記事を更新する
func PatchArticle(db *sql.DB, slug string, article models.UpdateArticleRequest) (bool, error) {
	query := `
	UPDATE articles
	SET title = $1, content = $2, image_url = $3, is_publish = $4, updated_at = NOW()
	WHERE slug = $5
	`
	result, err := db.Exec(query, article.Title, article.Content, article.ImageURL, article.IsPublish, slug)

	if err != nil {
		log.Println("記事の更新に失敗しました:", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("更新結果の取得に失敗しました:", err)
		return false, err
	}

	if rowsAffected == 0 {
		log.Println("更新対象の記事が見つかりません:", slug)
		return false, nil
	}

	log.Println("記事が更新されました:", slug)
	return true, nil
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
