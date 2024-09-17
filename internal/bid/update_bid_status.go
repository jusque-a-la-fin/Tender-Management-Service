package bid

import "fmt"

// UpdateBidStatus изменяет статус предложения по его уникальному идентификатору
func (repo *BidDBRepository) UpdateBidStatus(bidID, username string, newStatus StatusEnum) (*Bid, int, error) {
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
        UPDATE bid
        SET status = $1
        WHERE id = $2 AND author_id = $3;
    `

	result, err := repo.dtb.Exec(query, newStatus, bidID, userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: обновление статуса предложения: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка от метода `RowsAffected`, пакет sql: %v", err)
	}

	if rowsAffected == 0 {
		return nil, 403, nil
	}

	bid, err := GetBid(repo.dtb, bidID)
	if err != nil {
		return nil, -1, err
	}
	return bid, 200, nil
}
