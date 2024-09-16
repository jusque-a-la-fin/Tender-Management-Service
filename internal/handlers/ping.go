package handlers

import (
	"log"
	"net/http"
)

var IsReady bool = true

// CheckServer проверяет доступность сервера
func CheckServer(wrt http.ResponseWriter, rqt *http.Request) {
	if rqt.Method != http.MethodGet {
		wrt.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsReady {
		wrt.WriteHeader(http.StatusInternalServerError)
		return
	}

	wrt.Header().Set("Content-Type", "text/plain")
	wrt.WriteHeader(http.StatusOK)
	_, err := wrt.Write([]byte("ok"))
	if err != nil {
		log.Printf("ошибка записи в тело ответа %v\n", err)
		return
	}
}
