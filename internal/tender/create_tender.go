package tender

import (
	"database/sql"
	"fmt"
	"tendermanagement/internal"
)

// CreateTender создает новый тендер с заданными параметрами
func (repo *TenderDBRepository) CreateTender(tci TenderCreationInput, creatorUsername, organizationID string) (*TenderCreationOutput, int, error) {
	valid, err := internal.CheckUser(repo.dtb, creatorUsername)
	if !valid || err != nil {
		return nil, 401, err
	}

	query := `
            SELECT org.user_id
            FROM organization_responsible org
            JOIN employee emp ON org.user_id = emp.id
            WHERE emp.username = $1 AND org.organization_id = $2;`

	var userID string
	err = repo.dtb.QueryRow(query, creatorUsername, organizationID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 403, nil
		}
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: проверка прав доступа: %v", err)
	}

	var status StatusEnum = "Created"
	var version int32 = 1
	query = `INSERT INTO tender (status, current_version, created_at, user_id, organization_id)
	         VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	var tenderID string
	err = repo.dtb.QueryRow(query, status, version, tci.CreatedAt, userID, organizationID).Scan(&tenderID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: создание тендера: %v", err)
	}

	query = `INSERT INTO tender_versions (version, name, description, service_type, tender_id)
	         VALUES ($1, $2, $3, $4, $5);`

	_, err = repo.dtb.Exec(query, version, tci.Name, tci.Description, tci.ServiceType, tenderID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: создание версии тендера: %v", err)
	}

	tou := TenderCreationOutput{
		ID:      tenderID,
		Status:  status,
		Version: version,
	}

	return &tou, 200, nil
}
