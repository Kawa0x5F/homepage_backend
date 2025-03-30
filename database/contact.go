package database

import (
	"database/sql"
	"errors"
	"kawa_blog/models"
	"log"
)

// 全データを取得するコード
func GetAllContact(db *sql.DB) ([]models.Contact, error) {
	query := `
	SELECT id, name, link, has_image, created_at, updated_at
	FROM contact
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.Contact

	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(&contact.ID, &contact.Name, &contact.Link, &contact.HasImage, &contact.CreatedAt, &contact.UpdatedAt)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

// IDごとにデータを取得する
func GetContactByID(db *sql.DB, id int) (*models.Contact, error) {
	var contact models.Contact
	query := `
	SELECT id, name, link, has_image
	FROM contact
	WHERE id = $1
	`

	err := db.QueryRow(query, id).Scan(&contact.Name, &contact.Link, &contact.HasImage)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("IDのデータが見つかりません：", id)
			return nil, nil
		}
		log.Println("Contactデータの取得に失敗しました：", err)
		return nil, err
	}

	return &contact, nil
}

// DBにデータを追加する
func InsertContact(db *sql.DB, contact *models.Contact) error {
	query := `
	INSERT INTO contact (name, link, has_image, created_at, updated_at)
	VALUES($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(query, contact.Name, contact.Link, contact.HasImage).Scan(&contact.ID, &contact.CreatedAt, &contact.UpdatedAt)
	if err != nil {
		log.Println("Contactの挿入に失敗:", err)
		return err
	}
	return nil
}

// IDごとにデータを更新する
func PatchContact(db *sql.DB, id int, contact models.Contact) (bool, error) {
	query := `
	UPDATE contact
	SET name = $1, link = $2, has_image = $3, updated_at = NOW()
	WHERE id = $4
	`
	result, err := db.Exec(query, contact.Name, contact.Link, contact.HasImage)

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

// DBからデータを削除する
func DeleteContactByID(db *sql.DB, id int) (bool, error) {
	result, err := db.Exec("DELETE FROM contact WHERE id = $1", id)
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
