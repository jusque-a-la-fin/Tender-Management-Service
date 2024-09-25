package test

import (
	"net/http"
	"net/http/httptest"
	"tendermanagement/internal/handlers"
	"testing"
)

var url = "/ping"

// TestCheckServerOK тестирует успешный ответ
func TestCheckServerOK(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CheckServer)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, rr.Code)
	}

	if mime := rr.Header().Get("Content-Type"); mime != "text/plain" {
		t.Errorf("Заголовок Content-Type должен иметь MIME-тип text/plain, но имеет %s", mime)
	}

	errDesc := "ok"
	if body := rr.Body.String(); body != "ok" {
		t.Errorf("В теле ответа ожидалось: %s, но получено %s", errDesc, body)
	}
}

// TestCheckServerBadRequest тестирует случай, когда запрос не соответствует требованиям
func TestCheckServerBadRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CheckServer)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusBadRequest, rr.Code)
	}
}

// TestCheckServerServerNotReady тестирует внутреннюю ошибку сервера
func TestCheckServerServerNotReady(t *testing.T) {
	handlers.IsReady = false
	defer func() { handlers.IsReady = true }()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CheckServer)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusInternalServerError, rr.Code)
	}
}
