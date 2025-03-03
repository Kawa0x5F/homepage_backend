package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// DB接続の初期化
func InitDB() (*sql.DB, error) {
	// 環境変数から接続情報を取得
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// PostgreSQL の接続文字列
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// DBに接続
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB接続エラー:", err)
		return nil, err
	}

	// 接続確認
	if err := db.Ping(); err != nil {
		log.Fatal("DB接続確認エラー:", err)
		return nil, err
	}

	log.Println("DB接続成功")
	return db, nil
}
