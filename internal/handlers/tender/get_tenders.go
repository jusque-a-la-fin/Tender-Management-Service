package tender

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tendermanagement/internal/handlers"
	tnd "tendermanagement/internal/tender"
	"unicode/utf8"
)

// GetTenders получает список тендеров с возможностью фильтрации по типу услуг
func (hnd *TenderHandler) GetTenders(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodGet {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
	}

	var limit int32 = 0
	limitStr := rqt.URL.Query().Get("limit")
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			}
		}

		limit = int32(limitInt)
		if limit < 0 || limit > 50 {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			}
		}
	}

	offset := tnd.NoValue
	offsetStr := rqt.URL.Query().Get("offset")
	if offsetStr != "" {
		offsetInt, err := strconv.Atoi(offsetStr)
		if err != nil {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			}
		}

		offset = int32(offsetInt)
		if offset < 0 {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			}
		}
	}

	endIndex := offset + limit

	serviceTypesRaw := rqt.URL.Query()["service_type"]
	var serviceTypes []tnd.ServiceTypeEnum
	for _, service := range serviceTypesRaw {
		checkServiceTypes(wrt, service)
		serviceTypes = append(serviceTypes, tnd.ServiceTypeEnum(service))
	}

	tenders, err := hnd.TenderRepo.GetTenders(offset, endIndex, []tnd.ServiceTypeEnum(serviceTypes))
	if err != nil {
		log.Println(err)
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(tenders)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}

func checkServiceTypes(wrt http.ResponseWriter, serviceType string) {
	serviceTypeLen := utf8.RuneCountInString(serviceType)
	if serviceTypeLen == 0 {
		wrt.WriteHeader(http.StatusBadRequest)
		return
	}

	fail := tnd.CheckServiceType(serviceType)
	if fail {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
	}
}
