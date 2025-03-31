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

	r.HandleFunc("/tags/all", handlers.GetALLTags(db)).Methods("GET")
	r.HandleFunc("/tags", handlers.CreateTag(db)).Methods("POST")
	r.HandleFunc("/tags/article", handlers.InsertArticleTags(db)).Methods("POST")
	r.HandleFunc("/tags/article/{slug}", handlers.UpdateArticleTags(db)).Methods("PATCH")

	r.HandleFunc("/login", handlers.LoginHandler()).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler()).Methods("POST")
	r.HandleFunc("/auth/check", handlers.AuthCheckHandler()).Methods("GET")

	r.HandleFunc("/image", handlers.UploadFile()).Methods("POST")
	r.HandleFunc("/image", handlers.DeleteFile()).Methods("DELETE")
	r.HandleFunc("/image", handlers.PatchFile()).Methods("PATCH")

	r.HandleFunc("/about", handlers.CreateAbout(db)).Methods("POST")
	r.HandleFunc("/about/all", handlers.GetAllAbout(db)).Methods("GET")
	r.HandleFunc("/about/{id:[0-9]+}", handlers.GetAboutByID(db)).Methods("GET")
	r.HandleFunc("/about/{id:[0-9]+}", handlers.PatchAbout(db)).Methods("PATCH")

	r.HandleFunc("/skills", handlers.CreateSkill(db)).Methods("POST")
	r.HandleFunc("/skills/all", handlers.GetAllSkills(db)).Methods("GET")
	r.HandleFunc("/skills/{id:[0-9]+}", handlers.DeleteSkillByID(db)).Methods("DELETE")

	r.HandleFunc("/contact", handlers.CreateContact(db)).Methods("POST")
	r.HandleFunc("/contact/all", handlers.GetAllContact(db)).Methods("GET")
	r.HandleFunc("/contact/{id:[0-9]+}", handlers.GetContactByID(db)).Methods("GET")
	r.HandleFunc("/contact/{id:[0-9]+}", handlers.PatchContact(db)).Methods("PATCH")
	r.HandleFunc("/contact/{id:[0-9]+}", handlers.DeleteContactByID(db)).Methods("DELETE")

	r.HandleFunc("/product", handlers.CreateProduct(db)).Methods("POST")
	r.HandleFunc("/product/all", handlers.GetAllProduct(db)).Methods("GET")
	r.HandleFunc("/product/{id:[0-9]+}", handlers.GetProductByID(db)).Methods("GET")
	r.HandleFunc("/product/{id:[0-9]+}", handlers.PatchProduct(db)).Methods("PATCH")
	r.HandleFunc("/product/{id:[0-9]+}", handlers.DeleteProductByID(db)).Methods("DELETE")

	return r
}
