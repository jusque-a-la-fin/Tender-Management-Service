package tender

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

// checkTender проверяет, существует ли тендер.
func CheckTender(dtb *sql.DB, tenderID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tender WHERE id = $1)`
	err := dtb.QueryRow(query, tenderID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка запроса к базе данных: проверка существования тендера: %v", err)
	}

	return exists, nil
}

// CheckTenderAndVersion проверяет, существует ли тендер или его версия
func CheckTenderAndVersion(dtb *sql.DB, version int32, tenderID string) (bool, error) {
	var exists bool
	query := `
             SELECT EXISTS (
                 SELECT 1
                 FROM tender_versions
                 WHERE tender_id = $1 AND version = $2
             );`

	err := dtb.QueryRow(query, tenderID, version).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка запроса к базе данных: проверка существования тендера и его версии: %v", err)
	}

	return exists, nil
}
