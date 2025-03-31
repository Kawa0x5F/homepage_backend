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

// 全スキルデータを取得するハンドラー
func GetAllSkills(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		skills, err := database.GetAllSkills(db)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		respondWithJSON(w, http.StatusOK, skills)
	}
}

// スキルを追加するハンドラー
func CreateSkill(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var skill models.Skills
		if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト")
			return
		}

		if err := database.InsertSkills(db, &skill); err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		respondWithJSON(w, http.StatusCreated, skill)
	}
}

// スキルを削除するハンドラー
func DeleteSkillByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なID")
			return
		}

		deleted, err := database.DeleteSkillByID(db, id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		if !deleted {
			respondWithError(w, http.StatusNotFound, "指定されたスキルがありませんでした")
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Skill deleted"})
	}
}
