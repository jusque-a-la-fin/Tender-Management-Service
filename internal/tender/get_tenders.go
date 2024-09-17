package tender

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// GetTenders получает список тендеров с возможностью фильтрации по типу услуг
func (repo *TenderDBRepository) GetTenders(startIndex, endIndex int32, serviceTypes []ServiceTypeEnum) ([]Tender, error) {
	tenders := make([]Tender, 0)
	query := `
    SELECT tdr.id, tdr.status, tdr.current_version, tdr.created_at, 
           tdr.organization_id, tvs.name, tvs.description, tvs.service_type 
    FROM tender tdr
    JOIN tender_versions tvs ON tdr.id = tvs.tender_id AND tdr.current_version = tvs.version`

	var args []interface{}

	if len(serviceTypes) > 0 {
		placeholders := make([]string, len(serviceTypes))
		args = make([]interface{}, len(serviceTypes))

		for i, serviceType := range serviceTypes {
			placeholders[i] = "?"
			args[i] = serviceType
		}
		query += " WHERE tvs.service_type IN (" + strings.Join(placeholders, ", ") + ")"
	}

	stmt, err := repo.dtb.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка подготовки запроса: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: извлечение параметров тендеров: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		tnd := Tender{}
		err := rows.Scan(&tnd.ID, &tnd.Status, &tnd.Version, &tnd.CreatedAt,
			&tnd.OrganizationID, &tnd.Name, &tnd.Description, &tnd.ServiceType)
		if err != nil {
			return nil, fmt.Errorf("ошибка от метода `Scan`, пакет sql: %v", err)
		}
		tenders = append(tenders, tnd)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время итерирования по строкам, возвращенным запросом: %v", err)
	}

	clr := collate.New(language.Russian)
	sort.Slice(tenders, func(i, j int) bool {
		return clr.CompareString(tenders[i].Name, tenders[j].Name) < 0
	})

	if startIndex != NoValue && endIndex != NoValue {
		return tenders[startIndex:endIndex], nil
	}
	return tenders, nil
}
