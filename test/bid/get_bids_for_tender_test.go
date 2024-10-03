package bid_test

import (
	"net/http"
	"tendermanagement/test"
	"testing"
)

var getBidsForTenderPath = "/bids/{tenderId}/list"
var getBidsForTenderHttpMethod = http.MethodGet
var getBidsForTenderMethod = "GetBidsForTender"

// TestGetBidsForTenderUnauthorized тестирует случай, когда пользователь не существует или некорректен
func TestGetBidsForTenderUnauthorized(t *testing.T) {
	url := "/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/list?username=user34"
	handleGetBidsForTenderIncorrectUrl(t, url, http.StatusUnauthorized)
}

// TestGetBidsForTenderForbidden тестирует случай, когда у пользователя недостаточно прав
func TestGetBidsForTenderForbidden(t *testing.T) {
	url := "/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/list?username=user12"
	handleGetBidsForTenderIncorrectUrl(t, url, http.StatusForbidden)
}

var getBidsForTenderNotFoundUrls = []string{
	// тендер не найден
	"/bids/00000000-0000-0000-0000-000000000000/list?username=user12",
	// предложение не найдено
	"/bids/1ffac8e1-42d3-4351-a778-751b9dacf4b0/list?username=user11",
}

// TestGetBidsForTenderNotFound тестирует случай, когда тендер или предложение не найдены
func TestGetBidsForTenderNotFound(t *testing.T) {
	handleGetBidsForTenderIncorrectUrls(t, getBidsForTenderNotFoundUrls, http.StatusNotFound)
}

var getBidsForTenderBadRequestUrls = []string{
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daedca402b26-1bc6-4f97-b41a-0373aec3daedca402b26-1bc6-4f97-b41a-0373aec3daed/list?username=user11",
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/list?username=",
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/list?username=user11&limit=t&offset=2",
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/list?username=user11&limit=-1&offset=2",
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/list?username=user11&limit=51&offset=2",
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/list?username=user11&limit=3&offset=t",
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/list?username=user11&limit=3&offset=-2",
}

// TestGetBidsForTenderBadRequest тестирует некорректный запрос
func TestGetBidsForTenderBadRequest(t *testing.T) {
	handleGetBidsForTenderIncorrectUrls(t, getBidsForTenderBadRequestUrls, http.StatusBadRequest)
}

var getBidsForTenderOKTests = []test.OKTest{
	{Url: "/bids/0ef0296a-e6cb-4167-b524-fa13d989d95c/list?username=user3", ExpectedQuantity: 9},
	{Url: "/bids/0ef0296a-e6cb-4167-b524-fa13d989d95c/list?username=user3&limit=3", ExpectedQuantity: 3},
	{Url: "/bids/0ef0296a-e6cb-4167-b524-fa13d989d95c/list?username=user3&limit=3&offset=2", ExpectedQuantity: 3},
	{Url: "/bids/0ef0296a-e6cb-4167-b524-fa13d989d95c/list?username=user3&offset=3", ExpectedQuantity: 6},
}

// TestGetBidsForTenderOK тестирует успешное получение списка предложений для тендера
func TestGetBidsForTenderOK(t *testing.T) {
	test.ProcessGetBidsCorrectUrls(t, getBidsForTenderPath, getBidsForTenderHttpMethod, getBidsForTenderMethod, getBidsForTenderOKTests)
}

func handleGetBidsForTenderIncorrectUrl(t *testing.T, url string, expectedCode int) {
	test.ProcessIncorrectUrl(t, getBidsForTenderPath, getBidsForTenderHttpMethod, getBidsForTenderMethod, url, expectedCode)
}

func handleGetBidsForTenderIncorrectUrls(t *testing.T, urls []string, expectedCode int) {
	test.ProcessIncorrectUrls(t, getBidsForTenderPath, getBidsForTenderHttpMethod, getBidsForTenderMethod, urls, expectedCode)
}
