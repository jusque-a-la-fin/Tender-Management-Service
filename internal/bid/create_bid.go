package bid

import (
	"fmt"
	"tendermanagement/internal/tender"
)

// CreateBid cоздает предложение для существующего тендера
func (repo *BidDBRepository) CreateBid(bci BidCreationInput) (*Bid, int, error) {
	valid, err := checkAuthor(repo.dtb, bci.AuthorId, bci.AuthorType)
	if !valid || err != nil {
		return nil, 401, err
	}

	valid, err = checkAuthorRights(repo.dtb, bci.AuthorId)
	if !valid || err != nil {
		return nil, 403, err
	}

	valid, err = tender.CheckTender(repo.dtb, bci.TenderId)
	if !valid || err != nil {
		return nil, 404, err
	}

	var status StatusEnum = "Created"
	var version int32 = 1
	query := `INSERT INTO bid (status, tender_id, author_type, author_id, current_version)
	         VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at;`

	var bidID, createdAt string
	err = repo.dtb.QueryRow(query, status, bci.TenderId, bci.AuthorType, bci.AuthorId, version).Scan(&bidID, &createdAt)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: создание тендера: %v", err)
	}

	query = `INSERT INTO bid_versions (version, name, description, bid_id)
	         VALUES ($1, $2, $3, $4);`

	_, err = repo.dtb.Exec(query, version, bci.Name, bci.Description, bidID)
	if err != nil {
		return nil, -1, fmt.Errorf("ошибка запроса к базе данных: создание версии тендера: %v", err)
	}

	bid := &Bid{
		ID:          bidID,
		Name:        bci.Name,
		Description: bci.Description,
		Status:      status,
		TenderId:    bci.TenderId,
		AuthorType:  bci.AuthorType,
		AuthorId:    bci.AuthorId,
		Version:     version,
		CreatedAt:   createdAt,
	}
	return bid, 200, nil
}
