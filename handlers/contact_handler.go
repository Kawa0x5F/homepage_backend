package handlers

import (
	"database/sql"
	"encoding/json"
	"kawa_blog/database"
	"kawa_blog/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// 全コンタクトデータを取得するハンドラー
func GetAllContact(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contacts, err := database.GetAllContact(db)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		respondWithJSON(w, http.StatusOK, contacts)
	}
}

// コンタクトを追加するハンドラー
func CreateContact(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contact models.Contact
		if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			respondWithError(w, http.StatusBadRequest, "無効内リクエスト")
			return
		}

		if err := database.InsertContact(db, &contact); err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		respondWithJSON(w, http.StatusCreated, contact)
	}
}

// コンタクトを削除するハンドラー
func DeleteContactByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なID")
			return
		}

		deleted, err := database.DeleteContactByID(db, id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		if !deleted {
			respondWithError(w, http.StatusNotFound, "指定されたコンタクトがありませんでした")
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Contact deleted"})
	}
}
