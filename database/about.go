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

		defaultData := []struct {
			Name        string
			Roma        string
			Description string
			Color       string
			ImageURL    string
		}{
			{
				utils.GetEnv("PORTFOLIO_NAME"),
				utils.GetEnv("PORTFOLIO_ROMA"),
				utils.GetEnv("PORTFOLIO_DESCRIPTION"),
				utils.GetEnv("PORTFOLIO_COLOR"),
				utils.GetEnv("PORTFOLIO_IMAGE_URL"),
			},
			{
				utils.GetEnv("PORTFOLIO_NAME_2"),
				utils.GetEnv("PORTFOLIO_ROMA_2"),
				utils.GetEnv("PORTFOLIO_DESCRIPTION_2"),
				utils.GetEnv("PORTFOLIO_COLOR_2"),
				utils.GetEnv("PORTFOLIO_IMAGE_URL_2"),
			},
		}

		query := `
			INSERT INTO about (name, roma, description, color, image_url, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		`

		for _, data := range defaultData {
			_, err = db.Exec(query, data.Name, data.Roma, data.Description, data.Color, data.ImageURL)
			if err != nil {
				log.Println("デフォルトデータの挿入に失敗:", err)
				return err
			}
		}
	} else {
		fmt.Println("データは既に存在しています。")
	}
	return nil
}

// 全データを取得するコード
func GetAllAbout(db *sql.DB) ([]models.About, error) {
	query := `
	SELECT id, name, roma, description, color, image_url, created_at, updated_at
	FROM about
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var abouts []models.About

	for rows.Next() {
		var about models.About
		err := rows.Scan(&about.ID, &about.Name, &about.Roma, &about.Description, &about.Color, &about.ImageURL, &about.CreatedAt, &about.UpdatedAt)
		if err != nil {
			return nil, err
		}

		abouts = append(abouts, about)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return abouts, nil
}

// idごとのデータを取得するコード
func GetAboutByID(db *sql.DB, id int) (*models.About, error) {
	var about models.About
	query := `
	SELECT id, name, roma, description, color, image_url
	FROM about
	WHERE id = $1
	`

	err := db.QueryRow(query, id).Scan(&about.ID, &about.Name, &about.Roma, &about.Description, &about.Color, &about.ImageURL)

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
	INSERT INTO about (name, roma, description, image_url, color, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(query, about.Name, about.Roma, about.Description, about.Color, about.ImageURL).Scan(&about.ID, &about.CreatedAt, &about.UpdatedAt)
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
	SET name = $1, roma = $2, description = $3, color = $4, image_url = $5, updated_at = NOW()
	WHERE id = $6
	`
	result, err := db.Exec(query, about.Name, about.Roma, about.Description, about.Color, about.ImageURL, id)

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
