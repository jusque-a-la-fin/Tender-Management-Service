package bid_test

import (
	"net/http"
	"tendermanagement/test"
	"testing"
)

var getUserBidsPath = "/bids/my"
var getUserBidsHttpMethod = http.MethodGet
var getUserBidsMethod = "GetUserBids"

// TestGetUserBidsUnauthorized тестирует случай, когда пользователь не существует или некорректен
func TestGetUserBidsUnauthorized(t *testing.T) {
	url := "/bids/my?username=user34"
	handleGetUserBidsIncorrectUrl(t, url, http.StatusUnauthorized)
}

func handleGetUserBidsIncorrectUrl(t *testing.T, url string, expectedCode int) {
	test.ProcessIncorrectUrl(t, getUserBidsPath, getUserBidsHttpMethod, getUserBidsMethod, url, expectedCode)
}

var getUserBidsBadRequestUrls = []string{
	"/bids/my?username=",
	"/bids/my?limit=t&offset=2",
	"/bids/my?limit=-1&offset=2",
	"/bids/my?limit=51&offset=2",
	"/bids/my?limit=3&offset=t",
	"/bids/my?limit=3&offset=-2",
}

// TestGetUserBidsBadRequest тестирует некорректный запрос
func TestGetUserBidsBadRequest(t *testing.T) {
	handleGetUserBidsIncorrectUrls(t, getUserBidsBadRequestUrls, http.StatusBadRequest)
}

var getUserBidsOKTests = []test.OKTest{
	{Url: "/bids/my?username=user4", ExpectedQuantity: 3},
	{Url: "/bids/my?username=user4&limit=2", ExpectedQuantity: 2},
	{Url: "/bids/my?username=user4&limit=1&offset=1", ExpectedQuantity: 1},
	{Url: "/bids/my?username=user4&offset=2", ExpectedQuantity: 1},
}

// TestGetUserBidsOK тестирует успешное получение списка предложений текущего пользователя
func TestGetUserBidsOK(t *testing.T) {
	test.ProcessGetBidsCorrectUrls(t, getUserBidsPath, getUserBidsHttpMethod, getUserBidsMethod, getUserBidsOKTests)
}

func handleGetUserBidsIncorrectUrls(t *testing.T, urls []string, expectedCode int) {
	test.ProcessIncorrectUrls(t, getUserBidsPath, getUserBidsHttpMethod, getUserBidsMethod, urls, expectedCode)
}
