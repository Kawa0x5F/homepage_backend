package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"kawa_blog/utils"
)

// リクエストの構造
type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// エラーレスポンスの構造
type ErrorResponse struct {
	Error string `json:"error"`
}

// LoginHandler はログイン処理を行う
func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONResponse(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "Invalid request method"})
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
			return
		}

		adminID := utils.GetEnv("ADMIN_ID")
		adminPassword := utils.GetEnv("ADMIN_PASSWORD")

		// 認証情報が一致するか確認
		if req.ID != adminID || req.Password != adminPassword {
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
			return
		}

		// セキュアなトークンを生成
		token, err := utils.GenerateToken()
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token"})
			return
		}

		// Cookie に保存
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Now().Add(1 * time.Hour),
		})

		writeJSONResponse(w, http.StatusOK, map[string]string{"result": "success"})
	}
}

// 認証チェック用のハンドラー
func AuthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Cookieの取得
		cookie, err := r.Cookie("admin_token")
		if err != nil {
			// Cookieがない場合は認証エラー
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
			return
		}

		// トークンの検証
		if !utils.ValidateToken(cookie.Value) {
			// トークンが無効な場合
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
			return
		}

		// 認証成功時のレスポンス
		writeJSONResponse(w, http.StatusOK, map[string]string{"authenticated": "true"})
	}
}

// JSONレスポンスを返すための関数
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
