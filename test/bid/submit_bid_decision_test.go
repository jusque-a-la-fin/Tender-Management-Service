package bid_test

import (
	"log"
	"net/http"
	"tendermanagement/internal/bid"
	"tendermanagement/internal/datastore"
	"tendermanagement/test"
	"testing"
)

type SubmitDecisionResult struct {
	url       string
	decision  bid.DecisionEnum
	bidStatus bid.StatusEnum
}

var sdtests = []SubmitDecisionResult{
	{
		url:       "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Approved&username=user11",
		decision:  "Approved",
		bidStatus: "Canceled",
	},
	{
		url:       "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Rejected&username=user11",
		decision:  "Rejected",
		bidStatus: "Canceled",
	},
}

// TestSubmitBidDecisionOK тестирует принятия решения по предложению
func TestSubmitBidDecisionOK(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	for _, sdtest := range sdtests {
		path := "/bids/{bidId}/submit_decision"
		rr := test.ProcessReq(t, dtb, nil, sdtest.url, path, http.MethodPut, "SubmitBidDecision")
		bid := handleBidResponse(t, rr)
		status := sdtest.bidStatus
		expectedStatus := bid.Status
		if status != expectedStatus {
			t.Errorf("Ожидалось значение Status у предложения: %s, но получено: %s", expectedStatus, status)
		}

		var decision string
		query := "SELECT decision FROM bid_decisions WHERE bid_id = $1;"
		_ = dtb.QueryRow(query, bid.ID).Scan(&decision)

		expectedDecision := sdtest.decision
		if status != expectedStatus {
			t.Errorf("Ожидалось решение по предложению: %s, но получено: %s", expectedDecision, decision)
		}

		var tenderID string
		query = "SELECT tender_id FROM bid WHERE id = $1;"
		_ = dtb.QueryRow(query, bid.ID).Scan(&tenderID)

		var tenderStatus string
		query = "SELECT status FROM tender WHERE id = $1;"
		_ = dtb.QueryRow(query, tenderID).Scan(&tenderStatus)

		expectedTenderStatus := "Closed"
		if decision == "Approved" && tenderStatus != expectedTenderStatus {
			t.Errorf("Ожидался статус тендера: %s, но получен: %s", expectedTenderStatus, tenderStatus)
		}
	}
}
