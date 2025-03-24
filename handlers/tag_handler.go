package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"kawa_blog/database"
	"kawa_blog/models"
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
