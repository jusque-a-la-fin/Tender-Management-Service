package tender_test

import (
	"fmt"
	"log"
	"net/http"
	"tendermanagement/internal/datastore"
	thd "tendermanagement/internal/handlers/tender"
	"tendermanagement/internal/tender"
	"tendermanagement/test"
	"testing"
)

// TestUpdateTenderStatusDoesntExist тестирует случай, когда пользователь не существует или некорректен
func TestUpdateTenderStatusDoesntExist(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tenderID, _ := GetParams(dtb)
	url := fmt.Sprintf("/tenders/%s/status?status=%s&username=%s", tenderID, "Published", "user31")
	path := "/tenders/{tenderId}/status"

	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodPut, "UpdateTenderStatus")
	code := rr.Code
	if code != http.StatusUnauthorized {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusUnauthorized, code)
	}

	expected := "Пользователь не существует или некорректен"
	test.HandleError(t, rr, expected)
}

// TestUpdateTenderStatusForbidden тестирует случай, когда у пользователя недостаточно прав
func TestUpdateTenderStatusForbidden(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tenderID, _ := GetParams(dtb)
	url := fmt.Sprintf("/tenders/%s/status?status=%s&username=%s", tenderID, "Published", "user10")
	path := "/tenders/{tenderId}/status"

	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodPut, "UpdateTenderStatus")
	code := rr.Code
	if code != http.StatusForbidden {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusForbidden, code)
	}

	expected := "Недостаточно прав для выполнения действия"
	test.HandleError(t, rr, expected)
}

// TestUpdateTenderStatusNotFound тестирует случай, когда тендер не найден
func TestUpdateTenderStatusNotFound(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tenderID := "00000000-0000-0000-0000-000000000000"
	url := fmt.Sprintf("/tenders/%s/status?status=%s&username=%s", tenderID, "Published", "user10")
	path := "/tenders/{tenderId}/status"

	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodPut, "UpdateTenderStatus")
	code := rr.Code
	if code != http.StatusNotFound {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusNotFound, code)
	}

	expected := "Тендер не найден"
	test.HandleError(t, rr, expected)
}

var testUpdateStatusIncorrectUrls = []string{
	"/tenders//status?status=Published&username=user11",
	"/tenders/00000000-0000-0000-0000-000000000000/status?status=&username=user11",
	"/tenders/00000000-0000-0000-0000-000000000000/status?status=Published&username=",

	"/tenders/something/status?status=Published&username=user11",
	"/tenders/00000000-0000-0000-0000-000000000000/status?status=Default&username=user11",
	"/tenders/00000000-0000-0000-0000-000000000000/status?status=Published&username=33",
}

// TestUpdateTenderStatusBadRequest тестирует некорректный запрос
func TestUpdateTenderStatusBadRequest(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tdr := tender.NewDBRepo(dtb)
	tenderHandler := &thd.TenderHandler{
		TenderRepo: tdr,
	}

	handler := http.HandlerFunc(tenderHandler.UpdateTenderStatus)
	expected := "Неверный формат запроса или его параметры"

	// некорректные параметры url
	for _, testUrl := range testUpdateStatusIncorrectUrls {
		req, err := http.NewRequest(http.MethodPut, testUrl, nil)
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		rr := test.ServeRequest(handler, req)
		test.HandleBadReq(t, rr, expected)
	}

	url := "/tenders/0a78ce19-546a-438b-89f4-20b9f32a610d/status?status=Published&username=user11"

	// некорректный метод запроса
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := test.ServeRequest(handler, req)
	test.HandleBadReq(t, rr, expected)
}

// TestUpdateTenderStatusOK тестирует изменение статуса тендера по его идентификатору
func TestUpdateTenderStatusOK(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	url := "/tenders/ac8656ae-da41-4bfa-a3b5-c4cad399fcdf/status?status=Published&username=user11"
	path := "/tenders/{tenderId}/status"
	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodPut, "UpdateTenderStatus")
	code := rr.Code
	if code != http.StatusOK {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, code)
	}
}
