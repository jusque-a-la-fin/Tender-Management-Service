package bid

import (
	"database/sql"
	"fmt"
	"tendermanagement/internal"
	"tendermanagement/internal/tender"
)

// SubmitBidDecision добавляет решение (одобрение или отклонение) по предложению
func (repo *BidDBRepository) AddBidDecisions(bsi BidSubmissionInput) (string, string, int, error) {
	valid, err := сheckBid(repo.dtb, bsi.BidID)
	if !valid || err != nil {
		return "", "", 404, err
	}

	valid, err = internal.CheckUser(repo.dtb, bsi.Username)
	if !valid || err != nil {
		return "", "", 401, err
	}

	var tenderID string
	query := `SELECT tender_id FROM bid WHERE id = $1;`
	err = repo.dtb.QueryRow(query, bsi.BidID).Scan(&tenderID)
	if err != nil {
		return "", "", -1, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", bsi.Username).Scan(&userID)
	if err != nil {
		return "", "", -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var organizationID string
	err = repo.dtb.QueryRow("SELECT organization_id FROM tender WHERE id = $1", tenderID).Scan(&organizationID)
	if err != nil {
		return "", "", -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	var hasRights bool
	query = `SELECT EXISTS (
                  SELECT 1
                  FROM organization_responsible
                  WHERE organization_id = $1 AND user_id = $2
              ) AS has_rights;`

	err = repo.dtb.QueryRow(query, organizationID, userID).Scan(&hasRights)
	if err != nil {
		return "", "", -1, fmt.Errorf("ошибка запроса к базе данных во время проверки прав доступа на отправку решения по предложению: %v", err)
	}

	if !hasRights {
		return "", "", 403, nil
	}

	query = `
	    INSERT INTO bid_decisions (decision, organization_id, user_id, bid_id)
	    VALUES ($1, $2, $3, $4);`

	result, err := repo.dtb.Exec(query, bsi.Decision, organizationID, userID, bsi.BidID)
	if err != nil {
		return "", "", -1, fmt.Errorf("ошибка запроса к базе данных: добавление нового решения по предложению: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", "", -1, fmt.Errorf("ошибка от метода `RowsAffected`, пакет sql: %v", err)
	}

	if rowsAffected == 0 {
		return "", "", -1, fmt.Errorf("ошибка запроса к базе данных: не добавилось решение по предложению")
	}
	return tenderID, organizationID, 200, nil
}

func (repo *BidDBRepository) MakeFinalDecision(bidID, tenderID, organizationID string) (*Bid, int, error) {
	rejectedCount, err := getDecisionCount(repo.dtb, Rejected, organizationID, bidID)
	if err != nil {
		return nil, -1, err
	}

	if rejectedCount == 0 {
		// согласование предложения
		approvedCount, err := getDecisionCount(repo.dtb, Approved, organizationID, bidID)
		if err != nil {
			return nil, -1, err
		}
		userCount, err := countResponsibleUsers(repo.dtb, organizationID)
		if err != nil {
			return nil, -1, err
		}
		if approvedCount >= findMin(3, userCount) {
			err = CloseTender(repo.dtb, tenderID)
			if err != nil {
				return nil, -1, err
			}
		}
	}

	bid, err := GetBid(repo.dtb, bidID)
	if err != nil {
		return nil, -1, err
	}
	return bid, 200, nil
}

// CancelBid отменяет предложение
// func CancelBid(dtb *sql.DB, bidID string) error {
// 	query := `
// 	UPDATE bid
// 	SET status = $1
// 	WHERE id = $2;`

// 	result, err := dtb.Exec(query, Canceled, bidID)
// 	if err != nil {
// 		return fmt.Errorf("ошибка запроса к базе данных: отклонение предложения: %v", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("ошибка от метода `RowsAffected`, пакет sql: %v", err)
// 	}

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("ошибка запроса к базе данных: не отклонилось предложение")
// 	}
// 	return nil
// }

// GetDecisionCount считает количество решений для одного из типов: 'Approved', 'Rejected'
func getDecisionCount(dtb *sql.DB, decisionType DecisionEnum, organizationID, bidID string) (int, error) {
	var decisionCount int
	query := `
        SELECT COUNT(*) AS decision_count
        FROM bid_decisions
        WHERE decision = $1
          AND organization_id = $2
          AND bid_id = $3;`

	err := dtb.QueryRow(query, decisionType, organizationID, bidID).Scan(&decisionCount)
	if err != nil {
		return -1, fmt.Errorf("ошибка запроса к базе данных: подсчет решений `Approved` или `Rejected` для данного предложения")
	}
	return decisionCount, nil
}

// countResponsibleUsers считает количество ответственных в организации
func countResponsibleUsers(dtb *sql.DB, organizationID string) (int, error) {
	var userCount int
	query := `
		SELECT COUNT(DISTINCT user_id) AS user_count
		FROM organization_responsible
		WHERE organization_id = $1;`

	err := dtb.QueryRow(query, organizationID).Scan(&userCount)
	if err != nil {
		return -1, fmt.Errorf("ошибка запроса к базе данных: подсчет количества ответственных в организации")
	}

	return userCount, nil
}

func findMin(digit1, digit2 int) int {
	if digit1 < digit2 {
		return digit1
	}
	return digit2
}

// CloseTender закрывает тендер
func CloseTender(dtb *sql.DB, tenderID string) error {
	query := `
	UPDATE tender
	SET status = $1
	WHERE id = $2;`

	result, err := dtb.Exec(query, tender.Closed, tenderID)
	if err != nil {
		return fmt.Errorf("ошибка запроса к базе данных: закрытие тендера: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка от метода `RowsAffected`, пакет sql: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ошибка запроса к базе данных: не закрылся тендер")
	}
	return nil
}
