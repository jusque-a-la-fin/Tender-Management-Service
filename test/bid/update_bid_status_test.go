package bid_test

import (
	"log"
	"net/http"
	"tendermanagement/internal/datastore"
	"tendermanagement/test"
	"testing"
)

// TestUpdateBidStatus тестирует изменение статуса предложения по его уникальному идентификатору
func TestUpdateBidStatus(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	url := "/bids/5820469e-0a87-43d2-b139-0eb5e253cbfa/status?status=Published&username=user14"
	path := "/bids/{bidId}/status"
	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodPut, "UpdateBidStatus")
	test.HandleBidResponse(t, rr)
}
