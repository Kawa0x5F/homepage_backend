package main

import (
	"kawa_blog/database"
	"kawa_blog/routes"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	// DB接続
	database, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// ルーター設定
	router := routes.NewRouter(database)

	// CORSの設定を適用
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),                   // 許可するオリジン
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // 許可するHTTPメソッド
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // 許可するヘッダー
	)

	// サーバー起動
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsOptions(router)))
}
