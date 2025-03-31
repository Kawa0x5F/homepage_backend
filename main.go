package main

import (
	"kawa_blog/database"
	"kawa_blog/routes"
	"log"
	"net/http"

	"kawa_blog/client"
	"kawa_blog/utils"

	"github.com/gorilla/handlers"
)

func main() {
	// 環境変数を読み込む
	utils.LoadEnv()

	// Cloudflare R2 の S3 クライアントを初期化
	if err := client.InitS3Client(); err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	// DB接続
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.SeedAbout(db)

	// ルーター設定
	router := routes.NewRouter(db)

	// CORSの設定を適用
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://kawa0x5f.com/"}),                             // 許可するオリジン
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}), // 許可するHTTPメソッド
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),                    // 許可するヘッダー
		handlers.AllowCredentials(),
	)

	// サーバー起動
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", corsOptions(router)))
}
