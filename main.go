package main

import (
	"kawa_blog/database"
	"kawa_blog/routes"
	"log"
	"net/http"

	"kawa_blog/utils"

	"github.com/gorilla/handlers"
)

func main() {
	// 環境変数を読み込む
	utils.LoadEnv()

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
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),                            // 許可するオリジン
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}), // 許可するHTTPメソッド
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),                    // 許可するヘッダー
		handlers.AllowCredentials(),
	)

	// サーバー起動
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", corsOptions(router)))
}
