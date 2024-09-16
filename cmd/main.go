package main

import (
	"log"
	"net/http"
	"os"
	"tendermanagement/internal/bid"
	"tendermanagement/internal/datastore"
	"tendermanagement/internal/handlers"
	"tendermanagement/internal/tender"

	bdh "tendermanagement/internal/handlers/bid"
	thd "tendermanagement/internal/handlers/tender"

	"github.com/gorilla/mux"
)

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/ping", handlers.CheckServer).Methods("GET")

	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	tdr := tender.NewDBRepo(dtb)
	tendersHandler := &thd.TenderHandler{
		TenderRepo: tdr,
	}

	bdr := bid.NewDBRepo(dtb)
	bids := &bdh.BidHandler{
		BidRepo: bdr,
	}

	rtr.HandleFunc("/tenders", tendersHandler.GetTenders).Methods("GET")
	rtr.HandleFunc("/tenders/new", tendersHandler.CreateTender).Methods("POST")
	rtr.HandleFunc("/tenders/my", tendersHandler.GetUserTenders).Methods("GET")
	rtr.HandleFunc("/tenders/{tenderId}/status", tendersHandler.GetTenderStatus).Methods("GET")
	rtr.HandleFunc("/tenders/{tenderId}/status", tendersHandler.UpdateTenderStatus).Methods("PUT")
	rtr.HandleFunc("/tenders/{tenderId}/edit", tendersHandler.EditTender).Methods("PATCH")
	rtr.HandleFunc("/tenders/{tenderId}/rollback/{version}", tendersHandler.RollbackTender).Methods("PUT")
	rtr.HandleFunc("/bids/new", bids.CreateBid).Methods("POST")
	rtr.HandleFunc("/bids/{bidId}/edit", bids.EditBid).Methods("PATCH")
	rtr.HandleFunc("/bids/{bidId}/submit_decision", bids.SubmitBidDecision).Methods("PUT")
	rtr.HandleFunc("/bids/{bidId}/feedback", bids.SubmitBidFeedback).Methods("PUT")
	rtr.HandleFunc("/bids/{tenderId}/reviews", bids.GetBidReviews).Methods("GET")

	port := os.Getenv("SERVER_ADDRESS")

	certFile := "/etc/tls/tls.crt"
	keyFile := "/etc/tls/tls.key"
	certExists := false
	keyExists := false

	if _, err := os.Stat(certFile); err == nil {
		certExists = true
	} else if !os.IsNotExist(err) {
		log.Fatalf("ошибка проверки сертификата: %v", err)
	}

	if _, err := os.Stat(keyFile); err == nil {
		keyExists = true
	} else if !os.IsNotExist(err) {
		log.Fatalf("ошибка проверки ключа: %v", err)
	}

	if certExists && keyExists {
		err = http.ListenAndServeTLS(port, certFile, keyFile, rtr)
		if err != nil {
			log.Fatalf("ListenAndServe error: %#v", err)
		}
	} else {
		err = http.ListenAndServe(port, rtr)
		if err != nil {
			log.Fatalf("ListenAndServe error: %#v", err)
		}
	}
}
