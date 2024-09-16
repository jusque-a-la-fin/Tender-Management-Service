package tender

import (
	"database/sql"
	"fmt"
)

// GetTenderStatus получает статус тендера по его уникальному идентификатору
func (repo *TenderDBRepository) GetTenderStatus(username, tenderID string) (string, int, error) {
	valid, err := CheckUser(repo.dtb, username)
	if !valid || err != nil {
		return "", 401, err
	}

	valid, err = CheckTender(repo.dtb, tenderID)
	if !valid || err != nil {
		return "", 404, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return "", -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var status string
	err = repo.dtb.QueryRow(`SELECT status FROM tender WHERE id = $1 AND user_id = $2`, tenderID, userID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 403, nil
		}
		return "", -1, fmt.Errorf("ошибка запроса к базе данных: извлечение status тендера: %v", err)
	}

	return status, 200, nil
}
