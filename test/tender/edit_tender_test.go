package tender_test

import (
	"bytes"
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

// TestCreateTenderDoesntExist тестирует случай, когда пользователь не существует или некорректен
func TestEditTenderDoesntExist(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		Name:        "Доставка товары Москва - Санкт-Петербург",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType: "Delivery",
	}

	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tenderID, _ := GetParams(dtb)
	rr := receiveCreationResponseRecorder(t, dtb, body, tenderID, "user31")
	code := rr.Code
	if code != http.StatusUnauthorized {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusUnauthorized, code)
	}

	expected := "Пользователь не существует или некорректен"
	test.HandleError(t, rr, expected)
}

// TestEditTenderForbidden тестирует случай, когда у пользователя недостаточно прав
func TestEditTenderForbidden(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		Name:        "Доставка товары Москва - Санкт-Петербург",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType: "Delivery",
	}

	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tenderID, _ := GetParams(dtb)
	rr := receiveCreationResponseRecorder(t, dtb, body, tenderID, "user15")
	code := rr.Code
	if code != http.StatusForbidden {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusForbidden, code)
	}

	expected := "Недостаточно прав для выполнения действия"
	test.HandleError(t, rr, expected)
}

// TestEditTenderForbidden тестирует случай, когда тендер не найден
func TestEditTenderNotFound(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		Name:        "Доставка товары Москва - Санкт-Петербург",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType: "Delivery",
	}

	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	_, username := GetParams(dtb)
	tenderID := "0ef0296a-e6cb-4167-b524-fa13d989d99c"
	rr := receiveCreationResponseRecorder(t, dtb, body, tenderID, username)
	code := rr.Code
	if code != http.StatusNotFound {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusNotFound, code)
	}

	expected := "Тендер не найден"
	test.HandleError(t, rr, expected)
}

var testsEdition = []tender.TenderEditionInput{
	{
		Name:        "ДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставка",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType: "Delivery",
	},
	{
		Name: "Доставка товары Казань - Москва",
		Description: `Нужно доставить оборудование для олимпиады по робототехнике
                      Нужно доставить оборудование для олимпиады по робототехнике
                      Нужно доставить оборудование для олимпиады по робототехнике
                      Нужно доставить оборудование для олимпиады по робототехнике
                      Нужно доставить оборудование для олимпиады по робототехнике
                      Нужно доставить оборудование для олимпиады по робототехнике
                      Нужно доставить оборудование для олимпиады по робототехнике
                      Нужно доставить оборудование для олимпиады по робототехнике
                      Нужно доставить оборудование для олимпиады по робототехнике`,
		ServiceType: "Delivery",
	},
}

var testEditionUrls = []string{
	"/tenders//edit?username=user11",
	"/tenders/550e8400-e29b-41d4-a716-446655440023/edit?username=",
}

// TestEditTenderBadRequest тестирует некорректный запрос
func TestEditTenderBadRequest(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tdr := tender.NewDBRepo(dtb)
	tenderHandler := &thd.TenderHandler{
		TenderRepo: tdr,
	}

	handler := http.HandlerFunc(tenderHandler.EditTender)
	expected := "Неверный формат запроса или его параметры"

	// некорректные параметры url
	for _, testUrl := range testEditionUrls {
		req, err := http.NewRequest(http.MethodPatch, testUrl, nil)
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		rr := test.ServeRequest(handler, req)
		test.HandleBadReq(t, rr, expected)
	}

	tenderID, username := GetParams(dtb)
	url := fmt.Sprintf("/tenders/%s/edit?username=%s", tenderID, username)

	// некорректные параметры тела запроса
	for _, testEd := range testsEdition {
		data, err := json.Marshal(testEd)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
		}

		req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		rr := test.ServeRequest(handler, req)
		test.HandleBadReq(t, rr, expected)
	}

	// некорректный метод запроса
	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := test.ServeRequest(handler, req)
	test.HandleBadReq(t, rr, expected)
}

// TestEditTenderOK1 тестирует успешное редактирование полного названия тендера
func TestEditTenderOK1(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		Name: "Доставка товары Казань - Санкт-Петербург",
	}

	expected := tender.Tender{
		Name: body.Name,
	}

	rr := getCreationResponseRecorder(t, body)
	tend := handleTenderResponse(t, rr)
	if expected.Name != tend.Name {
		t.Errorf("Ожидалось значение Name: %s, но получено: %s", expected.Name, tend.Name)
	}
}

func GetParams(dtb *sql.DB) (string, string) {
	var tenderID, userID string
	query := "SELECT id, user_id FROM tender LIMIT 1;"
	_ = dtb.QueryRow(query).Scan(&tenderID, &userID)

	var username string
	query = "SELECT username FROM employee WHERE id = $1;"
	_ = dtb.QueryRow(query, userID).Scan(&username)
	return tenderID, username
}

