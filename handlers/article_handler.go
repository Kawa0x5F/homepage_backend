package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"kawa_blog/database"
	"kawa_blog/models"

	"github.com/gorilla/mux"
)

// 記事の一覧を取得する処理
func GetArticles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, slug, content, image_url, created_at, updated_at FROM articles")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var articles []models.Article
		for rows.Next() {
			var article models.Article
			err := rows.Scan(&article.ID, &article.Title, &article.Slug, &article.Content, &article.ImageURL, &article.CreatedAt, &article.UpdatedAt)
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

// 特定の記事を取得する処理
func GetArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ルートパラメータを取得
		vars := mux.Vars(r)
		target_slug := vars["slug"]

		var article models.Article
		err := db.QueryRow("SELECT id, title, content, image_url, created_at, updated_at FROM articles WHERE slug = $1", target_slug).
			Scan(&article.ID, &article.Title, &article.Content, &article.ImageURL, &article.CreatedAt, &article.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Article not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		// JSON を返す
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)
	}
}

// 記事を作成する処理
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
