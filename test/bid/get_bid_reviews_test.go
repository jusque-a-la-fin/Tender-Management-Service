package bid

import (
	"encoding/json"
	"log"
	"net/http"
	"tendermanagement/internal/bid"
	"tendermanagement/internal/datastore"
	"tendermanagement/test"
	"testing"
)

var getBidReviewsPath = "/bids/{tenderId}/reviews"
var getBidReviewsHttpMethod = http.MethodGet
var getBidReviewsMethod = "GetBidReviews"

var getBidReviewsUnauthorizedUrls = []string{
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=user92",
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user92&requesterUsername=user5",
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/reviews?authorUsername=Organization 7&requesterUsername=user11",
}

// TestGetBidReviewsUnauthorized тестирует случай, когда пользователь не существует или некорректен
func TestGetBidReviewsUnauthorized(t *testing.T) {
	handleGetBidReviewsIncorrectUrls(t, getBidReviewsUnauthorizedUrls, http.StatusUnauthorized)
}

var getBidReviewsForbiddenUrls = []string{
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=user6",
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/reviews?authorUsername=Organization 3&requesterUsername=user10",

	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user2&requesterUsername=user5",
	"/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/reviews?authorUsername=Organization 4&requesterUsername=user11",
}

// TestGetBidReviewsForbidden тестирует случай, когда у пользователя недостаточно прав
func TestGetBidReviewsForbidden(t *testing.T) {
	handleGetBidReviewsIncorrectUrls(t, getBidReviewsForbiddenUrls, http.StatusForbidden)
}

var getBidReviewsNotFoundUrls = []string{
	// тендер не найден
	"/bids/00000000-0000-0000-0000-000000000000/reviews?authorUsername=user12&requesterUsername=user11",
	// тендер, не связанный ни с одним предложением
	"/bids/1ffac8e1-42d3-4351-a778-751b9dacf4b0/reviews?authorUsername=user12&requesterUsername=user11",
	// отзывы не найдены
	"/bids/67309f56-3d5f-45ee-873c-11262ca16543/reviews?authorUsername=user4&requesterUsername=user11",
}

// TestCreateBidNotFound тестирует случай, когда тендер или отзывы не найдены
func TestGetBidReviewsNotFound(t *testing.T) {
	handleGetBidReviewsIncorrectUrls(t, getBidReviewsNotFoundUrls, http.StatusNotFound)
}

var getBidReviewsBadRequestUrls = []string{
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb50c523105-30f3-4f0b-b3b2-a23775ea9eb50c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=user5",
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=&requesterUsername=user5",
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=",
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=user5&limit=t&offset=2",
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=user5&limit=-1&offset=2",
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=user5&limit=51&offset=2",
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=user5&limit=3&offset=t",
	"/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=user5&limit=3&offset=-2",
}

// TestGetBidReviewsBadRequest тестирует некорректный запрос
func TestGetBidReviewsBadRequest(t *testing.T) {
	handleGetBidReviewsIncorrectUrls(t, getBidReviewsBadRequestUrls, http.StatusBadRequest)
}

type getBidReviewsOKTest struct {
	url              string
	expectedQuantity int
}

var getBidReviewsOKTests = []getBidReviewsOKTest{
	{url: "/bids/0c523105-30f3-4f0b-b3b2-a23775ea9eb5/reviews?authorUsername=user1&requesterUsername=user5", expectedQuantity: 3},
	{url: "/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/reviews?authorUsername=user8&requesterUsername=user11", expectedQuantity: 5},
	{url: "/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/reviews?authorUsername=user8&requesterUsername=user11&limit=2", expectedQuantity: 2},
	{url: "/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/reviews?authorUsername=user8&requesterUsername=user11&limit=3&offset=1", expectedQuantity: 3},
	{url: "/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/reviews?authorUsername=user8&requesterUsername=user11&&offset=3", expectedQuantity: 2},
	{url: "/bids/ca402b26-1bc6-4f97-b41a-0373aec3daed/reviews?authorUsername=Organization 3&requesterUsername=user11", expectedQuantity: 3},
}

// TestGetBidReviewsOK тестирует успешный просмотр отзывов на прошлые предложения
func TestGetBidReviewsOK(t *testing.T) {
	processGetBidReviewsCorrectUrls(t, getBidReviewsPath, getBidReviewsHttpMethod, getBidReviewsMethod)
}

func processGetBidReviewsCorrectUrls(t *testing.T, path, httpMethod, method string) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	for idx := 0; idx < len(getBidReviewsOKTests); idx++ {
		url := getBidReviewsOKTests[idx].url
		rr := test.ProcessReq(t, dtb, nil, url, path, httpMethod, method)
		test.CheckCodeAndMime(t, rr)

		var rws []bid.BidReview
		decoder := json.NewDecoder(rr.Body)
		err := decoder.Decode(&rws)
		if err != nil {
			t.Fatalf("Ошибка десериализации тела ответа сервера: %v", err)
		}

		expectedQuantity := getBidReviewsOKTests[idx].expectedQuantity
		quantity := len(rws)
		if quantity != expectedQuantity {
			t.Fatalf("Ожидалось количество отзывов: %d, но получено: %d", expectedQuantity, quantity)
		}
	}
}

func handleGetBidReviewsIncorrectUrls(t *testing.T, urls []string, expectedCode int) {
	test.ProcessIncorrectUrls(t, getBidReviewsPath, getBidReviewsHttpMethod, getBidReviewsMethod, urls, expectedCode)
}
