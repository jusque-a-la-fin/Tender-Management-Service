package bid

import "fmt"

// GetUserBids получает список предложений текущего пользователя
func (repo *BidDBRepository) GetUserBids(username string, startIndex, endIndex int32) ([]*Bid, int, error) {
	valid, err := checkUsername(repo.dtb, username)
	if !valid || err != nil {
		return nil, 401, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	query := `
        SELECT 
            b.id,
            b.status,
            b.tender_id,
            b.author_type,
            b.author_id,
            b.current_version,
            b.created_at,
            bv.name,
            bv.version
        FROM 
            bid b
        JOIN 
            bid_versions bv ON b.id = bv.bid_id
        WHERE 
            b.author_id = $1;`

	rows, err := repo.dtb.Query(query, userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение предложений текущего пользователя: %v", err)
	}
	defer rows.Close()

	var bids []*Bid
	for rows.Next() {
		bid := &Bid{}
		if err := rows.Scan(&bid.ID, &bid.Status, &bid.TenderID, &bid.AuthorType, &bid.AuthorID,
			&bid.Version, &bid.CreatedAt, &bid.Name, &bid.Version); err != nil {
			return nil, -1, fmt.Errorf("ошибка от метода `Scan`, пакет sql: %v", err)
		}
		bids = append(bids, bid)
	}

	if err := rows.Err(); err != nil {
		return nil, -1, fmt.Errorf("ошибка во время итерирования по строкам, возвращенным запросом: %v", err)
	}

	return bids, 200, nil
}
