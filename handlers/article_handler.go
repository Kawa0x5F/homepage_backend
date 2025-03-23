package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"kawa_blog/database"
	"kawa_blog/models"

	"github.com/gorilla/mux"
)

// 記事を作成する
func CreateArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var article models.Article
		if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト")
			return
		}

		if err := database.InsertArticle(db, &article); err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		respondWithJSON(w, http.StatusCreated, article)
	}
}

// 記事一覧を取得
func GetArticles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := database.GetAllArticles(db)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		respondWithJSON(w, http.StatusOK, articles)
	}
}

// 公開記事一覧を取得
func GetPublicArticles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := database.GetPublicArticles(db)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		respondWithJSON(w, http.StatusOK, articles)
	}
}

// 特定の記事を取得
func GetArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]

		article, err := database.GetArticleBySlug(db, slug)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		if article == nil {
			respondWithError(w, http.StatusNotFound, "Article not found")
			return
		}
		respondWithJSON(w, http.StatusOK, article)
	}
}

// 記事を削除する処理
func DeleteArticle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]

		deleted, err := database.DeleteArticleBySlug(db, slug)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		if !deleted {
			respondWithError(w, http.StatusNotFound, "Article not found")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Article deleted"})
	}
}

// 共通のJSONレスポンス関数
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// 共通のエラーレスポンス関数
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
