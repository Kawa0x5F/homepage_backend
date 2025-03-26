package database

import (
	"database/sql"
	"errors"
	"fmt"
	"kawa_blog/models"
	"kawa_blog/utils"
	"log"
)

// デフォルトのデータを挿入するコード
func SeedAbout(db *sql.DB) error {
	var count int64
	err := db.QueryRow("SELECT COUNT(*) FROM about").Scan(&count)
	if err != nil {
		log.Println("デフォルトデータの存在チェックに失敗:", err)
		return err
	}

	if count == 0 {
		fmt.Println("デフォルトデータを挿入します...")
		query := `
			INSERT INTO about (name, roma, description, image_url, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
		`
		_, err = db.Exec(query, utils.GetEnv("PORTFOLIO_NAME"), utils.GetEnv("PORTFOLIO_ROMA"), utils.GetEnv("PORTFOLIO_DESCRIPTION"), utils.GetEnv("PORTFOLIO_IMAGE_URL"))
		if err != nil {
			log.Println("デフォルトデータの挿入に失敗:", err)
			return err
		}
	} else {
		fmt.Println("データは既に存在しています。")
	}
	return nil
}

// idごとのデータを取得するコード
func GetAboutByID(db *sql.DB, id int) (*models.About, error) {
	var about models.About
	query := `
	SELECT id, name, roma, description, image_url
	FROM about
	WHERE id = $1
	`

	err := db.QueryRow(query, id).Scan(&about.ID, &about.Name, &about.Roma, &about.Description, &about.ImageURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("IDのデータが見つかりません：", id)
			return nil, nil
		}
		log.Println("Aboutデータの取得に失敗しました：", err)
		return nil, err
	}

	return &about, nil
}

// プロフィールデータをデータベースに追加
func InsertAbout(db *sql.DB, about *models.About) error {
	query := `
	INSERT INTO about (name, roma, description, image_url, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(query, about.Name, about.Roma, about.Description, about.ImageURL).Scan(&about.ID, &about.CreatedAt, &about.UpdatedAt)
	if err != nil {
		log.Println("Aboutの挿入に失敗:", err)
		return err
	}
	return nil
}

// プロフィールデータを更新
func PatchAbout(db *sql.DB, id int, about models.About) (bool, error) {
	query := `
	UPDATE about
	SET name = $1, roma = $2, description = $3, image_url = $4, updated_at = NOW()
	WHERE id = $5
	`
	result, err := db.Exec(query, about.Name, about.Roma, about.Description, about.ImageURL, id)

	if err != nil {
		log.Println("プロフィールの更新に失敗しました:", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("更新結果の取得に失敗しました:", err)
		return false, err
	}

	if rowsAffected == 0 {
		log.Println("更新対象のプロフィールが見つかりません:", id)
		return false, nil
	}

	log.Println("プロフィールが更新されました:", id)
	return true, nil
}
