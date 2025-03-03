package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"kawa_blog/database"
	"kawa_blog/models"
)

// 記事を取得する処理
func GetArticles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, content, image_url, created_at FROM articles")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var articles []models.Article
		for rows.Next() {
			var article models.Article
			err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.ImageURL, &article.CreatedAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			articles = append(articles, article)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)
	}
}

// 記事を作成するハンドラー
func CreateArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var article models.Article

		// JSON をデコード
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			http.Error(w, "無効なリクエスト", http.StatusBadRequest)
			return
		}

		// 記事をデータベースに保存
		if err := database.InsertArticle(db, &article); err != nil {
			http.Error(w, "データベースエラー", http.StatusInternalServerError)
			return
		}

		// レスポンス
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)
	}
}
