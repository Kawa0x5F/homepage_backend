package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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

// SQLファイルを読み込んで初期化する
func ApplySchema(db *sql.DB, schemaFilePath string) error {
	// schema.sql ファイルの読み込み
	schema, err := ioutil.ReadFile(schemaFilePath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	// SQL ステートメントを実行
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to apply schema: %v", err)
	}

	log.Println("Schema applied successfully")
	return nil
}
