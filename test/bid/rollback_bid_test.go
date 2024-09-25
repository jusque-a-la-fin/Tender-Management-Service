package bid_test

import (
	"log"
	"net/http"
	"tendermanagement/internal/datastore"
	"tendermanagement/test"
	"testing"
)

var testRollbackUrls = []string{
	"/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/rollback/1?username=user12",
	"/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/rollback/2?username=user12",
}

// TestRollbackBidOK тестирует успешный откат параметров тендера к указанной версии
func TestRollbackBidOK(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	path := "/bids/{bidId}/rollback/{version}"

	for _, testUrl := range testRollbackUrls {
		rr := test.ProcessReq(t, dtb, nil, testUrl, path, http.MethodPut, "RollbackBid")
		if rr.Code != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, rr.Code)
		}
	}
}
