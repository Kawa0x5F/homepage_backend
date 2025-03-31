package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// DB接続の初期化
func InitDB() (*sql.DB, error) {
	// DBに接続するためのURL
	connStr := os.Getenv("DATABASE_URL")

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
