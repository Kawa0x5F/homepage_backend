package handlers

import (
	"encoding/json"
	"kawa_blog/cloud"
	"net/http"
)

// ファイルをアップロードするハンドラー
func UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		uploadedFile, header, err := r.FormFile("file") // "file" はフロントエンド側のキー名
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト: 画像が取得できません")
			return
		}
		defer uploadedFile.Close()

		// ファイル名を取得
		filename := header.Filename

		var imageURL string
		imageURL, err = cloud.UploadFile(filename, uploadedFile)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "データベースエラー")
			return
		}

		// レスポンスを JSON で返す
		response := map[string]string{"message": "File uploaded successfully", "image_url": imageURL}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
