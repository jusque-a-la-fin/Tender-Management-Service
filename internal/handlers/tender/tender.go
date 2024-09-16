package tender

import (
	tns "tendermanagement/internal/tender"
)

type TenderHandler struct {
	TenderRepo tns.TenderRepo
}

// TenderCreationReq - тело запроса для создания тендера
type TenderCreationReq struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	OrganizationID  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}
