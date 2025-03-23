package handlers

import (
	"encoding/json"
	"net/http"

	"kawa_blog/utils"
)

// LoginRequest はリクエストの構造
type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Loginの処理を行う
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

		if req.ID == adminID && req.Password == adminPassword {
			writeJSONResponse(w, http.StatusOK, map[string]string{"result": "success"})
		} else {
			writeJSONResponse(w, http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		}
	}
}

// errorをjsonで返すための関数
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
