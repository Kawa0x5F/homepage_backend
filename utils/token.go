package utils

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

// トークンを一時保存するマップ
var (
	tokenStore = make(map[string]bool)
	mutex      = &sync.Mutex{}
)

// ランダムなトークンを生成
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)

	// 保存（ログアウト時に削除するため）
	mutex.Lock()
	tokenStore[token] = true
	mutex.Unlock()

	return token, nil
}

// トークンを検証
func ValidateToken(token string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, exists := tokenStore[token]
	return exists
}

// ログアウト時にトークンを削除
func InvalidateToken(token string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(tokenStore, token)
}
