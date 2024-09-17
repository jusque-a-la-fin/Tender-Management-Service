package bid

import (
	"fmt"
)

// GetBidStatus получает статус предложения по его уникальному идентификатору
func (repo *BidDBRepository) GetBidStatus(bidID, username string) (StatusEnum, int, error) {
	valid, err := checkUsername(repo.dtb, username)
	if !valid || err != nil {
		return "", 401, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return "", -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	valid, err = checkEditionRights(repo.dtb, bidID, userID)
	if !valid || err != nil {
		return "", 403, err
	}

	valid, err = сheckBid(repo.dtb, bidID)
	if !valid || err != nil {
		return "", 404, err
	}

	var status StatusEnum
	err = repo.dtb.QueryRow(`SELECT status FROM bid WHERE id = $1 AND author_id = $2`, bidID, userID).Scan(&status)
	if err != nil {
		return "", -1, fmt.Errorf("ошибка запроса к базе данных: извлечение status предложения: %v", err)
	}

	return status, 200, nil
}
