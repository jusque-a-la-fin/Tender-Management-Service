package bid

import (
	"encoding/json"
	"log"
	"net/http"
	"tendermanagement/internal/bid"
	"tendermanagement/internal/handlers"
	"unicode/utf8"
)

// CreateBid cоздает предложение для существующего тендера
func (hnd *BidHandler) CreateBid(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodPost {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	var brq BidCreationReq
	err := json.NewDecoder(rqt.Body).Decode(&brq)
	if err != nil {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
		}
		return
	}

	bdName := utf8.RuneCountInString(brq.Name)
	bdDesc := utf8.RuneCountInString(brq.Description)
	tndID := utf8.RuneCountInString(brq.TenderId)
	atID := utf8.RuneCountInString(brq.AuthorId)

	switch {
	case bdName == 0 || bdName > 100:
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}

	case bdDesc == 0 || bdDesc > 500:
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}

	case tndID == 0 || tndID > 100:
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}

	case atID == 0 || atID > 100:
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	fail := bid.CheckAuthorType(brq.AuthorType)
	if fail {
		errSend := handlers.SendBadReq(wrt)
		if errSend != nil {
			log.Printf("ошибка отправки сообщения о bad request: %v\n", errSend)
			return
		}
	}

	bdi := bid.BidCreationInput{
		Name:        brq.Name,
		Description: brq.Description,
		TenderId:    brq.TenderId,
		AuthorType:  bid.AuthorTypeEnum(brq.AuthorType),
		AuthorId:    brq.AuthorId,
	}

	nbd, code, err := hnd.BidRepo.CreateBid(bdi)
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
		err := "Тендер не найден"
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
