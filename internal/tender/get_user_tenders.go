package tender

import (
	"fmt"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// GetUserTenders получает список тендеров текущего пользователя
func (repo *TenderDBRepository) GetUserTenders(startIndex, endIndex int32, username string) ([]Tender, int, error) {
	valid, err := CheckUser(repo.dtb, username)
	if !valid || err != nil {
		return nil, 401, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	tenders := make([]Tender, 0)
	query := `
    SELECT tdr.id, tdr.status, tdr.current_version, tdr.created_at, 
           tdr.organization_id, tvs.name, tvs.description, tvs.service_type 
    FROM tender tdr
    JOIN tender_versions tvs ON tdr.id = tvs.tender_id AND tdr.current_version = tvs.version
    WHERE tdr.user_id = $1`

	rows, err := repo.dtb.Query(query, userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение параметров тендеров: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return tenders, 200, nil
	}

	for rows.Next() {
		tnd := Tender{}
		err := rows.Scan(&tnd.ID, &tnd.Status, &tnd.Version, &tnd.CreatedAt,
			&tnd.OrganizationID, &tnd.Name, &tnd.Description, &tnd.ServiceType)
		if err != nil {
			return nil, -1, fmt.Errorf("ошибка от метода `Scan`, пакет sql: %v", err)
		}
		tenders = append(tenders, tnd)
	}

	if err = rows.Err(); err != nil {
		return nil, -1, fmt.Errorf("ошибка во время итерирования по строкам, возвращенным запросом: %v", err)
	}

	clr := collate.New(language.Russian)
	sort.Slice(tenders, func(i, j int) bool {
		return clr.CompareString(tenders[i].Name, tenders[j].Name) < 0
	})

	if startIndex != NoValue && endIndex != NoValue {
		return tenders[startIndex:endIndex], 200, nil
	}

	return tenders, 200, nil
}
