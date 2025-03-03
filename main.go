package main

import (
	"kawa_blog/database"
	"kawa_blog/routes"
	"log"
	"net/http"
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

	// サーバー起動
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
