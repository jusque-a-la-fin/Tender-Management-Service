package bid_test

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"tendermanagement/internal/bid"
	"tendermanagement/internal/datastore"
	"tendermanagement/test"
	"testing"
)

// TestEditBidOK1 тестирует успешное редактирование параметров предложения
func TestEditTenderOK1(t *testing.T) {
	test.SetVars()
	body := bid.BidEditionInput{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
	}

	rr := getBidEditionResponseRecorder(t, body)

	bid := test.HandleBidResponse(t, rr)
	name := bid.Name
	bName := body.Name
	if name != bName {
		t.Errorf("Ожидалось значение Name: %s, но получено: %s", bName, name)
	}

	description := bid.Description
	bDescription := body.Description
	if description != bDescription {
		t.Errorf("Ожидалось значение Name: %s, но получено: %s", bDescription, description)
	}
}

// TestEditBidOK2 тестирует успешное редактирование полного названия предложения
func TestEditTenderOK2(t *testing.T) {
	test.SetVars()
	body := bid.BidEditionInput{
		Name:        "Доставка товары Казань - Москва",
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
	}

	rr := getBidEditionResponseRecorder(t, body)

	bid := test.HandleBidResponse(t, rr)
	name := bid.Name
	bName := body.Name
	if name != bName {
		t.Errorf("Ожидалось значение Name: %s, но получено: %s", bName, name)
	}
}

// TestEditBidOK2 тестирует успешное редактирование описания предложения
func TestEditTenderOK3(t *testing.T) {
	test.SetVars()
	body := bid.BidEditionInput{
		Description: "Нужно доставить оборудование для олимпиады по робототехнике",
	}

	rr := getBidEditionResponseRecorder(t, body)

	bid := test.HandleBidResponse(t, rr)
	description := bid.Description
	bDescription := body.Description
	if description != bDescription {
		t.Errorf("Ожидалось значение Name: %s, но получено: %s", bDescription, description)
	}
}

func getBidParams(dtb *sql.DB) (string, string) {
	var bidID, authorType, authorID string
	query := "SELECT id, author_type, author_id FROM bid LIMIT 1;"
	_ = dtb.QueryRow(query).Scan(&bidID, &authorType, &authorID)

	var username string
	switch authorType {
	case "User":
		query = "SELECT username FROM employee WHERE id = $1;"
		_ = dtb.QueryRow(query, authorID).Scan(&username)

	case "Organization":
		query = "SELECT name FROM organization WHERE id = $1;"
		_ = dtb.QueryRow(query, authorID).Scan(&username)
	}

	return bidID, username
}

func getBidEditionResponseRecorder(t *testing.T, body any) *httptest.ResponseRecorder {
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	bidID, username := getBidParams(dtb)
	return receiveBidEditionResponseRecorder(t, dtb, body, bidID, username)
}

func receiveBidEditionResponseRecorder(t *testing.T, dtb *sql.DB, body any, bidID, username string) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/bids/%s/edit?username=%s", bidID, username)

	path := "/bids/{bidId}/edit"
	return test.ProcessReq(t, dtb, body, url, path, http.MethodPatch, "EditBid")
}
