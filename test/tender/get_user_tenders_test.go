package tender_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"tendermanagement/internal/datastore"
	thd "tendermanagement/internal/handlers/tender"
	"tendermanagement/internal/tender"
	"tendermanagement/test"
	"testing"
)

// TestGetUserTendersDoesntExist тестирует случай, когда пользователь не существует или некорректен
func TestGetUserTendersDoesntExist(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	url := "/tenders/my?limit=4&offset=2&username=user31"
	path := "/tenders/my"

	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodGet, "GetUserTenders")
	code := rr.Code
	if code != http.StatusUnauthorized {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusUnauthorized, code)
	}

	expected := "Пользователь не существует или некорректен"
	test.HandleError(t, rr, expected)
}

var testGetUserTendersIncorrectUrls = []string{
	"/tenders/my?limit=t&offset=2&username=user11",
	"/tenders/my?limit=-1&offset=2&username=user11",
	"/tenders/my?limit=51&offset=2&username=user11",
	"/tenders/my?limit=3&offset=t&username=user11",
	"/tenders/my?limit=3&offset=-2&username=user11",
	"/tenders/my?limit=3&offset=2&username=",
}

// TestGetUserTendersBadRequest тестирует некорректный запрос
func TestGetUserTendersBadRequest(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tdr := tender.NewDBRepo(dtb)
	tenderHandler := &thd.TenderHandler{
		TenderRepo: tdr,
	}

	handler := http.HandlerFunc(tenderHandler.GetUserTenders)
	expected := "Неверный формат запроса или его параметры"

	// некорректные параметры url
	for _, testUrl := range testGetUserTendersIncorrectUrls {
		req, err := http.NewRequest(http.MethodGet, testUrl, nil)
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		rr := test.ServeRequest(handler, req)
		test.HandleBadReq(t, rr, expected)
	}

	url := "/tenders/my?limit=4&offset=2&username=user11"
	// некорректный метод запроса
	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := test.ServeRequest(handler, req)
	test.HandleBadReq(t, rr, expected)
}

var testGetUserTendersUrls = []string{
	"/tenders/my?limit=3&offset=2&username=user11",
	"/tenders/my?limit=3&username=user11",
	"/tenders/my?offset=2&username=user11",
	"/tenders/my?username=user11",
}

// TestGetUserTenders тестирует успешное получение списка тендеров пользователя
func TestGetUserTendersOK(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	for _, testUrl := range testGetUserTendersUrls {
		req, err := http.NewRequest(http.MethodGet, testUrl, nil)
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		tdr := tender.NewDBRepo(dtb)
		tenderHandler := &thd.TenderHandler{
			TenderRepo: tdr,
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tenderHandler.GetUserTenders)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, rr.Code)
		}
	}
}
