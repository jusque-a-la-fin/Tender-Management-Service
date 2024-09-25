package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"tendermanagement/internal/bid"
	"tendermanagement/internal/datastore"
	"tendermanagement/internal/handlers"
	bhd "tendermanagement/internal/handlers/bid"
	thd "tendermanagement/internal/handlers/tender"
	"tendermanagement/internal/tender"
	"testing"

	"github.com/gorilla/mux"
)

func SetVars() {
	err := os.Setenv("POSTGRES_USERNAME", "postgres")
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}
	err = os.Setenv("POSTGRES_PASSWORD", "bmw")
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}
	err = os.Setenv("POSTGRES_HOST", "localhost")
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}
	err = os.Setenv("POSTGRES_PORT", "5432")
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}
	err = os.Setenv("POSTGRES_DATABASE", "tenders")
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return
	}
}

func ProcessReq(t *testing.T, dtb *sql.DB, body any, url, path, httpMethod, method string) *httptest.ResponseRecorder {
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(data))
	if err != nil {
		t.Fatal("Ошибка создания объекта *http.Request:", err)
	}

	if dtb == nil {
		dtb, err = datastore.CreateNewDB()
		if err != nil {
			log.Fatalf("ошибка подключения к базе данных: %v", err)
		}
	}

	var handlerFunc func(http.ResponseWriter, *http.Request)

	var rr *httptest.ResponseRecorder
	tenderMethods := []string{"CreateTender", "EditTender", "GetTenderStatus", "GetUserTenders", "RollbackTender", "UpdateTenderStatus"}
	bidMethods := []string{"CreateBid", "EditBid", "UpdateBidStatus", "RollbackBid", "SubmitBidDecision", "SubmitBidFeedback"}

	if contains(tenderMethods, method) {
		tenderHandler := getTenderHandler(dtb)
		switch method {
		case "CreateTender":
			handlerFunc = tenderHandler.CreateTender
			rr = ServeRequest(handlerFunc, req)

		case "EditTender":
			handlerFunc = tenderHandler.EditTender
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)

		case "GetTenderStatus":
			handlerFunc = tenderHandler.GetTenderStatus
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)

		case "UpdateTenderStatus":
			handlerFunc = tenderHandler.UpdateTenderStatus
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)

		case "GetUserTenders":
			handlerFunc = tenderHandler.GetUserTenders
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)

		case "RollbackTender":
			handlerFunc = tenderHandler.RollbackTender
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)
		}
	}

	if contains(bidMethods, method) {
		bidHandler := GetBidHandler(dtb)
		switch method {
		case "CreateBid":
			handlerFunc = bidHandler.CreateBid
			rr = ServeRequest(handlerFunc, req)

		case "EditBid":
			handlerFunc = bidHandler.EditBid
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)

		case "RollbackBid":
			handlerFunc = bidHandler.RollbackBid
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)

		case "SubmitBidDecision":
			handlerFunc = bidHandler.SubmitBidDecision
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)

		case "SubmitBidFeedback":
			handlerFunc = bidHandler.SubmitBidFeedback
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)

		case "UpdateBidStatus":
			handlerFunc = bidHandler.UpdateBidStatus
			rr = setupRouterAndServe(path, httpMethod, handlerFunc, req)
		}
	}

	return rr
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func getTenderHandler(dtb *sql.DB) *thd.TenderHandler {
	tdr := tender.NewDBRepo(dtb)
	tenderHandler := &thd.TenderHandler{
		TenderRepo: tdr,
	}
	return tenderHandler
}

func GetBidHandler(dtb *sql.DB) *bhd.BidHandler {
	nbd := bid.NewDBRepo(dtb)
	bidHandler := &bhd.BidHandler{
		BidRepo: nbd,
	}
	return bidHandler
}

func ServeRequest(handlerFunc func(http.ResponseWriter, *http.Request), req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rr, req)
	return rr
}

func setupRouterAndServe(path, httpMethod string, handlerFunc func(http.ResponseWriter, *http.Request), req *http.Request) *httptest.ResponseRecorder {
	rtr := mux.NewRouter()
	rtr.HandleFunc(path, handlerFunc).Methods(httpMethod)
	rr := httptest.NewRecorder()
	rtr.ServeHTTP(rr, req)
	return rr
}

func HandleError(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	if mime := rr.Header().Get("Content-Type"); mime != "application/json" {
		t.Errorf("Заголовок Content-Type должен иметь MIME-тип application/json, но имеет %s", mime)
	}

	var errResp handlers.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errResp)
	if err != nil {
		t.Fatalf("Ошибка десериализации тела ответа сервера: %v", err)
	}

	result := errResp.Reason
	if result != expected {
		t.Errorf("Ожидалось %s, но получено %s", expected, result)
	}
}

func HandleBadReq(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	code := rr.Code
	if code != http.StatusBadRequest {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusBadRequest, code)
	}

	HandleError(t, rr, expected)
}

func CheckCodeAndMime(t *testing.T, rr *httptest.ResponseRecorder) {
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, rr.Code)
	}

	if mime := rr.Header().Get("Content-Type"); mime != "application/json" {
		t.Errorf("Заголовок Content-Type должен иметь MIME-тип application/json, но имеет %s", mime)
	}
}
