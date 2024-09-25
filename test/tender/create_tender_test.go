package tender_test

import (
	"bytes"
	"encoding/json"
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
func TestCreateTenderDoesntExist(t *testing.T) {
	test.SetVars()
	body := thd.TenderCreationReq{
		Name:            "Доставка товары Казань - Москва",
		Description:     "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType:     "Delivery",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "user31",
	}

	url := "/tenders/new"
	rr := test.ProcessReq(t, nil, body, url, url, http.MethodPost, "CreateTender")
	code := rr.Code
	if code != http.StatusUnauthorized {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusUnauthorized, code)
	}

	expected := "Пользователь не существует или некорректен"
	test.HandleError(t, rr, expected)
}

// TestCreateTenderForbidden тестирует случай, когда у пользователя недостаточно прав
func TestCreateTenderForbidden(t *testing.T) {
	test.SetVars()
	body := thd.TenderCreationReq{
		Name:            "Доставка товары Казань - Москва",
		Description:     "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType:     "Delivery",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "user15",
	}

	url := "/tenders/new"
	rr := test.ProcessReq(t, nil, body, url, url, http.MethodPost, "CreateTender")
	code := rr.Code
	if code != http.StatusForbidden {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusForbidden, code)
	}

	expected := "Недостаточно прав для выполнения действия"
	test.HandleError(t, rr, expected)
}

var testsCreation = []thd.TenderCreationReq{
	{
		Name:            "",
		Description:     "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType:     "Delivery",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "user11",
	},
	{
		Name:            "ДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставка",
		Description:     "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType:     "Delivery",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "user11",
	},
	{
		Name:            "Доставка товары Казань - Москва",
		Description:     "",
		ServiceType:     "Delivery",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "user11",
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
		ServiceType:     "Delivery",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "user11",
	},
	{
		Name:            "Доставка товары Казань - Москва",
		Description:     "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType:     "Delivery",
		OrganizationID:  "",
		CreatorUsername: "user11",
	},
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType: "Delivery",
		OrganizationID: `550e8400-e29b-41d4-a716-446655440023550e8400-e29b-41d4-a716-446655440023
		                 550e8400-e29b-41d4-a716-446655440023550e8400-e29b-41d4-a716-446655440023`,
		CreatorUsername: "user11",
	},
	{
		Name:            "Доставка товары Казань - Москва",
		Description:     "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType:     "Delivery",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "",
	},
	{
		Name:            "Доставка товары Казань - Москва",
		Description:     "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType:     "",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "user11",
	},
	{
		Name:            "Доставка товары Казань - Москва",
		Description:     "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType:     "Deliver",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "user11",
	},
}

// TestCreateTenderBadRequest тестирует некорректный запрос
func TestCreateTenderBadRequest(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tdr := tender.NewDBRepo(dtb)
	tenderHandler := &thd.TenderHandler{
		TenderRepo: tdr,
	}

	handler := http.HandlerFunc(tenderHandler.CreateTender)
	expected := "Неверный формат запроса или его параметры"
	url := "/tenders/new"
	// некорректные параметры тела запроса
	for _, testCr := range testsCreation {
		data, err := json.Marshal(testCr)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		rr := test.ServeRequest(handler, req)
		test.HandleBadReq(t, rr, expected)
	}

	// некорректный метод запроса
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := test.ServeRequest(handler, req)
	test.HandleBadReq(t, rr, expected)
}

// TestCreateTenderOK тестирует успешный ответ
func TestCreateTenderOK(t *testing.T) {
	test.SetVars()
	body := thd.TenderCreationReq{
		Name:            "Доставка товары Казань - Москва",
		Description:     "Нужно доставить оборудование для олимпиады по робототехнике",
		ServiceType:     "Delivery",
		OrganizationID:  "550e8400-e29b-41d4-a716-446655440023",
		CreatorUsername: "user11",
	}

	url := "/tenders/new"
	rr := test.ProcessReq(t, nil, body, url, url, http.MethodPost, "CreateTender")
	tend := handleTenderResponse(t, rr)

	name := tend.Name
	bName := body.Name
	if name != bName {
		t.Errorf("Ожидалось значение Name: %s, но получено: %s", name, bName)
	}

	description := tend.Description
	bDescription := body.Description
	if description != bDescription {
		t.Errorf("Ожидалось значение Description: %s, но получено: %s", description, bDescription)
	}

	serviceType := tend.ServiceType
	bServiceType := body.ServiceType
	if string(serviceType) != bServiceType {
		t.Errorf("Ожидалось значение ServiceType: %s, но получено: %s", serviceType, bServiceType)
	}

	organizationID := tend.OrganizationID
	bOrganizationID := body.OrganizationID
	if organizationID != bOrganizationID {
		t.Errorf("Ожидалось значение OrganizationID: %s, но получено: %s", organizationID, bOrganizationID)
	}

	status := tend.Status
	expectedStatus := "Created"
	if string(status) != expectedStatus {
		t.Errorf("Ожидалось значение Status: %s, но получено: %s", status, expectedStatus)
	}
}

func handleTenderResponse(t *testing.T, rr *httptest.ResponseRecorder) *tender.Tender {
	test.CheckCodeAndMime(t, rr)

	var tend tender.Tender
	err := json.Unmarshal(rr.Body.Bytes(), &tend)
	if err != nil {
		t.Fatalf("Ошибка десериализации тела ответа сервера: %v", err)
	}

	return &tend
}
