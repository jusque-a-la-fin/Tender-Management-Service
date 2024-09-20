package bid

import (
	"database/sql"
	"fmt"
)

// EditBid редактирует параметры предложения
func (repo *BidDBRepository) EditBid(bdi BidEditionInput, bidID, username string) (*Bid, int, error) {
	valid, err := checkUsername(repo.dtb, username)
	if !valid || err != nil {
		return nil, 401, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	valid, err = checkEditionRights(repo.dtb, bidID, userID)
	if !valid || err != nil {
		return nil, 403, err
	}

	valid, err = сheckBid(repo.dtb, bidID)
	if !valid || err != nil {
		return nil, 404, err
	}

	if bdi.Name == "" && bdi.Description == "" {
		bid, err := GetBid(repo.dtb, bidID)
		if err != nil {
			return nil, -1, err
		}
		return bid, 200, nil
	}

	query := `
		SELECT MAX(version)
		FROM bid_versions
		WHERE bid_id = $1;
	`

	var latestVersion int32
	err = repo.dtb.QueryRow(query, bidID).Scan(&latestVersion)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение последней версии предложения: %v", err)
	}

	// maxArgs - максимальное количество аргументов
	maxArgs := 4
	args := make([]interface{}, maxArgs)
	query = "INSERT bid_versions (version, name, description, bid_id) VALUES ($1, $2, $3, $4);"

	if bdi.Name == "" {
		args[1], err = getParam(repo.dtb, "name", bidID, latestVersion)
		if err != nil {
			return nil, -1, err
		}
	} else {
		args[1] = bdi.Name
	}

	if bdi.Description == "" {
		args[2], err = getParam(repo.dtb, "description", bidID, latestVersion)
		if err != nil {
			return nil, -1, err
		}
	} else {
		args[2] = bdi.Description
	}

	latestVersion++
	args[0] = latestVersion
	args[3] = bidID

	result, err := repo.dtb.Exec(query, args...)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: добавление параметров новой версии предложения: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка от метода `RowsAffected`, пакет sql: %v", err)
	}

	if rowsAffected == 0 {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: не добавились параметры новой версии предложения")
	}

	query = `
		     UPDATE bid
		     SET current_version = $1
		     WHERE id = $2;`

	_, err = repo.dtb.Exec(query, latestVersion, bidID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: обновление текущей версии предложения: %v", err)
	}

	bid, err := GetBid(repo.dtb, bidID)
	if err != nil {
		return nil, -1, err
	}
	return bid, 200, nil
}

func getParam(dtb *sql.DB, param, bidID string, version int32) (string, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM bid_versions
		WHERE bid_id = $1 AND version = $2;`, param)

	var paramVal string
	err := dtb.QueryRow(query, bidID, version).Scan(&paramVal)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса к базе данных: извлечение параметра, который не должен быть изменен: %v", err)
	}
	return paramVal, nil
}
