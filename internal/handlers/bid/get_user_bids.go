package bid

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tendermanagement/internal/handlers"
	"tendermanagement/internal/tender"
	"unicode/utf8"
)

// GetUserBids получает список предложений текущего пользователя
func (hnd *BidHandler) GetUserBids(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodGet {
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

	var limit int32 = 0
	limitStr := rqt.URL.Query().Get("limit")
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
				return
			}
		}

		limit = int32(limitInt)
		if limit < 0 || limit > 50 {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
				return
			}
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
				return
			}
		}

		offset = int32(offsetInt)
		if offset < 0 {
			errSend := handlers.SendBadReq(wrt)
			if errSend != nil {
				log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
				return
			}
		}
	}
	endIndex := offset + limit

	bids, code, err := hnd.BidRepo.GetUserBids(username, offset, endIndex)
	if err != nil {
		log.Println(err)
		return
	}

	if code == 401 {
		err := "Пользователь не существует или некорректен"
		errResp := handlers.RespondWithError(wrt, err, http.StatusUnauthorized)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
			return
		}
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(bids)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
