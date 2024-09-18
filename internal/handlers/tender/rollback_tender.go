package tender

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tendermanagement/internal/handlers"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// RollbackTender откатывает параметры тендера к указанной версии
func (hnd *TenderHandler) RollbackTender(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodPut {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	vars := mux.Vars(rqt)
	tenderID := vars["tenderId"]
	tenderIDLen := utf8.RuneCountInString(tenderID)
	if tenderIDLen == 0 || tenderIDLen > 100 {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	versionStr := vars["version"]
	versionInt, err := strconv.Atoi(versionStr)
	if err != nil || versionInt < 1 {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	version := int32(versionInt)

	username := rqt.URL.Query().Get("username")
	usernameLen := utf8.RuneCountInString(username)
	if usernameLen == 0 {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	tdr, code, err := hnd.TenderRepo.RollbackTender(version, tenderID, username)
	if err != nil {
		log.Println(err)
		return
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

	case 404:
		err := "Тендер или версия не найдены"
		errResp := handlers.RespondWithError(wrt, err, http.StatusNotFound)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
		}
		return
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(tdr)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
