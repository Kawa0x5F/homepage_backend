package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"kawa_blog/database"
	"kawa_blog/models"
)

// タグを取得する処理
func GetTags(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM tags")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tags []models.Tag
		for rows.Next() {
			var tag models.Tag
			err := rows.Scan(&tag.ID, &tag.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tags = append(tags, tag)
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
			http.Error(w, "無効なリクエスト", http.StatusBadRequest)
			return
		}

		// タグをデータベースに保存
		if err := database.InsertTag(db, &tag); err != nil {
			http.Error(w, "データベースエラー", http.StatusInternalServerError)
			return
		}

		// レスポンス
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tag)
	}
}
