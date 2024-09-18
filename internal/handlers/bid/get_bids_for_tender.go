package bid

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	bds "tendermanagement/internal/bid"
	"tendermanagement/internal/handlers"
	"tendermanagement/internal/tender"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// GetBidsForTender получает предложения, связанные с указанным тендером
func (hnd *BidHandler) GetBidsForTender(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodGet {
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

	username := rqt.URL.Query().Get("username")
	usernameLen := utf8.RuneCountInString(username)
	if usernameLen == 0 {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
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
			return
		}

		limit = int32(limitInt)
		if limit < 0 || limit > 50 {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			}
			return
		}
	}

	offset := tender.NoValue
	offsetStr := rqt.URL.Query().Get("offset")
	if offsetStr != "" {
		offsetInt, err := strconv.Atoi(offsetStr)
		if err != nil {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			}
			return
		}

		offset = int32(offsetInt)
		if offset < 0 {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			}
			return
		}
	}
	endIndex := offset + limit

	gbi := bds.GetBidsInput{
		TenderId: tenderID,
		Username: username,
		Offset:   offset,
		EndIndex: endIndex,
	}

	nbds, code, err := hnd.BidRepo.GetBidsForTender(gbi)
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
		err := "Тендер или предложение не найдено"
		errResp := handlers.RespondWithError(wrt, err, http.StatusForbidden)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
		}
		return
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(nbds)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
