package bid

import (
	"fmt"
	"tendermanagement/internal"
)

// SubmitBidFeedback отправляет отзыв по предложению
func (repo *BidDBRepository) SubmitBidFeedback(bfi BidFeedbackInput) (*Bid, int, error) {
	valid, err := сheckBid(repo.dtb, bfi.BidID)
	if !valid || err != nil {
		return nil, 404, err
	}

	valid, err = internal.CheckUser(repo.dtb, bfi.Username)
	if !valid || err != nil {
		return nil, 401, err
	}

	var tenderID string
	query := `SELECT tender_id FROM bid WHERE id = $1;`
	err = repo.dtb.QueryRow(query, bfi.BidID).Scan(&tenderID)
	if err != nil {
		return nil, -1, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1;", bfi.Username).Scan(&userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var organizationID string
	err = repo.dtb.QueryRow("SELECT organization_id FROM tender WHERE id = $1;", tenderID).Scan(&organizationID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var hasRights bool
	query = `SELECT EXISTS (
                  SELECT 1
                  FROM tender
                  WHERE organization_id = $1 AND user_id = $2
              ) AS has_rights;`

	err = repo.dtb.QueryRow(query, organizationID, userID).Scan(&hasRights)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных во время проверки прав доступа на отправку решения по предложению: %v", err)
	}

	if !hasRights {
		return nil, 403, nil
	}

	query = `
	    INSERT INTO bid_review (description, user_id, bid_id)
	    VALUES ($1, $2, $3);`

	result, err := repo.dtb.Exec(query, bfi.BidFeedback, userID, bfi.BidID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: добавление нового отзыва по предложению: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка от метода `RowsAffected`, пакет sql: %v", err)
	}

	if rowsAffected == 0 {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: не добавился отзыв по предложению")
	}

	bid, err := GetBid(repo.dtb, bfi.BidID)
	if err != nil {
		return nil, -1, err
	}
	return bid, 200, nil
}
