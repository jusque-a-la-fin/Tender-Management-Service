package tender_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"tendermanagement/internal/datastore"
	thd "tendermanagement/internal/handlers/tender"
	"tendermanagement/internal/tender"
	"tendermanagement/test"
	"testing"
)

// TestGetTenderStatusDoesntExist тестирует случай, когда пользователь не существует или некорректен
func TestGetTenderStatusDoesntExist(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tenderID, _ := GetParams(dtb)
	rr := receiveStatusResponseRecorder(t, dtb, nil, tenderID, "user31")
	code := rr.Code
	if code != http.StatusUnauthorized {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusUnauthorized, code)
	}

	expected := "Пользователь не существует или некорректен"
	test.HandleError(t, rr, expected)
}

// TestGetTenderStatusForbidden тестирует случай, когда у пользователя недостаточно прав
func TestGetTenderStatusForbidden(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tenderID, _ := GetParams(dtb)
	rr := receiveStatusResponseRecorder(t, dtb, nil, tenderID, "user15")
	code := rr.Code
	if code != http.StatusForbidden {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusForbidden, code)
	}

	expected := "Недостаточно прав для выполнения действия"
	test.HandleError(t, rr, expected)
}

// TestGetTenderStatusNotFound тестирует случай, когда тендер не найден
func TestGetTenderStatusNotFound(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	_, username := GetParams(dtb)
	tenderID := "0ef0296a-e6cb-4167-b524-fa13d989d99c"
	rr := receiveStatusResponseRecorder(t, dtb, nil, tenderID, username)
	code := rr.Code
	if code != http.StatusNotFound {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusNotFound, code)
	}

	expected := "Тендер не найден"
	test.HandleError(t, rr, expected)
}

var testGetStatusUrls = []string{
	"/tenders//status?username=user11",
	"/tenders/550e8400-e29b-41d4-a716-446655440023/status?username=",
}

// TestGetTenderStatusBadRequest тестирует некорректный запрос
func TestGetTenderStatusBadRequest(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tdr := tender.NewDBRepo(dtb)
	tenderHandler := &thd.TenderHandler{
		TenderRepo: tdr,
	}

	handler := http.HandlerFunc(tenderHandler.GetTenderStatus)
	expected := "Неверный формат запроса или его параметры"

	// некорректные параметры url
	for _, testUrl := range testGetStatusUrls {
		req, err := http.NewRequest(http.MethodPatch, testUrl, nil)
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		rr := test.ServeRequest(handler, req)
		test.HandleBadReq(t, rr, expected)
	}

	tenderID, username := GetParams(dtb)
	url := fmt.Sprintf("/tenders/%s/status?username=%s", tenderID, username)
	// некорректный метод запроса
	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := test.ServeRequest(handler, req)
	test.HandleBadReq(t, rr, expected)
}

// TestGetTenderStatusOK проверяет успешное получение статуса тендера
func TestGetTenderStatusOK(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tenderID, username := GetParams(dtb)
	rr := receiveStatusResponseRecorder(t, dtb, nil, tenderID, username)
	test.CheckCodeAndMime(t, rr)

	type tenderStatus struct {
		TenderStatus string `json:"tenderStatus"`
	}

	var tst tenderStatus
	err = json.Unmarshal(rr.Body.Bytes(), &tst)
	if err != nil {
		t.Fatalf("Ошибка десериализации тела ответа сервера: %v", err)
	}

	fail := tender.CheckStatus(string(tst.TenderStatus))
	if fail {
		t.Errorf("Ожидалось, что статус тендера будет иметь значение из списка (Created, Published, Closed), но получено: %s", tst.TenderStatus)
	}
}

func receiveStatusResponseRecorder(t *testing.T, dtb *sql.DB, body any, tenderID, username string) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/tenders/%s/status?username=%s", tenderID, username)
	path := "/tenders/{tenderId}/status"
	return test.ProcessReq(t, dtb, body, url, path, http.MethodGet, "GetTenderStatus")
}
