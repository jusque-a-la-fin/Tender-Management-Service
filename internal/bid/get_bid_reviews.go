package bid

import (
	"database/sql"
	"fmt"
	"tendermanagement/internal/tender"
)

// GetBidReviews просматривает отзывы на прошлые предложения
func (repo *BidDBRepository) GetBidReviews(bri BidReviewsInput) ([]*BidReview, int, error) {
	valid, err := tender.CheckTender(repo.dtb, bri.TenderId)
	if !valid || err != nil {
		return nil, 404, err
	}

	valid, err = checkUsername(repo.dtb, bri.RequesterUsername)
	if !valid || err != nil {
		return nil, 401, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", bri.RequesterUsername).Scan(&userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var organizationID string
	err = repo.dtb.QueryRow("SELECT organization_id FROM tender WHERE id = $1", bri.TenderId).Scan(&organizationID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var hasRights bool
	query := `SELECT EXISTS (
                  SELECT 1
                  FROM organization_responsible
                  WHERE organization_id = $1 AND user_id = $2
              ) AS has_rights;`

	err = repo.dtb.QueryRow(query, organizationID, userID).Scan(&hasRights)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных во время проверки прав доступа на отправку решения по предложению: %v", err)
	}

	if !hasRights {
		return nil, 403, nil
	}

	valid, err = checkUsername(repo.dtb, bri.AuthorUsername)
	if !valid || err != nil {
		return nil, 401, err
	}

	reviews, err := getReviews(repo.dtb, bri.TenderId, bri.AuthorUsername)
	if err != nil {
		return nil, -1, err
	}

	if len(reviews) == 0 {
		return nil, 404, nil
	}

	if bri.Offset != tender.NoValue && bri.EndIndex != tender.NoValue {
		return reviews[bri.Offset:bri.EndIndex], 200, nil
	}

	return reviews, 200, nil
}

func getReviews(dtb *sql.DB, tenderId, authorUsername string) ([]*BidReview, error) {
	var bidID string
	query := `SELECT id FROM bid WHERE tender_id = $1`
	err := dtb.QueryRow(query, tenderId).Scan(&bidID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: извлечение id предложения: %v", err)
	}

	var userID string
	err = dtb.QueryRow("SELECT id FROM employee WHERE username = $1", authorUsername).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var reviews []*BidReview

	query = `
        SELECT brw.id, brw.description, brw.created_at
        FROM bid_review brw
        JOIN bid_versions bdv ON bdv.bid_id = brw.bid_id
        WHERE bdv.bid_id = $1 AND brw.user_id = $2;
    `

	rows, err := dtb.Query(query, bidID, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: извлечение параметров отзывов: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		review := &BidReview{}
		if err := rows.Scan(&review.ID, &review.Description, &review.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка от метода `Scan`, пакет sql: %v", err)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время итерирования по строкам, возвращенным запросом: %v", err)
	}

	return reviews, nil
}
