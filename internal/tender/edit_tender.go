package tender

import (
	"fmt"
	"strings"
)

// EditTender изменяет параметры существующего тендера
func (repo *TenderDBRepository) EditTender(tei TenderEditionInput, tenderID, username string) (*Tender, int, error) {
	valid, err := CheckUser(repo.dtb, username)
	if !valid || err != nil {
		return nil, 401, err
	}

	valid, err = CheckTender(repo.dtb, tenderID)
	if !valid || err != nil {
		return nil, 404, err
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

	query := `
	        SELECT current_version
		    FROM tender
		    WHERE id = $1;`

	var currentVersion int
	err = repo.dtb.QueryRow(query, tenderID).Scan(&currentVersion)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение текущей версии тендера: %v", err)
	}

	currentVersion++

	query = `
		     UPDATE tender
		     SET current_version = $1
		     WHERE id = $2;`

	_, err = repo.dtb.Exec(query, currentVersion, tenderID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: обновление текущей версии тендера: %v", err)
	}

	// maxArgs - максимальное количество аргументов
	maxArgs := 5
	args := make([]interface{}, 0, maxArgs)

	query = "INSERT INTO tender_versions (version, "
	args = append(args, currentVersion)
	noChanges := true
	counter := 1
	switch {
	case tei.Name != "":
		args = append(args, tei.Name)
		query = fmt.Sprintf("%sname, ", query)
		noChanges = false
		counter++

	case tei.Description != "":
		args = append(args, tei.Description)
		query = fmt.Sprintf("%sdescription, ", query)
		noChanges = false
		counter++

	case tei.ServiceType != "":
		args = append(args, string(tei.ServiceType))
		query = fmt.Sprintf("%sservice_type, ", query)
		noChanges = false
		counter++
	}

	if noChanges {
		tnd, err := GetTender(repo.dtb, tenderID)
		if err != nil {
			return nil, -1, err
		}
		return tnd, 200, nil
	}

	counter++
	args = append(args, tenderID)
	query = fmt.Sprintf("%stender_id)", query)
	query = fmt.Sprintf("%s VALUES (", query)

	for cnt := 1; cnt <= counter; cnt++ {
		query = fmt.Sprintf("%s$%d, ", query, cnt)
	}

	query = strings.TrimSuffix(query, ", ")
	query = fmt.Sprintf("%s);", query)

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

	tnd, err := GetTender(repo.dtb, tenderID)
	if err != nil {
		return nil, -1, err
	}

	return tnd, 200, nil
}
