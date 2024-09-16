package tender

import (
	"encoding/json"
	"log"
	"net/http"
	"tendermanagement/internal/handlers"
	tnd "tendermanagement/internal/tender"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// EditTender изменяет параметры существующего тендера
func (hnd *TenderHandler) EditTender(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodPatch {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	vars := mux.Vars(rqt)
	tenderID := vars["tenderId"]
	tenderIDLen := utf8.RuneCountInString(tenderID)
	if tenderIDLen == 0 || tenderIDLen > 100 {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	username := rqt.URL.Query().Get("username")
	usernameLen := utf8.RuneCountInString(username)
	if usernameLen == 0 {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	var trq tnd.TenderEditionInput
	err := json.NewDecoder(rqt.Body).Decode(&trq)
	if err != nil {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	tdNameLen := utf8.RuneCountInString(trq.Name)
	tdDescLen := utf8.RuneCountInString(trq.Description)

	switch {
	case tdNameLen > 100:
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}

	case tdDescLen > 500:
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	fail := tnd.CheckServiceType(string(trq.ServiceType))
	if fail {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
	}

	tdr, code, err := hnd.TenderRepo.EditTender(trq, tenderID, username)
	if err != nil {
		log.Println(err)
	}

	switch code {
	case 401:
		err := "Пользователь не существует или некорректен"
		errResp := handlers.RespondWithError(wrt, err, http.StatusUnauthorized)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
			return
		}

	case 403:
		err := "Недостаточно прав для выполнения действия"
		errResp := handlers.RespondWithError(wrt, err, http.StatusForbidden)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
			return
		}

	case 404:
		err := "Предложение не найдено"
		errResp := handlers.RespondWithError(wrt, err, http.StatusNotFound)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
			return
		}
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(tdr)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
