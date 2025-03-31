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

// 全作品情報を取得するハンドラー
func GetAllProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := database.GetAllProduct(db)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		respondWithJSON(w, http.StatusOK, products)
	}
}

// 作品情報を取得
func GetProductByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なID")
			return
		}

		product, err := database.GetProductByID(db, id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}
		if product == nil {
			respondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		respondWithJSON(w, http.StatusOK, product)
	}
}

// 情報を追加するハンドラー
func CreateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト")
			return
		}

		if err := database.InsertProduct(db, &product); err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		respondWithJSON(w, http.StatusCreated, product)
	}
}

// コンタクトを更新
func PatchProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なID")
			return
		}

		var req models.Product
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("JSONデコードエラー:", err)
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト: JSONの形式が正しくありません")
			return
		}

		// 名前とリンクが空でないかチェック
		if req.Title == "" || req.Description == "" {
			respondWithError(w, http.StatusBadRequest, "名前もしくは説明が空です")
			return
		}

		// データベースを更新
		updated, err := database.PatchProduct(db, id, req)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		if !updated {
			respondWithError(w, http.StatusNotFound, "Product not found")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product update"})
	}
}

// 作品情報を削除するハンドラー
func DeleteProductByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なID")
			return
		}

		deleted, err := database.DeleteProductByID(db, id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		if !deleted {
			respondWithError(w, http.StatusNotFound, "指定された作品がありませんでした")
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted"})
	}
}
