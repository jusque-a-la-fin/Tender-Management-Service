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

var testGetTendersIncorrectUrls = []string{
	"/tenders?limit=t&offset=71&service_type=Delivery",
	"/tenders?limit=-1&offset=71&service_type=Delivery",
	"/tenders?limit=51&offset=71&service_type=Delivery",
	"/tenders?limit=4&offset=t&service_type=Delivery",
	"/tenders?limit=4&offset=-5&service_type=Delivery",
	"/tenders?limit=4&offset=71&service_type=Deliver",
}

// TestGetTendersBadRequest тестирует некорректный запрос
func TestGetTendersBadRequest(t *testing.T) {
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
	for _, testUrl := range testGetTendersIncorrectUrls {
		req, err := http.NewRequest(http.MethodGet, testUrl, nil)
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		rr := test.ServeRequest(handler, req)
		test.HandleBadReq(t, rr, expected)
	}

	url := "/tenders?limit=4&offset=71&service_type=Delivery"
	// некорректный метод запроса
	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := test.ServeRequest(handler, req)
	test.HandleBadReq(t, rr, expected)
}

var testGetTendersUrls = []string{
	"/tenders?limit=4&offset=2",
	"/tenders?offset=2",
	"/tenders?limit=4",
	"/tenders",

	"/tenders?limit=4&offset=2&service_type=Delivery",
	"/tenders?offset=2&service_type=Delivery",
	"/tenders?limit=4&service_type=Delivery",
	"/tenders?service_type=Delivery",

	"/tenders?limit=4&offset=2&service_type=Delivery&service_type=Manufacture",
	"/tenders?offset=2&service_type=Delivery&service_type=Manufacture",
	"/tenders?limit=4&service_type=Delivery&service_type=Manufacture",
	"/tenders?service_type=Delivery&service_type=Manufacture",

	"/tenders?limit=4&offset=2&service_type=Delivery&service_type=Construction",
	"/tenders?offset=2&service_type=Delivery&service_type=Construction",
	"/tenders?limit=4&service_type=Delivery&service_type=Construction",
	"/tenders?service_type=Delivery&service_type=Construction",

	"/tenders?limit=4&offset=2&service_type=Manufacture&service_type=Construction",
	"/tenders?offset=2&service_type=Manufacture&service_type=Construction",
	"/tenders?limit=4&service_type=Manufacture&service_type=Construction",
	"/tenders?service_type=Manufacture&service_type=Construction",

	"/tenders?limit=4&offset=2&service_type=Delivery&service_type=Manufacture&service_type=Construction",
	"/tenders?offset=2&service_type=Delivery&service_type=Manufacture&service_type=Construction",
	"/tenders?limit=4&service_type=Delivery&service_type=Manufacture&service_type=Construction",
	"/tenders?service_type=Delivery&service_type=Manufacture&service_type=Construction",
}

// TestGetTendersOK тестирует получение списка тендеров с возможностью фильтрации по типу услуг
func TestGetTendersOK(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	for _, testUrl := range testGetTendersUrls {
		req, err := http.NewRequest(http.MethodGet, testUrl, nil)
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		tdr := tender.NewDBRepo(dtb)
		tenderHandler := &thd.TenderHandler{
			TenderRepo: tdr,
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tenderHandler.GetTenders)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, rr.Code)
		}
	}
}
