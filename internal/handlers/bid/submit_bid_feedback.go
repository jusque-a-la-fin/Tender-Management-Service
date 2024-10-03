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

// SubmitBidFeedback отправляет отзыв по предложению
func (hnd *BidHandler) SubmitBidFeedback(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodPut {
		errSend := handlers.SendFeedbackBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	vars := mux.Vars(rqt)
	bidID := vars["bidId"]
	bidIDLen := utf8.RuneCountInString(bidID)
	if bidIDLen == 0 || bidIDLen > 100 {
		errSend := handlers.SendFeedbackBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	bidFeedback := rqt.URL.Query().Get("bidFeedback")
	bidFeedbackLen := utf8.RuneCountInString(bidFeedback)
	if bidFeedbackLen == 0 || bidFeedbackLen > 1000 {
		errSend := handlers.SendFeedbackBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	username := rqt.URL.Query().Get("username")
	usernameLen := utf8.RuneCountInString(username)
	if usernameLen == 0 {
		errSend := handlers.SendFeedbackBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	bfi := bid.BidFeedbackInput{
		BidID:       bidID,
		BidFeedback: bidFeedback,
		Username:    username,
	}

	nbd, code, err := hnd.BidRepo.SubmitBidFeedback(bfi)
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

	case 404:
		err := "Предложение не найдено"
		errResp := handlers.RespondWithError(wrt, err, http.StatusForbidden)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
		}
		return
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(nbd)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
