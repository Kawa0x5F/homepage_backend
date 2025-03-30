package database

import (
	"database/sql"
	"kawa_blog/models"
	"log"
)

// 全データを取得するコード
func GetAllContact(db *sql.DB) ([]models.Contact, error) {
	query := `
	SELECT id, name, link, has_image, created_at
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
		err := rows.Scan(&contact.ID, &contact.Name, &contact.Link, &contact.HasImage, &contact.CreatedAt)
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

// DBにデータを追加する
func InsertContact(db *sql.DB, contact *models.Contact) error {
	query := `
	INSERT INTO contact (name, link, has_image, created_at)
	VALUES($1, $2, $3, NOW()) RETURNING id, created_at
	`

	err := db.QueryRow(query, contact.Name, contact.Link, contact.HasImage).Scan(&contact.ID, &contact.CreatedAt)
	if err != nil {
		log.Println("Contactの挿入に失敗:", err)
		return err
	}
	return nil
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
