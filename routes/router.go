package routes

import (
	"database/sql"
	"kawa_blog/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/articles/all", handlers.GetAllArticles(db)).Methods("GET")
	r.HandleFunc("/articles", handlers.CreateArticle(db)).Methods("POST")
	r.HandleFunc("/articles/publish", handlers.GetPublishArticles(db)).Methods("GET")
	r.HandleFunc("/article/{slug}", handlers.GetArticle(db)).Methods("GET")
	r.HandleFunc("/article/{slug}", handlers.PatchArticle(db)).Methods("PATCH")
	r.HandleFunc("/article/{slug}", handlers.DeleteArticle(db)).Methods("DELETE")
	r.HandleFunc("/tags", handlers.GetTags(db)).Methods("GET")
	r.HandleFunc("/tags", handlers.CreateTag(db)).Methods("POST")

	r.HandleFunc("/login", handlers.LoginHandler()).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler()).Methods("POST")
	r.HandleFunc("/auth/check", handlers.AuthCheckHandler()).Methods("GET")

	return r
}
