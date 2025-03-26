package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv は環境変数をロードする
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env ファイルが見つかりません。環境変数を直接使用します。")
	}
}

// GetEnv は環境変数を取得する
func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetSecretKey() string {
	return os.Getenv("SECRET_KEY")
}
