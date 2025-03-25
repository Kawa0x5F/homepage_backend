package handlers

import (
	"encoding/json"
	"kawa_blog/cloud"
	"kawa_blog/models"
	"net/http"
	"net/url"
	"path"
)

// ファイルをアップロードするハンドラー
func UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		uploadedFile, header, err := r.FormFile("file")
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

// ファイルをデリートするハンドラー
func DeleteFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.FileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト")
			return
		}

		// 削除対象のURLが空でないか確認
		if req.ImageURL == "" {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト: 画像URLが空です")
			return
		}

		// ファイル名を抽出
		fileName, err := extractFileName(req.ImageURL)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト：画像のファイル名の抽出に失敗しました")
			return
		}

		// ファイルを削除
		if err := cloud.DeleteFile(fileName); err != nil {
			respondWithError(w, http.StatusInternalServerError, "画像の削除に失敗しました")
			return
		}

		// 成功レスポンス
		w.WriteHeader(http.StatusNoContent)
	}
}

// ファイルを更新（削除してアップロード）するハンドラー
func PatchFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// アップロードされた画像を受け取る
		uploadedFile, header, err := r.FormFile("file")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト: 画像が取得できません")
			return
		}
		defer uploadedFile.Close()

		// 古い画像URLを取得
		imageURL := r.FormValue("image_url")
		if imageURL == "" {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト: 画像URLがありません")
			return
		}
		// ファイル名を抽出
		imageName, err := extractFileName(imageURL)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "無効なリクエスト：画像のファイル名の抽出に失敗しました")
			return
		}

		filename := header.Filename
		newImageURL, err := cloud.UploadFile(filename, uploadedFile)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "アップロードエラー")
			return
		}

		if err := cloud.DeleteFile(imageName); err != nil {
			respondWithError(w, http.StatusInternalServerError, "古い画像の削除に失敗しました")
			return
		}

		// レスポンスを JSON で返す
		response := map[string]string{
			"message":   "File uploaded successfully",
			"image_url": newImageURL,
		}
		respondWithJSON(w, http.StatusCreated, response)
	}
}

// URLからファイル名を抽出する関数
func extractFileName(imageURL string) (string, error) {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return "", err
	}
	return path.Base(parsedURL.Path), nil
}
