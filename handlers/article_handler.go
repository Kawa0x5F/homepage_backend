package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"kawa_blog/database"
	"kawa_blog/models"

	"github.com/gorilla/mux"
)

// 記事の一覧を取得する処理
func GetArticles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, content, image_url, created_at, updated_at FROM articles")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var articles []models.Article
		for rows.Next() {
			var article models.Article
			err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.ImageURL, &article.CreatedAt, &article.UpdatedAt)
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

// 記事を取得する処理
func GetArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ルートパラメータを取得
		vars := mux.Vars(r)
		idStr := vars["id"]

		// int に変換
		target_id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid article ID", http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT id, title, content, image_url, created_at, updated_at FROM articles")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var article models.Article
		for rows.Next() {
			var row models.Article
			err := rows.Scan(&row.ID, &row.Title, &row.Content, &row.ImageURL, &row.CreatedAt, &row.UpdatedAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if row.ID == target_id {
				article = row
				break
			}
		}

		w.Header().Set("Cointent-Type", "application/json")
		json.NewEncoder(w).Encode(article)
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
