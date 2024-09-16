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

// SubmitBidDecision отправляет решение (одобряет или отклоняет) по предложению
func (hnd *BidHandler) SubmitBidDecision(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodPut {
		errSend := handlers.SendSubmitBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	vars := mux.Vars(rqt)
	bidID := vars["bidId"]
	bidIDLen := utf8.RuneCountInString(bidID)
	if bidIDLen == 0 || bidIDLen > 100 {
		errSend := handlers.SendSubmitBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	decision := rqt.URL.Query().Get("decision")
	fail := bid.CheckDecision(decision)
	if fail {
		errSend := handlers.SendSubmitBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	username := rqt.URL.Query().Get("username")
	usernameLen := utf8.RuneCountInString(username)
	if usernameLen == 0 {
		errSend := handlers.SendSubmitBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	bsi := bid.BidSubmissionInput{
		BID:      bidID,
		Decision: bid.DecisionEnum(decision),
		Username: username,
	}

	nbd, code, err := hnd.BidRepo.SubmitBidDecision(bsi)
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
		errResp := handlers.RespondWithError(wrt, err, http.StatusForbidden)
		if errResp != nil {
			log.Printf("ошибка отправки сообщения об ошибке: %d (%s): %v\n", code, err, errResp)
			return
		}
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(nbd)
	if errJSON != nil {
		log.Printf("ошибка отправки тела ответа: %v\n", errJSON)
	}
}
