package bid

import (
	"database/sql"
	"fmt"
)

// RollbackBid откатывает параметры предложения к указанной версии
func (repo *BidDBRepository) RollbackBid(bri BidRollbackInput) (*Bid, int, error) {
	valid, err := checkAuthorName(repo.dtb, bri.Username)
	if !valid || err != nil {
		return nil, 401, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", bri.Username).Scan(&userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	valid, err = checkEditionRights(repo.dtb, bri.BidID, userID)
	if !valid || err != nil {
		return nil, 403, err
	}

	valid, err = сheckBid(repo.dtb, bri.BidID)
	if !valid || err != nil {
		return nil, 404, err
	}

	err = swapParams(repo.dtb, bri.Version, bri.BidID)
	if err != nil {
		return nil, -1, err
	}

	bri.Version++
	query := `
		     UPDATE bid
		     SET current_version = $1
		     WHERE id = $2;`

	_, err = repo.dtb.Exec(query, bri.Version, bri.BidID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: обновление текущей версии предложения: %v", err)
	}

	bid, err := GetBid(repo.dtb, bri.BidID)
	if err != nil {
		return nil, -1, err
	}
	return bid, 200, nil
}

// swapParams меняет местами версии предложения в таблице bid_versions
func swapParams(dtb *sql.DB, rollbackVersion int32, bidID string) error {
	var swappedVersion int32 = rollbackVersion + 1
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM bid_versions
            WHERE bid_id = $1 AND version = $2
        );`

	err := dtb.QueryRow(query, bidID, rollbackVersion).Scan(&exists)
	if err != nil {
		return fmt.Errorf(`ошибка запроса к базе данных: проверка существования версии предложения, 
		которая больше на 1, чем та версия(до инкремента), к которой нужно откатить предложение: %v`, err)
	}

	if !exists {
		return nil
	}

	rows, err := dtb.Query(`
        SELECT id, name, description
        FROM bid_versions
        WHERE bid_id = $1 AND version IN ($2, $3)`, bidID, rollbackVersion, swappedVersion)
	if err != nil {
		return fmt.Errorf("ошибка запроса к базе данных: извлечение параметров двух версий предложения: %v", err)
	}
	defer rows.Close()

	versions := make(map[int32]struct {
		ID          int
		Name        string
		Description string
	})

	for rows.Next() {
		var id int
		var name, description string
		var version int32

		if err := rows.Scan(&id, &name, &description); err != nil {
			return fmt.Errorf("ошибка от метода `Scan`, пакет sql: %v", err)
		}

		versions[version] = struct {
			ID          int
			Name        string
			Description string
		}{ID: id, Name: name, Description: description}
	}

	if verRB, okRB := versions[rollbackVersion]; okRB {
		if verSW, okSW := versions[swappedVersion]; okSW {
			_, err := dtb.Exec(`
                UPDATE bid_versions
                SET name = CASE
                    WHEN version = $1 THEN $2
                    WHEN version = $3 THEN $4
                END,
                description = CASE
                    WHEN version = $1 THEN $5
                    WHEN version = $3 THEN $6
                END
                WHERE tender_id = $7 AND version IN ($1, $3)`,
				rollbackVersion, verSW.Name, swappedVersion, verRB.Name,
				verSW.Description, verRB.Description,
				bidID)

			if err != nil {
				return fmt.Errorf("ошибка запроса к базе данных: обмен параметрами между двумя версиями предложения: %v", err)
			}
		}
	}
	return nil
}
