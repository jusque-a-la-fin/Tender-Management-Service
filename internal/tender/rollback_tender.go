package tender

import (
	"database/sql"
	"fmt"
)

// RollbackTender откатывает параметры тендера к указанной версии
func (repo *TenderDBRepository) RollbackTender(version int32, tenderID, username string) (*Tender, int, error) {
	valid, err := CheckUser(repo.dtb, username)
	if !valid || err != nil {
		return nil, 401, err
	}

	valid, err = CheckTenderAndVersion(repo.dtb, version, tenderID)
	if !valid || err != nil {
		return nil, 404, err
	}

	var userID string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	query := `
	         SELECT EXISTS 
	            (SELECT 1 
				FROM tender
                WHERE id = $1 AND user_id = $2) 
				AS result`

	var hasRights bool
	err = repo.dtb.QueryRow(query, tenderID, userID).Scan(&hasRights)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: проверка прав доступа: %v", err)
	}

	if !hasRights {
		return nil, 403, nil
	}

	err = swapParams(repo.dtb, version, tenderID)
	if err != nil {
		return nil, -1, err
	}

	version++
	query = `
		     UPDATE tender
		     SET current_version = $1
		     WHERE id = $2;`

	_, err = repo.dtb.Exec(query, version, tenderID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: обновление текущей версии тендера: %v", err)
	}

	tnd, err := GetTender(repo.dtb, tenderID)
	if err != nil {
		return nil, -1, err
	}

	return tnd, 200, nil
}

// swapParams меняет местами параметры тендера для двух версий:
// 1я версия - версия, к которой нужно откатить тендер
// 2я версия (если она есть в базе данных) - версия, которая больше 1ой версии на 1 и
// которая будет скрыта инкрементированной 1ой версией
// rollbackVersion - первая версия до инкремента
// swapParams предотвращает скрытие
func swapParams(dtb *sql.DB, rollbackVersion int32, tenderID string) error {
	// проверка, существует ли версия на 1 больше, чем версия, к которой нужно откатить тендер
	// swappedVersion - 2я версия
	var swappedVersion int32 = rollbackVersion + 1
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM tender_versions
            WHERE tender_id = $1 AND version = $2
        );
    `

	err := dtb.QueryRow(query, tenderID, rollbackVersion).Scan(&exists)
	if err != nil {
		return fmt.Errorf(`ошибка запроса к базе данных: проверка существования версии тендера, 
		которая больше на 1, чем та версия(до инкремента), к которой нужно откатить тендер: %v`, err)
	}

	if !exists {
		return nil
	}

	rows, err := dtb.Query(`
        SELECT id, name, description, service_type
        FROM tender_versions
        WHERE tender_id = $1 AND version IN ($2, $3)`, tenderID, rollbackVersion, swappedVersion)
	if err != nil {
		return fmt.Errorf("ошибка запроса к базе данных: извлечение параметров двух версий тендера: %v", err)
	}
	defer rows.Close()

	versions := make(map[int32]struct {
		ID          int
		Name        string
		Description string
		ServiceType string
	})

	for rows.Next() {
		var id int
		var name, description, serviceType string
		var version int32

		if err := rows.Scan(&id, &name, &description, &serviceType); err != nil {
			return fmt.Errorf("ошибка от метода `Scan`, пакет sql: %v", err)
		}

		versions[version] = struct {
			ID          int
			Name        string
			Description string
			ServiceType string
		}{ID: id, Name: name, Description: description, ServiceType: serviceType}
	}

	if verRB, okRB := versions[rollbackVersion]; okRB {
		if verSW, okSW := versions[swappedVersion]; okSW {
			_, err := dtb.Exec(`
                UPDATE tender_versions
                SET name = CASE
                    WHEN version = $1 THEN $2
                    WHEN version = $3 THEN $4
                END,
                description = CASE
                    WHEN version = $1 THEN $5
                    WHEN version = $3 THEN $6
                END,
                service_type = CASE
                    WHEN version = $1 THEN $7
                    WHEN version = $3 THEN $8
                END
                WHERE tender_id = $9 AND version IN ($1, $3)`,
				rollbackVersion, verSW.Name, swappedVersion, verRB.Name,
				verSW.Description, verRB.Description,
				verSW.ServiceType, verRB.ServiceType,
				tenderID)

			if err != nil {
				return fmt.Errorf("ошибка запроса к базе данных: обмен параметрами между двумя версиями тендера: %v", err)
			}
		}
	}
	return nil
}
