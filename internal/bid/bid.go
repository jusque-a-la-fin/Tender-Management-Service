package bid

import (
	"database/sql"
	"fmt"
)

// Bid - предложение
type Bid struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      StatusEnum     `json:"status"`
	TenderID    string         `json:"tenderId"`
	AuthorType  AuthorTypeEnum `json:"authorType"`
	AuthorID    string         `json:"authorId"`
	Version     int32          `json:"version"`
	CreatedAt   string         `json:"createdAt"`
}

type AuthorTypeEnum string

const (
	Organization AuthorTypeEnum = "Organization"
	User         AuthorTypeEnum = "User"
)

func CheckAuthorType(authorType string) bool {
	fail := true
	switch authorType {
	case "Organization":
		fail = false
	case "User":
		fail = false
	}
	return fail
}

// BidReview - отзыв о предложении
type BidReview struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type BidReviewsInput struct {
	TenderId          string
	AuthorUsername    string
	RequesterUsername string
	Offset            int32
	EndIndex          int32
}

type GetBidsInput struct {
	TenderId string
	Username string
	Offset   int32
	EndIndex int32
}

type BidRollbackInput struct {
	BidID    string
	Version  int32
	Username string
}

func CheckDecision(decision string) bool {
	fail := true
	switch decision {
	case "Approved":
		fail = false
	case "Rejected":
		fail = false
	}
	return fail
}

type StatusEnum string

const (
	Created   StatusEnum = "Created"
	Published StatusEnum = "Published"
	Canceled  StatusEnum = "Canceled"
)

func CheckStatusEnum(status string) bool {
	fail := true
	switch status {
	case "Created":
		fail = false
	case "Published":
		fail = false
	case "Canceled":
		fail = false
	}
	return fail
}

// BidCreationInput - параметры, переданные пользователем серверу для создания предложения
type BidCreationInput struct {
	Name        string
	Description string
	TenderID    string
	AuthorType  AuthorTypeEnum
	AuthorId    string
}

// BidEditionInput - параметры, переданные пользователем серверу для редактирования предложения
type BidEditionInput struct {
	Name        string
	Description string
}

type DecisionEnum string

const (
	Approved DecisionEnum = "Approved"
	Rejected DecisionEnum = "Rejected"
)

// BidEditionInput - параметры, переданные пользователем серверу для отправки решения (одобрения или отклонения) по предложению
type BidSubmissionInput struct {
	BidID    string
	Decision DecisionEnum
	Username string
}

// BidFeedbackInput - параметры, переданные пользователем серверу для отправки отзыва по предложению
type BidFeedbackInput struct {
	BidID       string
	BidFeedback string
	Username    string
}

func GetBid(dtb *sql.DB, bidID string) (*Bid, error) {
	nbd := Bid{
		ID: bidID,
	}

	query := `
		SELECT bd.status, bd.tender_id, bd.author_type, bd.author_id ,bd.current_version, bd.created_at, 
			   bvs.name, bvs.description
		FROM bid bd
		JOIN bid_versions bvs ON bd.id = bvs.bid_id AND bd.current_version = bvs.version
		WHERE bd.id = $1
	`

	err := dtb.QueryRow(query, bidID).Scan(&nbd.Status, &nbd.TenderID, &nbd.AuthorType,
		&nbd.AuthorID, &nbd.Version, &nbd.CreatedAt, &nbd.Name, &nbd.Description)

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: извлечение параметров тендера: %v", err)
	}

	return &nbd, nil
}
