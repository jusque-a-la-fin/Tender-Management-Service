package bid

import (
	"fmt"
	"strings"
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

	query := `
	        SELECT current_version
		    FROM bid
		    WHERE id = $1;`

	var currentVersion int
	err = repo.dtb.QueryRow(query, bidID).Scan(&currentVersion)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение текущей версии предложения: %v", err)
	}

	currentVersion++

	query = `
		     UPDATE bid
		     SET current_version = $1
		     WHERE id = $2;`

	_, err = repo.dtb.Exec(query, currentVersion, bidID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: обновление текущей версии предложения: %v", err)
	}

	// maxArgs - максимальное количество аргументов
	maxArgs := 4
	args := make([]interface{}, 0, maxArgs)
	query = "INSERT INTO bid_versions (version, "
	args = append(args, currentVersion)
	noChanges := true
	counter := 1

	if bdi.Name != "" {
		args = append(args, bdi.Name)
		query = fmt.Sprintf("%sname, ", query)
		noChanges = false
		counter++
	}

	if bdi.Description != "" {
		args = append(args, bdi.Description)
		query = fmt.Sprintf("%sdescription, ", query)
		noChanges = false
		counter++
	}

	if noChanges {
		bid, err := GetBid(repo.dtb, bidID)
		if err != nil {
			return nil, -1, err
		}
		return bid, 200, nil
	}

	counter++
	args = append(args, bidID)
	query = fmt.Sprintf("%sbid_id)", query)
	query = fmt.Sprintf("%s VALUES (", query)

	for cnt := 1; cnt <= counter; cnt++ {
		query = fmt.Sprintf("%s$%d, ", query, cnt)
	}

	query = strings.TrimSuffix(query, ", ")
	query = fmt.Sprintf("%s);", query)

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

	bid, err := GetBid(repo.dtb, bidID)
	if err != nil {
		return nil, -1, err
	}
	return bid, 200, nil
}
