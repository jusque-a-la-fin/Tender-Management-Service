package bid_test

import (
	"encoding/json"
	"log"
	"net/http"
	"tendermanagement/internal/bid"
	"tendermanagement/internal/datastore"
	"tendermanagement/test"
	"testing"
)

var getBidStatusPath = "/bids/{bidId}/status"
var getBidStatusHttpMethod = http.MethodGet
var getBidStatusMethod = "GetBidStatus"

var getBidStatusUnauthorizedUrls = []string{
	"/bids/4820469e-0a87-43d2-b139-0eb5e253cbfa/status?username=user35",
	"/bids/65d94d94-1319-4175-b2cc-43dae747bffa/status?username=Organization 9",
}

// TestGetBidStatusUnauthorized тестирует случай, когда пользователь не существует или некорректен
func TestGetBidStatusUnauthorized(t *testing.T) {
	handleGetBidStatusOKIncorrectUrls(t, getBidStatusUnauthorizedUrls, http.StatusUnauthorized)
}

var getBidStatusForbiddenUrls = []string{
	"/bids/4820469e-0a87-43d2-b139-0eb5e253cbfa/status?username=user14",
	"/bids/6820469e-0a87-43d2-b139-0eb5e253cbfa/status?username=Organization 2",
}

// TestGetBidStatusForbidden тестирует случай, когда у пользователя недостаточно прав
func TestGetBidStatusForbidden(t *testing.T) {
	handleGetBidStatusOKIncorrectUrls(t, getBidStatusForbiddenUrls, http.StatusForbidden)
}

// TestGetBidStatusNotFound тестирует случай, когда предложение не найдено
func TestGetBidStatusNotFound(t *testing.T) {
	url := "/bids/00000000-0000-0000-0000-000000000000/status?username=user14"
	handleGetBidStatusOKIncorrectUrl(t, url, http.StatusNotFound)
}

var getBidStatusBadRequestUrls = []string{
	"/bids/4820469e-0a87-43d2-b139-0eb5e253cbfa4820469e-0a87-43d2-b139-0eb5e253cbfa4820469e-0a87-43d2-b139-0eb5e253cbfa/status?username=user1",
	"/bids/4820469e-0a87-43d2-b139-0eb5e253cbfa/status?username=",
}

// TestGetBidStatusBadRequest тестирует некорректный запрос
func TestGetBidStatusBadRequest(t *testing.T) {
	handleGetBidStatusOKIncorrectUrls(t, getBidStatusBadRequestUrls, http.StatusBadRequest)
}

type getBidStatusOKTest struct {
	url    string
	status bid.StatusEnum
}

var getBidStatusOKTests = []getBidStatusOKTest{
	{url: "/bids/4820469e-0a87-43d2-b139-0eb5e253cbfa/status?username=user1", status: "Created"},
	{url: "/bids/2820469e-0a87-43d2-b139-0eb5e253cbfa/status?username=Organization 1", status: "Created"},
}

// TestGetBidStatusOK тестирует успешное получение текущего статуса предложения
func TestGetBidStatusOK(t *testing.T) {
	processGetStatusCorrectUrls(t, getBidStatusPath, getBidStatusHttpMethod, getBidStatusMethod)
}

func handleGetBidStatusOKIncorrectUrls(t *testing.T, urls []string, expectedCode int) {
	test.ProcessIncorrectUrls(t, getBidStatusPath, getBidStatusHttpMethod, getBidStatusMethod, urls, expectedCode)
}

func handleGetBidStatusOKIncorrectUrl(t *testing.T, url string, expectedCode int) {
	test.ProcessIncorrectUrl(t, getBidStatusPath, getBidStatusHttpMethod, getBidStatusMethod, url, expectedCode)
}

func processGetStatusCorrectUrls(t *testing.T, path, httpMethod, method string) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	for idx := 0; idx < len(getBidStatusOKTests); idx++ {
		url := getBidStatusOKTests[idx].url
		rr := test.ProcessReq(t, dtb, nil, url, path, httpMethod, method)
		test.CheckCodeAndMime(t, rr)

		var status bid.StatusEnum
		decoder := json.NewDecoder(rr.Body)
		err := decoder.Decode(&status)
		if err != nil {
			t.Fatalf("Ошибка десериализации тела ответа сервера: %v", err)
		}

		expectedStatus := getBidStatusOKTests[idx].status
		if status != expectedStatus {
			t.Fatalf("Ожидался статус: %s, но получен: %s", expectedStatus, status)
		}
	}
}
