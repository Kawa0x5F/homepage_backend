package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"kawa_blog/database"
	"kawa_blog/models"

	"github.com/gorilla/mux"
)

// タグを取得する処理
func GetALLTags(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := database.GetAllTags(db)
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "データベースエラー"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tags)
	}
}

// タグを作成するハンドラー
func CreateTag(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tag models.Tag

		// JSON をデコード
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
			return
		}

		// タグをデータベースに保存
		if err := database.InsertTag(db, &tag); err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "データベースエラー"})
			return
		}

		// レスポンス
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tag)
	}
}

// 記事と紐づいたタグを追加するハンドラー
func InsertArticleTags(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type Req struct {
			Slug string            `json:"slug"`
			Tags models.TagRequest `json:"tags"`
		}

		var req Req

		// JSON をデコード
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
			return
		}

		slug := req.Slug
		tags := req.Tags

		// タグをデータベースに保存
		if err := database.InsertArticleTags(db, slug, tags); err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "データベースエラー"})
			return
		}

		// レスポンス
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Article update"})
	}
}

// 記事と紐づいたタグを追加するハンドラー
func UpdateArticleTags(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]
		if slug == "" {
			respondWithError(w, http.StatusBadRequest, "Slugが指定されていません")
			return
		}

		var tagReq models.TagRequest

		// JSON をデコード
		if err := json.NewDecoder(r.Body).Decode(&tagReq); err != nil {
			writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
			return
		}

		// タグをデータベースに保存
		if err := database.UpdateArticleTags(db, slug, tagReq); err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "データベースエラー"})
			return
		}

		// レスポンス
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Article update"})
	}
}
