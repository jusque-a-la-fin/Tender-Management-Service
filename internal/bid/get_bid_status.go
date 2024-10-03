package bid

import (
	"database/sql"
	"fmt"
)

// GetBidStatus получает статус предложения по его уникальному идентификатору
func (repo *BidDBRepository) GetBidStatus(bidID, authorName string) (StatusEnum, int, error) {
	valid, err := checkAuthorName(repo.dtb, authorName)
	if !valid || err != nil {
		return "", 401, err
	}

	valid, err = сheckBid(repo.dtb, bidID)
	if !valid || err != nil {
		return "", 404, err
	}

	var authorId string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1;", authorName).Scan(&authorId)
	if err == sql.ErrNoRows {
		err = repo.dtb.QueryRow("SELECT id FROM organization WHERE name = $1;", authorName).Scan(&authorId)
		if err != nil {
			return "", -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для authorName: %v", err)
		}
	}

	if err != nil && err != sql.ErrNoRows {
		return "", -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для authorName: %v", err)
	}

	valid, err = checkEditionRights(repo.dtb, bidID, authorId)
	if !valid || err != nil {
		return "", 403, err
	}

	var status StatusEnum
	err = repo.dtb.QueryRow(`SELECT status FROM bid WHERE id = $1 AND author_id = $2;`, bidID, authorId).Scan(&status)
	if err != nil {
		return "", -1, fmt.Errorf("ошибка запроса к базе данных: извлечение status предложения: %v", err)
	}

	return status, 200, nil
}
