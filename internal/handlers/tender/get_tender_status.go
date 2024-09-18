package tender

import (
	"encoding/json"
	"log"
	"net/http"
	"tendermanagement/internal/handlers"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// GetTenderStatus получает статус тендера по его уникальному идентификатору
func (hnd *TenderHandler) GetTenderStatus(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodGet {
		wrt.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(rqt)
	tenderID := vars["tenderId"]
	tenderIDLen := utf8.RuneCountInString(tenderID)
	if tenderIDLen == 0 || tenderIDLen > 100 {
		wrt.WriteHeader(http.StatusBadRequest)
		return
	}

	username := rqt.URL.Query().Get("username")
	usernameLen := utf8.RuneCountInString(username)
	if usernameLen == 0 {
		wrt.WriteHeader(http.StatusBadRequest)
		return
	}

	status, code, err := hnd.TenderRepo.GetTenderStatus(username, tenderID)
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
		err := "Тендер не найден"
		errResp := handlers.RespondWithError(wrt, err, http.StatusNotFound)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
		}
		return
	}

	tenderStatus := struct {
		TenderStatus string `json:"tenderStatus"`
	}{
		TenderStatus: status,
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(tenderStatus)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
