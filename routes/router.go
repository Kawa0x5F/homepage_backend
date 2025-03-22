package routes

import (
	"database/sql"
	"kawa_blog/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/articles", handlers.GetArticles(db)).Methods("GET")
	r.HandleFunc("/articles", handlers.CreateArticle(db)).Methods("POST")
	r.HandleFunc("/article/{id}", handlers.GetArticle(db)).Methods("GET")
	r.HandleFunc("/tags", handlers.GetTags(db)).Methods("GET")
	r.HandleFunc("/tags", handlers.CreateTag(db)).Methods("POST")
	return r
}
