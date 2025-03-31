package database

import (
	"database/sql"
	"errors"
	"kawa_blog/models"
	"log"
)

// 作品の一覧をDBから取得する
func GetAllProduct(db *sql.DB) ([]models.Product, error) {
	query := `
	SELECT id, title, description, image_url, github, blog, updated_at
	FROM product
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Title, &product.Description, &product.ImageURL, &product.Github, &product.Blog, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// idごとのデータを取得するコード
func GetProductByID(db *sql.DB, id int) (*models.Product, error) {
	var product models.Product
	query := `
	SELECT id, title, description, image_url, github, blog
	FROM product
	WHERE id = $1
	`

	err := db.QueryRow(query, id).Scan(&product.ID, &product.Title, &product.Description, &product.ImageURL, &product.Github, &product.Blog)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("IDのデータが見つかりません：", id)
			return nil, nil
		}
		log.Println("作品情報の取得に失敗しました：", err)
		return nil, err
	}

	return &product, nil
}

// 作品をDBに追加
func InsertProduct(db *sql.DB, product *models.Product) error {
	query := `
	INSERT INTO product (title, description, image_url, github, blog, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(query, product.Title, product.Description, product.ImageURL, product.Github, product.Blog).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		log.Println("Productの挿入に失敗:", err)
		return err
	}
	return nil
}

// 作品情報を更新
func PatchProduct(db *sql.DB, id int, product models.Product) (bool, error) {
	query := `
	UPDATE product
	SET title = $1, description = $2, image_url = $3, github = $4, blog = $5, updated_at = NOW()
	WHERE id = $6
	`
	result, err := db.Exec(query, product.Title, product.Description, product.ImageURL, product.Github, product.Blog, id)

	if err != nil {
		log.Println("作品情報の更新に失敗しました:", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("更新結果の取得に失敗しました:", err)
		return false, err
	}

	if rowsAffected == 0 {
		log.Println("更新対象の作品が見つかりません:", id)
		return false, nil
	}

	log.Println("作品が更新されました:", id)
	return true, nil
}

// DBからデータを削除する
func DeleteProductByID(db *sql.DB, id int) (bool, error) {
	result, err := db.Exec("DELETE FROM product WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	// 削除された行数を確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, err
	}

	return true, nil
}