// TestEditTenderOK2 тестирует успешное редактирование описания тендера
func TestEditTenderOK2(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
	}

	expected := tender.Tender{
		Description: body.Description,
	}

	rr := getCreationResponseRecorder(t, body)
	tend := handleTenderResponse(t, rr)
	if expected.Description != tend.Description {
		t.Errorf("Ожидалось значение Description: %s, но получено: %s", expected.Description, tend.Description)
	}
}

// TestEditTenderOK3 тестирует успешное редактирование вида услуги, к которой относится тендер
func TestEditTenderOK3(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		ServiceType: "Delivery",
	}

	expected := tender.Tender{
		ServiceType: body.ServiceType,
	}

	rr := getCreationResponseRecorder(t, body)

	tend := handleTenderResponse(t, rr)
	if expected.ServiceType != tend.ServiceType {
		t.Errorf("Ожидалось значение ServiceType: %s, но получено: %s", expected.ServiceType, tend.ServiceType)
	}
}

// TestEditTenderOK4 тестирует успешное редактирование всех параметров тендера
func TestEditTenderOK4(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		Name:        "Доставка товары Москва - Санкт-Петербург",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType: "Delivery",
	}

	expected := tender.Tender{
		Name:        body.Name,
		Description: body.Description,
		ServiceType: body.ServiceType,
	}

	rr := getCreationResponseRecorder(t, body)

	tend := handleTenderResponse(t, rr)
	if expected.Name != tend.Name {
		t.Errorf("Ожидалось значение Name: %s, но получено: %s", expected.Name, tend.Name)
	}

	if expected.Description != tend.Description {
		t.Errorf("Ожидалось значение Description: %s, но получено: %s", expected.Description, tend.Description)
	}

	if expected.ServiceType != tend.ServiceType {
		t.Errorf("Ожидалось значение ServiceType: %s, но получено: %s", expected.ServiceType, tend.ServiceType)
	}
}

// TestEditTenderOK5 тестирует успешное редактирование полного названия и описания тендера
func TestEditTenderOK5(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		Name:        "Доставка товары Казань - Санкт-Петербург",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
	}

	expected := tender.Tender{
		Name:        body.Name,
		Description: body.Description,
	}

	rr := getCreationResponseRecorder(t, body)

	tend := handleTenderResponse(t, rr)
	if expected.Name != tend.Name {
		t.Errorf("Ожидалось значение Name: %s, но получено: %s", expected.Name, tend.Name)
	}

	if expected.Description != tend.Description {
		t.Errorf("Ожидалось значение Description: %s, но получено: %s", expected.Description, tend.Description)
	}
}

// TestEditTenderOK6 тестирует успешное редактирование полного описания тендера и вида услуги, к которой относится тендер
func TestEditTenderOK6(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		Name:        "Доставка товары Казань - Санкт-Петербург",
		ServiceType: "Delivery",
	}

	expected := tender.Tender{
		Name:        body.Name,
		ServiceType: body.ServiceType,
	}

	rr := getCreationResponseRecorder(t, body)

	tend := handleTenderResponse(t, rr)
	if expected.Name != tend.Name {
		t.Errorf("Ожидалось значение Name: %s, но получено: %s", expected.Name, tend.Name)
	}

	if expected.ServiceType != tend.ServiceType {
		t.Errorf("Ожидалось значение ServiceType: %s, но получено: %s", expected.ServiceType, tend.ServiceType)
	}
}

// TestEditTenderOK7 тестирует успешное редактирование описания тендера и вида услуги, к которой относится тендер
func TestEditTenderOK7(t *testing.T) {
	test.SetVars()
	body := tender.TenderEditionInput{
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType: "Delivery",
	}

	expected := tender.Tender{
		Description: body.Description,
		ServiceType: body.ServiceType,
	}

	rr := getCreationResponseRecorder(t, body)

	tend := handleTenderResponse(t, rr)
	if expected.Description != tend.Description {
		t.Errorf("Ожидалось значение Description: %s, но получено: %s", expected.Description, tend.Description)
	}

	if expected.ServiceType != tend.ServiceType {
		t.Errorf("Ожидалось значение ServiceType: %s, но получено: %s", expected.ServiceType, tend.ServiceType)
	}
}

// TestEditTenderOK8 тестирует редактирование тендера без передачи значений, которые нужно обновить
func TestEditTenderOK8(t *testing.T) {
	test.SetVars()

	rr := getCreationResponseRecorder(t, nil)
	handleTenderResponse(t, rr)
}

func getCreationResponseRecorder(t *testing.T, body any) *httptest.ResponseRecorder {
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tenderID, username := GetParams(dtb)
	return receiveCreationResponseRecorder(t, dtb, body, tenderID, username)
}

func receiveCreationResponseRecorder(t *testing.T, dtb *sql.DB, body any, tenderID, username string) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/tenders/%s/edit?username=%s", tenderID, username)
	path := "/tenders/{tenderId}/edit"
	return test.ProcessReq(t, dtb, body, url, path, http.MethodPatch, "EditTender")
}
