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

// GetUserTenders получает список тендеров текущего пользователя
func (hnd *TenderHandler) GetUserTenders(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodGet {
		wrt.WriteHeader(http.StatusBadRequest)
		return
	}

	var limit int32 = 0
	limitStr := rqt.URL.Query().Get("limit")
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			wrt.WriteHeader(http.StatusBadRequest)
			return
		}

		limit = int32(limitInt)
		if limit < 0 || limit > 50 {
			wrt.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	offset := tnd.NoValue
	offsetStr := rqt.URL.Query().Get("offset")
	if offsetStr != "" {
		offsetInt, err := strconv.Atoi(offsetStr)
		if err != nil {
			wrt.WriteHeader(http.StatusBadRequest)
			return
		}

		offset = int32(offsetInt)
		if offset < 0 {
			wrt.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	username := rqt.URL.Query().Get("username")
	usernameLen := utf8.RuneCountInString(username)
	if usernameLen == 0 {
		wrt.WriteHeader(http.StatusBadRequest)
		return
	}

	endIndex := offset + limit
	tenders, code, err := hnd.TenderRepo.GetUserTenders(offset, endIndex, username)
	if err != nil {
		log.Println(err)
		return
	}

	if code == 401 {
		err := "Пользователь не существует или некорректен"
		errResp := handlers.RespondWithError(wrt, err, http.StatusUnauthorized)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
		}
		return
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(tenders)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
