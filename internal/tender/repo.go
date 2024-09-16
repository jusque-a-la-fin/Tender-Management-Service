package tender

import (
	"database/sql"
)

type TenderRepo interface {
	GetTenders(startIndex, endIndex int32, serviceTypes []ServiceTypeEnum) ([]Tender, error)
	GetUserTenders(startIndex, endIndex int32, username string) ([]Tender, int, error)
	CreateTender(tci TenderCreationInput, creatorUsername, organizationID string) (*TenderCreationOutput, int, error)
	GetTenderStatus(username, tenderID string) (string, int, error)
	UpdateTenderStatus(tenderID, newStatus, username string) (*Tender, int, error)
	EditTender(tei TenderEditionInput, tenderID, username string) (*Tender, int, error)
	RollbackTender(version int32, tenderID, username string) (*Tender, int, error)
}

type TenderDBRepository struct {
	dtb *sql.DB
}

func NewDBRepo(sdb *sql.DB) *TenderDBRepository {
	return &TenderDBRepository{dtb: sdb}
}
