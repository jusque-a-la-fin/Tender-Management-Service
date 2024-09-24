package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse - это errorResponse
type ErrorResponse struct {
	Reason string `json:"reason"`
}

func RespondWithError(wrt http.ResponseWriter, err string, statusCode int) error {
	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(statusCode)
	errorResponse := ErrorResponse{Reason: err}
	errJSON := json.NewEncoder(wrt).Encode(errorResponse)
	return errJSON
}

func SendBadReq(wrt http.ResponseWriter) error {
	err := "Неверный формат запроса или его параметры"
	errResp := RespondWithError(wrt, err, http.StatusBadRequest)
	return errResp
}

func SendEditBadReq(wrt http.ResponseWriter) error {
	err := "Данные неправильно сформированы или не соответствуют требованиям"
	errResp := RespondWithError(wrt, err, http.StatusBadRequest)
	return errResp
}

func SendSubmitBadReq(wrt http.ResponseWriter) error {
	err := "Решение не может быть отправлено"
	errResp := RespondWithError(wrt, err, http.StatusBadRequest)
	return errResp
}

func SendFeedbackBadReq(wrt http.ResponseWriter) error {
	err := "Отзыв не может быть отправлен"
	errResp := RespondWithError(wrt, err, http.StatusBadRequest)
	return errResp
}

func CheckCode(code int) bool {
	switch code {
	case 401, 403, 404:
		return true
	}
	return false
}
