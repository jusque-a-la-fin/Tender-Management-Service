package bid_test

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"tendermanagement/internal/bid"
	"tendermanagement/internal/datastore"
	"tendermanagement/internal/tender"
	"tendermanagement/test"
	"testing"
)

type SubmitDecision struct {
	url      string
	decision bid.DecisionEnum
	userId   string
}

type SubmitDecisionTestOK struct {
	suds         []SubmitDecision
	tenderStatus tender.StatusEnum
}

var exampleTenderId = "67309f56-3d5f-45ee-873c-11262ca16543"

var sdtests = []SubmitDecisionTestOK{
	{suds: []SubmitDecision{{url: "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Approved&username=user11",
		decision: "Approved", userId: "550e8400-e29b-41d4-a716-44665544000b"},
		{url: "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Approved&username=user12",
			decision: "Approved", userId: "550e8400-e29b-41d4-a716-44665544000c"}},
		tenderStatus: "Created",
	},

	{suds: []SubmitDecision{{url: "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Approved&username=user11",
		decision: "Approved", userId: "550e8400-e29b-41d4-a716-44665544000b"},
		{url: "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Approved&username=user12",
			decision: "Approved", userId: "550e8400-e29b-41d4-a716-44665544000c"},
		{url: "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Approved&username=user13",
			decision: "Approved", userId: "550e8400-e29b-41d4-a716-44665544000d"}},
		tenderStatus: "Closed",
	},

	{suds: []SubmitDecision{{url: "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Approved&username=user11",
		decision: "Approved", userId: "550e8400-e29b-41d4-a716-44665544000b"},
		{url: "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Approved&username=user12",
			decision: "Approved", userId: "550e8400-e29b-41d4-a716-44665544000c"},
		{url: "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Approved&username=user13",
			decision: "Approved", userId: "550e8400-e29b-41d4-a716-44665544000d"},
		{url: "/bids/3820469e-0a87-43d2-b139-0eb5e253cbfa/submit_decision?decision=Rejected&username=user14",
			decision: "Rejected", userId: "550e8400-e29b-41d4-a716-44665544000e"}},
		tenderStatus: "Created",
	},
}

var path = "/bids/{bidId}/submit_decision"

// TestSubmitBidDecisionOK тестирует принятие решения по предложению
func TestSubmitBidDecisionOK(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	for _, sdtest := range sdtests {
		var wg sync.WaitGroup

		for _, sud := range sdtest.suds {
			wg.Add(1)

			go func(sud SubmitDecision) {
				defer wg.Done()

				rr := test.ProcessReq(t, dtb, nil, sud.url, path, http.MethodPut, "SubmitBidDecision")
				nbid := test.HandleBidResponse(t, rr)

				var decision string
				query := "SELECT decision FROM bid_decisions WHERE bid_id = $1 AND user_id = $2;"
				_ = dtb.QueryRow(query, nbid.ID, sud.userId).Scan(&decision)

				expectedDecision := sud.decision
				if bid.DecisionEnum(decision) != expectedDecision {
					t.Errorf("Ожидалось решение по предложению: %s, но получено: %s", expectedDecision, decision)
				}
			}(sud)
		}

		wg.Wait()

		var tenderStatus tender.StatusEnum
		query := "SELECT status FROM tender WHERE id = $1;"
		_ = dtb.QueryRow(query, exampleTenderId).Scan(&tenderStatus)

		expectedTenderStatus := sdtest.tenderStatus
		if tenderStatus != expectedTenderStatus {
			t.Errorf("Ожидался статус тендера: %s, но получен: %s", expectedTenderStatus, tenderStatus)
		}

		rollbackDecision(t, dtb)
	}
}

func rollbackDecision(t *testing.T, dtb *sql.DB) {
	query := `
	UPDATE tender
	SET status = 'Created'
	WHERE id = $1;`

	_, err := dtb.Exec(query, exampleTenderId)
	if err != nil {
		t.Errorf("ошибка запроса к базе данных: откат статуса тендера: %v", err)
	}

	query = `
	TRUNCATE TABLE bid_decisions;`

	_, err = dtb.Exec(query)
	if err != nil {
		t.Errorf("ошибка запроса к базе данных: очистка решений: %v", err)
	}

	var tenderStatus tender.StatusEnum
	query = "SELECT status FROM tender WHERE id = $1;"
	_ = dtb.QueryRow(query, exampleTenderId).Scan(&tenderStatus)
}
