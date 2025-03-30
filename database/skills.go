package database

import (
	"database/sql"
	"kawa_blog/models"
	"log"
)

// 全データを取得するコード
func GetAllSkills(db *sql.DB) ([]models.Skills, error) {
	query := `
	SELECT id, name, type, has_image, created_at
	FROM skills
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skills

	for rows.Next() {
		var skill models.Skills
		err := rows.Scan(&skill.ID, &skill.Name, &skill.Type, &skill.HasImage, &skill.CreatedAt)
		if err != nil {
			return nil, err
		}

		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}

// DBにデータを追加する
func InsertSkills(db *sql.DB, skills *models.Skills) error {
	query := `
	INSERT INTO skills (name, type, has_image, created_at)
	VALUES($1, $2, $3, NOW()) RETURNING id, created_at
	`

	err := db.QueryRow(query, skills.Name, skills.Type, skills.HasImage).Scan(&skills.ID, &skills.CreatedAt)
	if err != nil {
		log.Println("Skillの挿入に失敗:", err)
		return err
	}
	return nil
}

// DBからデータを削除する
func DeleteSkillByID(db *sql.DB, id int) (bool, error) {
	result, err := db.Exec("DELETE FROM skills WHERE id = $1", id)
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
