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

// プロフィールを作成する
func CreateAbout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var about models.About
		if err := json.NewDecoder(r.Body).Decode(&about); err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト")
			return
		}

		if err := database.InsertAbout(db, &about); err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		respondWithJSON(w, http.StatusCreated, about)
	}
}

// プロフィールを更新
func PatchAbout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なID")
			return
		}

		var req models.About
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("JSONデコードエラー:", err)
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト: JSONの形式が正しくありません")
			return
		}

		// 名前が空でないかチェック
		if req.Name == "" {
			respondWithError(w, http.StatusBadRequest, "名前が空です")
			return
		}

		// データベースを更新
		updated, err := database.PatchAbout(db, id, req)
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

// 全プロフィールを取得
func GetAllAbout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		abouts, err := database.GetAllAbout(db)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		respondWithJSON(w, http.StatusOK, abouts)
	}
}

// プロフィールを取得
func GetAboutByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なID")
			return
		}

		about, err := database.GetAboutByID(db, id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		if about == nil {
			respondWithError(w, http.StatusNotFound, "About not found")
			return
		}
		respondWithJSON(w, http.StatusOK, about)
	}
}
