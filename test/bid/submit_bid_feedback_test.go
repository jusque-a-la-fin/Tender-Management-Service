package bid_test

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"tendermanagement/internal/datastore"
	"tendermanagement/test"
	"testing"
)

// TestSubmitBidFeedbackOK тестирует успешную отправку отзыва по предложению
func TestSubmitBidFeedbackOK(t *testing.T) {
	test.SetVars()
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	bidId := "3820469e-0a87-43d2-b139-0eb5e253cbfa"
	username := "user11"
	expectedReview := "отзыв"
	url := fmt.Sprintf("/bids/%s/feedback?bidFeedback=%s&username=%s", bidId, expectedReview, username)
	path := "/bids/{bidId}/feedback"

	rr := test.ProcessReq(t, dtb, nil, url, path, http.MethodPut, "SubmitBidFeedback")
	bid := test.HandleBidResponse(t, rr)

	var userID string
	dtb.QueryRow("SELECT id FROM employee WHERE username = $1", username).Scan(&userID)

	var review string
	query := "SELECT description FROM bid_review WHERE bid_id = $1 AND user_id = $2;"
	dtb.QueryRow(query, bid.ID, userID).Scan(&review)

	if review != expectedReview {
		t.Errorf("Ожидался отзыв по предложению: %s, но получен: %s", expectedReview, review)
	}

	deleteBidReview(t, dtb)
}

func deleteBidReview(t *testing.T, dtb *sql.DB) {
	query := `DELETE FROM bid_review 
	          WHERE id = 
			  (SELECT id FROM bid_review
               ORDER BY created_at DESC LIMIT 1);`

	result, err := dtb.Exec(query)
	if err != nil {
		t.Fatalf("Ошибка удаления предложения: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected == 0 {
		t.Fatalf("Ошибка: предложение не было удалено: %v", err)
	}
}
