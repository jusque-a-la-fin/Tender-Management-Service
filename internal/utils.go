package internal

import (
	"database/sql"
	"fmt"
)

// CheckUser проверяет, существует ли пользователь или корректен ли он.
func CheckUser(dtb *sql.DB, creatorUsername string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM employee WHERE username = $1)`
	var exists bool
	err := dtb.QueryRow(query, creatorUsername).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка запроса к базе данных: проверка корректности пользователя: %v", err)
	}

	return exists, nil
}
