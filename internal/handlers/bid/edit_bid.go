package bid

import (
	"encoding/json"
	"log"
	"net/http"
	"tendermanagement/internal/bid"
	"tendermanagement/internal/handlers"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// EditBid редактирует параметры предложения
func (hnd *BidHandler) EditBid(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodPatch {
		errSend := handlers.SendEditBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	vars := mux.Vars(rqt)
	bidID := vars["bidId"]
	bidIDLen := utf8.RuneCountInString(bidID)
	if bidIDLen == 0 || bidIDLen > 100 {
		errSend := handlers.SendEditBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	username := rqt.URL.Query().Get("username")
	usernameLen := utf8.RuneCountInString(username)
	if usernameLen == 0 {
		errSend := handlers.SendEditBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	var bdi bid.BidEditionInput
	err := json.NewDecoder(rqt.Body).Decode(&bdi)
	if err != nil {
		errSend := handlers.SendEditBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	bdNameLen := utf8.RuneCountInString(bdi.Name)
	bdDescLen := utf8.RuneCountInString(bdi.Description)

	switch {
	case bdNameLen > 100:
		errSend := handlers.SendEditBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return

	case bdDescLen > 500:
		errSend := handlers.SendEditBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	bid, code, err := hnd.BidRepo.EditBid(bdi, bidID, username)
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
		err := "Предложение не найдено"
		errResp := handlers.RespondWithError(wrt, err, http.StatusNotFound)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
		}
		return
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(bid)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
