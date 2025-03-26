package handlers

import (
	"net/http"
	"time"

	"kawa_blog/utils"
)

// LogoutHandler はログアウト処理を行う
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_token")
		if err == nil {
			utils.InvalidateToken(cookie.Value)
		}

		// Cookie を無効化
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(-1 * time.Second),
		})

		writeJSONResponse(w, http.StatusOK, map[string]string{"result": "logout successful"})
	}
}
