package bid

import (
	"database/sql"
	"fmt"
)

// checkAuthor проверяет, существует ли пользователь/организация или корректен ли он/она.
func checkAuthor(dtb *sql.DB, authorId string, authorType AuthorTypeEnum) (bool, error) {
	query := ``
	var errStr string
	switch authorType {
	case User:
		query = `SELECT EXISTS (SELECT 1 FROM employee WHERE id = $1)`
		errStr = "ошибка запроса к базе данных: проверка корректности пользователя: "
	case Organization:
		query = `SELECT EXISTS (SELECT 1 FROM organization WHERE id = $1)`
		errStr = "ошибка запроса к базе данных: проверка корректности организации: "
	}

	var exists bool
	err := dtb.QueryRow(query, authorId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s%v", errStr, err)
	}

	return exists, nil
}

// checkCreationRights проверяет, достаточно ли прав для создания предложения
func checkCreationRights(dtb *sql.DB, authorId string, authorType AuthorTypeEnum) (bool, error) {
	if authorType == "Organization" {
		return true, nil
	}

	var hasRights bool

	query := `
        SELECT COUNT(DISTINCT organization_id) = 1 AS single_responsible
        FROM organization_responsible
        WHERE user_id = $1;
    `

	err := dtb.QueryRow(query, authorId).Scan(&hasRights)
	if err != nil {
		return false, fmt.Errorf("ошибка запроса к базе данных: проверка прав для создания предложения: %v", err)
	}

	return hasRights, nil
}

// checkEditionRights проверяет, достаточно ли прав для изменения предложения
func checkEditionRights(dtb *sql.DB, bidID, authorID string) (bool, error) {
	var hasRights bool
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM bid
            WHERE id = $1 AND author_id = $2
        );
    `
	err := dtb.QueryRow(query, bidID, authorID).Scan(&hasRights)
	if err != nil {
		return false, fmt.Errorf("ошибка запроса к базе данных: проверка прав для изменения предложения: %v", err)
	}
	return hasRights, nil
}

// checkUsername проверяет, существует ли c таким именем пользователь/организация или корректен ли он/она.
func checkUsername(dtb *sql.DB, username string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM employee WHERE username = $1)`
	var exists bool
	err := dtb.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s%v", "ошибка запроса к базе данных: проверка корректности пользователя: ", err)
	}

	if !exists {
		query = `SELECT EXISTS (SELECT 1 FROM organization WHERE name = $1)`
		err := dtb.QueryRow(query, username).Scan(&exists)
		if err != nil {
			return false, fmt.Errorf("%s%v", "ошибка запроса к базе данных: проверка корректности организации: ", err)
		}
	}

	return exists, nil
}

// checkBid проверяет, существует ли предложение
func сheckBid(dtb *sql.DB, bidID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM bid WHERE id = $1)`
	err := dtb.QueryRow(query, bidID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка запроса к базе данных: проверка существования предложения: %v", err)
	}

	return exists, nil
}
