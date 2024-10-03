package tender

import (
	"encoding/json"
	"log"
	"net/http"
	"tendermanagement/internal/handlers"
	tnd "tendermanagement/internal/tender"
	"time"
	"unicode/utf8"
)

// CreateTender создает новый тендер с заданными параметрами
func (hnd *TenderHandler) CreateTender(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodPost {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	var trq TenderCreationReq
	err := json.NewDecoder(rqt.Body).Decode(&trq)
	if err != nil {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	tdName := utf8.RuneCountInString(trq.Name)
	tdDesc := utf8.RuneCountInString(trq.Description)
	tdOrgID := utf8.RuneCountInString(trq.OrganizationID)

	switch {
	case tdName == 0 || tdName > 100:
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return

	case tdDesc == 0 || tdDesc > 500:
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return

	case tnd.CheckServiceType(trq.ServiceType):
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return

	case tdOrgID == 0 || tdOrgID > 100:
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return

	case trq.CreatorUsername == "":
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	tci := tnd.TenderCreationInput{
		Name:        trq.Name,
		Description: trq.Description,
		CreatedAt:   time.Now().Format(time.RFC3339),
		ServiceType: tnd.ServiceTypeEnum(trq.ServiceType),
	}

	tcu, code, err := hnd.TenderRepo.CreateTender(tci, trq.CreatorUsername, trq.OrganizationID)
	if err != nil {
		log.Println(err)
		if !handlers.CheckCode(code) {
			return
		}
	}

	switch code {
	case 401:
		err := "Пользователь не существует или некорректен"
		errResp := handlers.RespondWithError(wrt, err, http.StatusUnauthorized)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
		}
		return

	case 403:
		err := "Недостаточно прав для выполнения действия"
		errResp := handlers.RespondWithError(wrt, err, http.StatusForbidden)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
		}
		return
	}

	tender := tnd.GetCreatedTender(tci, *tcu, trq.OrganizationID)
	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(tender)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
