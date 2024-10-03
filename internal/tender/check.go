package tender

import (
	"database/sql"
	"fmt"
)

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

// CheckTenderAndVersion проверяет, существует ли версия тендера, к которой нужно откатить тендер
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

// CheckRights проверяет, есть ли у пользователя права на работу с тендером
func CheckRights(dtb *sql.DB, tenderID, userID string) (bool, error) {
	query := `
	         SELECT EXISTS 
	            (SELECT 1 
				FROM tender
                WHERE id = $1 AND user_id = $2) 
				AS result`

	var hasRights bool
	err := dtb.QueryRow(query, tenderID, userID).Scan(&hasRights)
	if err != nil {
		return false, fmt.Errorf("ошибка запроса к базе данных: проверка прав доступа: %v", err)
	}

	return hasRights, nil
}
