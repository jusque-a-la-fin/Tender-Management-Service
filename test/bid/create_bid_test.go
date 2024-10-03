package bid_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"tendermanagement/internal/datastore"
	bhd "tendermanagement/internal/handlers/bid"
	"tendermanagement/test"
	"testing"
)

var createBidUrl = "/bids/new"

var testsExistence = []bhd.BidCreationReq{
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "ac8656ae-da41-4bfa-a3b5-c4cad399fcdf",
		AuthorType:  "User",
		AuthorId:    "550e8400-t29b-41d4-a716-446655440002",
	},
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "ac8656ae-da41-4bfa-a3b5-c4cad399fcdf",
		AuthorType:  "User",
		AuthorId:    "551e8400",
	},
}

// TestCreateBidDoesntExist тестирует случай, когда пользователь не существует или некорректен
func TestCreateBidDoesntExist(t *testing.T) {
	test.SetVars()
	for _, testBid := range testsExistence {
		rr := test.ProcessReq(t, nil, testBid, createBidUrl, createBidUrl, http.MethodPost, "CreateBid")
		code := rr.Code
		if code != http.StatusUnauthorized {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusUnauthorized, code)
		}

		expected := "Пользователь не существует или некорректен"
		test.HandleError(t, rr, expected)
	}
}

// TestCreateBidForbidden тестирует случай, когда у пользователя недостаточно прав
func TestCreateBidForbidden(t *testing.T) {
	test.SetVars()
	body := bhd.BidCreationReq{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "1ffac8e1-42d3-4351-a778-751b9dacf4b0",
		AuthorType:  "User",
		AuthorId:    "550e8400-e29b-41d4-a716-44665544001b",
	}

	rr := test.ProcessReq(t, nil, body, createBidUrl, createBidUrl, http.MethodPost, "CreateBid")
	code := rr.Code
	if code != http.StatusForbidden {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusForbidden, code)
	}

	expected := "Недостаточно прав для выполнения действия"
	test.HandleError(t, rr, expected)
}

// TestCreateBidNotFound тестирует случай, когда тендер не найден
func TestCreateBidNotFound(t *testing.T) {
	test.SetVars()
	body := bhd.BidCreationReq{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "00000000-0000-0000-0000-000000000000",
		AuthorType:  "User",
		AuthorId:    "550e8400-e29b-41d4-a716-446655440002",
	}

	rr := test.ProcessReq(t, nil, body, createBidUrl, createBidUrl, http.MethodPost, "CreateBid")
	code := rr.Code
	if code != http.StatusNotFound {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusNotFound, code)
	}

	expected := "Тендер не найден"
	test.HandleError(t, rr, expected)
}

var testsBadRequest = []bhd.BidCreationReq{
	{
		Name:        "",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "6e29150e-ed3d-4c49-8699-362c8fff1f76",
		AuthorType:  "User",
		AuthorId:    "551e8400-e29b-41d4-a716-44665544001t",
	},
	{
		Name:        "ДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставкаДоставка",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "6e29150e-ed3d-4c49-8699-362c8fff1f76",
		AuthorType:  "User",
		AuthorId:    "551e8400-e29b-41d4-a716-44665544001t",
	},
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "",
		TenderID:    "6e29150e-ed3d-4c49-8699-362c8fff1f76",
		AuthorType:  "User",
		AuthorId:    "551e8400-e29b-41d4-a716-44665544001t",
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
		TenderID:   "6e29150e-ed3d-4c49-8699-362c8fff1f76",
		AuthorType: "User",
		AuthorId:   "551e8400-e29b-41d4-a716-44665544001t",
	},
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "6e29150e-ed3d-4c49-8699-362c8fff1f766e29150e-ed3d-4c49-8699-362c8fff1f766e29150e-ed3d-4c49-8699-362c8fff1f76",
		AuthorType:  "User",
		AuthorId:    "551e8400-e29b-41d4-a716-44665544001t",
	},
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "",
		AuthorType:  "User",
		AuthorId:    "551e8400-e29b-41d4-a716-44665544001t",
	},
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "6e29150e-ed3d-4c49-8699-362c8fff1f76",
		AuthorType:  "Usr",
		AuthorId:    "551e8400-e29b-41d4-a716-44665544001t",
	},
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "6e29150e-ed3d-4c49-8699-362c8fff1f76",
		AuthorType:  "User",
		AuthorId:    "",
	},
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "6e29150e-ed3d-4c49-8699-362c8fff1f76",
		AuthorType:  "User",
		AuthorId:    "551e8400-e29b-41d4-a716-44665544001t551e8400-e29b-41d4-a716-44665544001t551e8400-e29b-41d4-a716-44665544001t",
	},
}

