package tender

import (
	"database/sql"
	"fmt"
)

// NoValue - отсутствие значения
const NoValue int32 = -1

// Tender - тендер в теле ответа
type Tender struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	ServiceType    ServiceTypeEnum `json:"serviceType"`
	Status         StatusEnum      `json:"status"`
	OrganizationID string          `json:"organizationId"`
	Version        int32           `json:"version"`
	CreatedAt      string          `json:"createdAt"`
}

type ServiceTypeEnum string

const (
	Construction ServiceTypeEnum = "Construction"
	Delivery     ServiceTypeEnum = "Delivery"
	Manufacture  ServiceTypeEnum = "Manufacture"
)

type StatusEnum string

const (
	Created   StatusEnum = "Created"
	Published StatusEnum = "Published"
	Closed    StatusEnum = "Closed"
)

// TenderInput - параметры, установленные пользователем (за исключением CreatedAt) для создания тендера
type TenderCreationInput struct {
	Name        string
	Description string
	ServiceType ServiceTypeEnum
	CreatedAt   string
}

// TenderOutput - параметры, установленные сервером при создании тендера
type TenderCreationOutput struct {
	ID      string
	Status  StatusEnum
	Version int32
}

func GetCreatedTender(tci TenderCreationInput, tcu TenderCreationOutput, organizationID string) Tender {
	tender := Tender{}
	tender.ID = tcu.ID
	tender.Name = tci.Name
	tender.Description = tci.Description
	tender.ServiceType = tci.ServiceType
	tender.Status = tcu.Status
	tender.OrganizationID = organizationID
	tender.Version = tcu.Version
	tender.CreatedAt = tci.CreatedAt
	return tender
}

// TenderInput - параметры, установленные пользователем (за исключением CreatedAt) для изменения существующего тендера
type TenderEditionInput struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ServiceType ServiceTypeEnum `json:"serviceType"`
}

func GetTender(dtb *sql.DB, tenderID string) (*Tender, error) {
	tnd := Tender{
		ID: tenderID,
	}

	query := `
		SELECT tdr.status, tdr.current_version, tdr.created_at, tdr.organization_id, 
			   tvs.name, tvs.description, tvs.service_type
		FROM tender tdr
		JOIN tender_versions tvs ON tdr.id = tvs.tender_id AND tdr.current_version = tvs.version
		WHERE tdr.id = $1
	`

	err := dtb.QueryRow(query, tenderID).Scan(&tnd.Status, &tnd.Version, &tnd.CreatedAt,
		&tnd.OrganizationID, &tnd.Name,
		&tnd.Description, &tnd.ServiceType)

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: извлечение параметров тендера: %v", err)
	}

	return &tnd, nil
}

func CheckStatus(status string) bool {
	fail := true
	switch status {
	case "Created":
		fail = false
	case "Published":
		fail = false
	case "Closed":
		fail = false
	}
	return fail
}

func CheckServiceType(serviceType string) bool {
	fail := true
	switch serviceType {
	case "Construction":
		fail = false
	case "Delivery":
		fail = false
	case "Manufacture":
		fail = false
	}
	return fail
}
