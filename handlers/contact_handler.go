package handlers

import (
	"database/sql"
	"encoding/json"
	"kawa_blog/database"
	"kawa_blog/models"
	"log"
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

// コンタクトを取得
func GetContactByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なID")
			return
		}

		contact, err := database.GetContactByID(db, id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		if contact == nil {
			respondWithError(w, http.StatusNotFound, "About not found")
			return
		}
		respondWithJSON(w, http.StatusOK, contact)
	}
}

// コンタクトを追加するハンドラー
func CreateContact(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contact models.Contact
		if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト")
			return
		}

		if err := database.InsertContact(db, &contact); err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		respondWithJSON(w, http.StatusCreated, contact)
	}
}

// コンタクトを更新
func PatchContact(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なID")
			return
		}

		var req models.Contact
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("JSONデコードエラー:", err)
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト: JSONの形式が正しくありません")
			return
		}

		// 名前とリンクが空でないかチェック
		if req.Name == "" || req.Link == "" {
			respondWithError(w, http.StatusBadRequest, "名前もしくはリンクが空です")
			return
		}

		// データベースを更新
		updated, err := database.PatchContact(db, id, req)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		if !updated {
			respondWithError(w, http.StatusNotFound, "About not found")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "About update"})
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
