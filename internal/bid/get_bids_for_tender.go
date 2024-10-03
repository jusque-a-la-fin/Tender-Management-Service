package bid

import (
	"database/sql"
	"fmt"
	"sort"
	"tendermanagement/internal/tender"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// GetBidsForTender получает предложения, связанные с указанным тендером
func (repo *BidDBRepository) GetBidsForTender(gbi GetBidsInput) ([]*Bid, int, error) {
	valid, err := checkAuthorName(repo.dtb, gbi.Username)
	if !valid || err != nil {
		return nil, 401, err
	}

	valid, err = tender.CheckTender(repo.dtb, gbi.TenderId)
	if !valid || err != nil {
		return nil, 404, err
	}

	var userId string
	err = repo.dtb.QueryRow("SELECT id FROM employee WHERE username = $1;", gbi.Username).Scan(&userId)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: извлечение id для username: %v", err)
	}

	valid, err = tender.CheckRights(repo.dtb, gbi.TenderId, userId)
	if !valid || err != nil {
		return nil, 403, err
	}

	nbids, err := getBids(repo.dtb, gbi.TenderId, gbi.Limit, gbi.Offset)
	if err != nil {
		return nil, -1, err
	}

	if len(nbids) == 0 {
		return nil, 404, nil
	}

	clr := collate.New(language.Russian)
	sort.Slice(nbids, func(i, j int) bool {
		return clr.CompareString(nbids[i].Name, nbids[j].Name) < 0
	})

	return nbids, 200, nil
}

// getBids получает предложения
func getBids(dtb *sql.DB, tenderID string, limit, offset int32) ([]*Bid, error) {
	var nbids []*Bid
	query := `
        SELECT 
            b.id,
            b.status,
            b.author_type,
            b.author_id,
            b.current_version,
            b.created_at,
            bv.name,
            bv.description
        FROM 
            bid b
        JOIN 
            bid_versions bv ON b.id = bv.bid_id AND b.current_version = bv.version
        WHERE 
            b.tender_id = $1`

	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	if offset > 0 {
		query = fmt.Sprintf("%s OFFSET %d", query, offset)
	}

	query = fmt.Sprintf("%s;", query)

	rows, err := dtb.Query(query, tenderID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: извлечение предложений, связанных с указанным тендером: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		nbid := Bid{}
		err := rows.Scan(&nbid.ID, &nbid.Status, &nbid.AuthorType, &nbid.AuthorID,
			&nbid.Version, &nbid.CreatedAt, &nbid.Name, &nbid.Description)
		if err != nil {
			return nil, fmt.Errorf("ошибка от метода `Scan`, пакет sql: %v", err)
		}
		nbids = append(nbids, &nbid)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время итерирования по строкам, возвращенным запросом: %v", err)
	}
	return nbids, nil
}
