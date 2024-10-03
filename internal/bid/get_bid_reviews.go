package bid

import (
	"database/sql"
	"fmt"
	"tendermanagement/internal"
)

// GetBidReviews просматривает отзывы на прошлые предложения
func (repo *BidDBRepository) GetBidReviews(bri BidReviewsInput) ([]BidReview, int, error) {
	valid, err := internal.CheckUser(repo.dtb, bri.RequesterUsername)
	if !valid || err != nil {
		return nil, 401, err
	}

	valid, err = checkAuthorName(repo.dtb, bri.AuthorUsername)
	if !valid || err != nil {
		return nil, 401, err
	}

	valid, err = CheckTender(repo.dtb, bri.TenderId)
	if !valid || err != nil {
		return nil, 404, err
	}

	var authorId string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1;", bri.AuthorUsername).Scan(&authorId)
	if err == sql.ErrNoRows {
		err = repo.dtb.QueryRow("SELECT id FROM organization WHERE name = $1;", bri.AuthorUsername).Scan(&authorId)
		if err != nil {
			return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
		}
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var hasRights bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM bid
		WHERE tender_id = $1 AND author_id = $2
	) AS has_rights;`

	err = repo.dtb.QueryRow(query, bri.TenderId, authorId).Scan(&hasRights)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных во время проверки прав доступа на отправку решения по предложению: %v", err)
	}

	if !hasRights {
		return nil, 403, nil
	}

	var requesterUserID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1;", bri.RequesterUsername).Scan(&requesterUserID)
	if err != nil {

		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var organizationID string
	err = repo.dtb.QueryRow("SELECT organization_id FROM tender WHERE id = $1", bri.TenderId).Scan(&organizationID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	query = `SELECT EXISTS (
                  SELECT 1
                  FROM tender
                  WHERE organization_id = $1 AND user_id = $2
              ) AS has_rights;`

	err = repo.dtb.QueryRow(query, organizationID, requesterUserID).Scan(&hasRights)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных во время проверки прав доступа на отправку решения по предложению: %v", err)
	}

	if !hasRights {
		return nil, 403, nil
	}

	reviews, err := getReviews(repo.dtb, authorId, bri.TenderId, bri.Limit, bri.Offset)
	if err != nil {
		return nil, -1, err
	}

	if len(reviews) == 0 {
		return nil, 404, nil
	}

	return reviews, 200, nil
}

func getReviews(dtb *sql.DB, authorId, tenderId string, limit, offset int32) ([]BidReview, error) {
	var bidId string
	err := dtb.QueryRow("SELECT id FROM bid WHERE author_id = $1 AND tender_id = $2;", authorId, tenderId).Scan(&bidId)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: извлечение id для предложения: %v", err)
	}

	var reviews []BidReview

	query := "SELECT id, description, created_at FROM bid_review WHERE bid_id = $1"
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	if offset > 0 {
		query = fmt.Sprintf("%s OFFSET %d", query, offset)
	}

	query = fmt.Sprintf("%s;", query)
	rows, err := dtb.Query(query, bidId)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: извлечение параметров отзывов: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		review := BidReview{}
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
