package tender_test

import (
	"log"
	"net/http"
	"tendermanagement/internal/datastore"
	thd "tendermanagement/internal/handlers/tender"
	"tendermanagement/internal/tender"
	"tendermanagement/test"
	"testing"
)

// TestRollbackTenderDoesntExist тестирует случай, когда пользователь не существует или некорректен
func TestRollbackTenderDoesntExist(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	url := "/tenders/ac8656ae-da41-4bfa-a3b5-c4cad399fcdf/rollback/1?username=user31"
	path := "/tenders/{tenderId}/rollback/{version}"

	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodPut, "RollbackTender")
	code := rr.Code
	if code != http.StatusUnauthorized {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusUnauthorized, code)
	}

	expected := "Пользователь не существует или некорректен"
	test.HandleError(t, rr, expected)
}

// TestRollbackTenderForbidden тестирует случай, когда у пользователя недостаточно прав
func TestRollbackTenderForbidden(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	url := "/tenders/ac8656ae-da41-4bfa-a3b5-c4cad399fcdf/rollback/1?username=user15"
	path := "/tenders/{tenderId}/rollback/{version}"

	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodPut, "RollbackTender")
	code := rr.Code
	if code != http.StatusForbidden {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusForbidden, code)
	}

	expected := "Недостаточно прав для выполнения действия"
	test.HandleError(t, rr, expected)
}

// TestRollbackTenderNotFound тестирует случай, когда тендер не найден
func TestRollbackTenderNotFound(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	url := "/tenders/0a78ce19-543a-438b-89f4-20b9f32a610d/rollback/1?username=user11"
	path := "/tenders/{tenderId}/rollback/{version}"

	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodPut, "RollbackTender")
	code := rr.Code
	if code != http.StatusNotFound {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusNotFound, code)
	}

	expected := "Тендер или версия не найдены"
	test.HandleError(t, rr, expected)
}

var testRollbackIncorrectUrls = []string{
	"/tenders//status?username=user11",
	"/tenders/550e8400-e29b-41d4-a716-446655440023/status?username=",
}

// TestRollbackTenderBadRequest тестирует некорректный запрос
func TestRollbackTenderBadRequest(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tdr := tender.NewDBRepo(dtb)
	tenderHandler := &thd.TenderHandler{
		TenderRepo: tdr,
	}

	handler := http.HandlerFunc(tenderHandler.RollbackTender)
	expected := "Неверный формат запроса или его параметры"

	// некорректные параметры url
	for _, testUrl := range testRollbackIncorrectUrls {
		req, err := http.NewRequest(http.MethodPut, testUrl, nil)
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		rr := test.ServeRequest(handler, req)
		test.HandleBadReq(t, rr, expected)
	}

	url := "/tenders/ac8656ae-da41-4bfa-a3b5-c4cad399fcdf/rollback/1?username=user11"

	// некорректный метод запроса
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := test.ServeRequest(handler, req)
	test.HandleBadReq(t, rr, expected)
}

var testRollbackUrls = []string{
	"/tenders/ac8656ae-da41-4bfa-a3b5-c4cad399fcdf/rollback/1?username=user11",
	"/tenders/ac8656ae-da41-4bfa-a3b5-c4cad399fcdf/rollback/2?username=user11",
}

// TestRollbackTender тестирует успешный откат параметров тендера к указанной версии
func TestRollbackTenderOK(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	path := "/tenders/{tenderId}/rollback/{version}"
	for _, testUrl := range testRollbackUrls {
		rr := test.ProcessReq(t, dtb, nil, testUrl, path, http.MethodPut, "RollbackTender")
		if rr.Code != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, rr.Code)
		}
	}
}