// TestCreateBidBadRequest тестирует некорректный запрос
func TestCreateBidBadRequest(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	bidHandler := test.GetBidHandler(dtb)
	handler := http.HandlerFunc(bidHandler.CreateBid)
	expected := "Неверный формат запроса или его параметры"

	// некорректные параметры тела запроса
	for _, testBR := range testsBadRequest {
		data, err := json.Marshal(testBR)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, createBidUrl, bytes.NewBuffer(data))
		if err != nil {
			t.Fatal("Ошибка создания объекта *http.Request:", err)
		}

		rr := test.ServeRequest(handler, req)
		test.HandleBadReq(t, rr, expected)
	}

	// некорректный метод запроса
	req, err := http.NewRequest(http.MethodGet, createBidUrl, nil)
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	rr := test.ServeRequest(handler, req)
	test.HandleBadReq(t, rr, expected)
}

var testsCreation = []bhd.BidCreationReq{
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "8a05f051-88f2-4113-a0ef-f72424f07919",
		AuthorType:  "User",
		AuthorId:    "550e8400-e29b-41d4-a716-446655440002",
	},
	{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
		TenderID:    "8a05f051-88f2-4113-a0ef-f72424f07919",
		AuthorType:  "Organization",
		AuthorId:    "550e8400-e29b-41d4-a716-446655440022",
	},
}

// TestCreateBidOK тестирует успешный ответ
func TestCreateBidOK(t *testing.T) {
	test.SetVars()

	for _, testCreation := range testsCreation {
		rr := test.ProcessReq(t, nil, testCreation, createBidUrl, createBidUrl, http.MethodPost, "CreateBid")
		bid := test.HandleBidResponse(t, rr)

		name := bid.Name
		bName := testCreation.Name
		if name != bName {
			t.Errorf("Ожидалось значение Name: %s, но получено: %s", bName, name)
		}

		description := bid.Description
		bDescription := testCreation.Description
		if description != bDescription {
			t.Errorf("Ожидалось значение Description: %s, но получено: %s", bDescription, description)
		}

		status := bid.Status
		expectedStatus := "Created"
		if string(status) != expectedStatus {
			t.Errorf("Ожидалось значение Status: %s, но получено: %s", expectedStatus, status)
		}

		tenderID := bid.TenderID
		bTenderID := testCreation.TenderID
		if tenderID != bTenderID {
			t.Errorf("Ожидалось значение Description: %s, но получено: %s", bTenderID, tenderID)
		}

		authorType := bid.AuthorType
		bAuthorType := testCreation.AuthorType
		if string(authorType) != bAuthorType {
			t.Errorf("Ожидалось значение Description: %s, но получено: %s", bAuthorType, authorType)
		}

		authorID := bid.AuthorID
		bAuthorID := testCreation.AuthorId
		if authorID != bAuthorID {
			t.Errorf("Ожидалось значение Description: %s, но получено: %s", bAuthorID, authorID)
		}

		version := bid.Version
		if version != 1 {
			t.Errorf("Ожидалось значение Description: %d, но получено: %d", 1, version)
		}
	}
}
