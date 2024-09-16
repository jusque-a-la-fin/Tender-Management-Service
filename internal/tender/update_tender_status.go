package tender

import (
	"fmt"
)

// UpdateTenderStatus изменяет статус тендера по его идентификатору
func (repo *TenderDBRepository) UpdateTenderStatus(tenderID, newStatus, username string) (*Tender, int, error) {
	valid, err := CheckUser(repo.dtb, username)
	if !valid || err != nil {
		return nil, 401, err
	}

	valid, err = CheckTender(repo.dtb, tenderID)
	if !valid || err != nil {
		return nil, 404, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	query := `
        UPDATE tender
        SET status = $1
        WHERE id = $2 AND user_id = $3;
    `

	result, err := repo.dtb.Exec(query, newStatus, tenderID, userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: обновление статуса тендера: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка от метода `RowsAffected`, пакет sql: %v", err)
	}

	if rowsAffected == 0 {
		return nil, 403, nil
	}

	tnd := Tender{
		ID:     tenderID,
		Status: StatusEnum(newStatus),
	}

	query = `
		SELECT tdr.current_version, tdr.created_at, tdr.organization_id, 
			   tvs.name, tvs.description, tvs.service_type
		FROM tender tdr
		JOIN tender_versions tvs ON tdr.id = tvs.tender_id AND tdr.current_version = tvs.version
		WHERE tdr.id = $1`

	err = repo.dtb.QueryRow(query, tenderID).Scan(&tnd.Version, &tnd.CreatedAt,
		&tnd.OrganizationID, &tnd.Name,
		&tnd.Description, &tnd.ServiceType)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение данных тендера: %v", err)
	}

	return &tnd, 200, nil
}
