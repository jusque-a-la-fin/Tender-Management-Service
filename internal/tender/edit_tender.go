package tender

import (
	"database/sql"
	"fmt"
)

// EditTender изменяет параметры существующего тендера
func (repo *TenderDBRepository) EditTender(tei TenderEditionInput, tenderID, username string) (*Tender, int, error) {
	valid, err := CheckUser(repo.dtb, username)
	if !valid || err != nil {
		return nil, 401, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	valid, err = CheckRights(repo.dtb, tenderID, userID)
	if !valid || err != nil {
		return nil, 403, err
	}

	valid, err = CheckTender(repo.dtb, tenderID)
	if !valid || err != nil {
		return nil, 404, err
	}

	query := `
		SELECT MAX(version)
		FROM tender_versions
		WHERE tender_id = $1;
	`

	var latestVersion int32
	err = repo.dtb.QueryRow(query, tenderID).Scan(&latestVersion)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение последней версии тендера: %v", err)
	}

	// maxArgs - максимальное количество аргументов
	maxArgs := 5
	args := make([]interface{}, maxArgs)
	query = "INSERT INTO tender_versions (version, name, description, service_type, tender_id) VALUES ($1, $2, $3, $4, $5);"
	noChanges := true

	if tei.Name == "" {
		args[1], err = getParam(repo.dtb, "name", tenderID, latestVersion)
		if err != nil {
			return nil, -1, err
		}
		noChanges = false
	} else {
		args[1] = tei.Name
	}

	if tei.Description == "" {
		args[2], err = getParam(repo.dtb, "description", tenderID, latestVersion)
		if err != nil {
			return nil, -1, err
		}
		noChanges = false
	} else {
		args[2] = tei.Description
	}

	if tei.ServiceType == "" {
		args[3], err = getParam(repo.dtb, "service_type", tenderID, latestVersion)
		if err != nil {
			return nil, -1, err
		}
		noChanges = false
	} else {
		args[3] = tei.ServiceType
	}

	latestVersion++
	args[0] = latestVersion
	args[4] = tenderID

	if noChanges {
		tnd, err := GetTender(repo.dtb, tenderID)
		if err != nil {
			return nil, -1, err
		}
		return tnd, 200, nil
	}

	result, err := repo.dtb.Exec(query, args...)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: изменение параметров существующего тендера: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка от метода `RowsAffected`, пакет sql: %v", err)
	}

	if rowsAffected == 0 {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: не обновились параметры существующего тендера")
	}

	query = `UPDATE tender
	         SET current_version = $1
	         WHERE id = $2;`

	_, err = repo.dtb.Exec(query, latestVersion, tenderID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: обновление текущей версии тендера: %v", err)
	}

	tnd, err := GetTender(repo.dtb, tenderID)
	if err != nil {
		return nil, -1, err
	}

	return tnd, 200, nil
}

func getParam(dtb *sql.DB, param, tenderID string, version int32) (string, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM tender_versions
		WHERE tender_id = $1 AND version = $2;`, param)

	var paramVal string
	err := dtb.QueryRow(query, tenderID, version).Scan(&paramVal)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса к базе данных: извлечение параметра, который не должен быть изменен: %v", err)
	}
	return paramVal, nil
}
