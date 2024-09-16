package bid

import (
	"database/sql"
)

type BidRepo interface {
	CreateBid(bci BidCreationInput) (*Bid, int, error)
	EditBid(bdi BidEditionInput, bidID, username string) (*Bid, int, error)
	SubmitBidDecision(bsi BidSubmissionInput) (*Bid, int, error)
	SubmitBidFeedback(bfi BidFeedbackInput) (*Bid, int, error)
	GetBidReviews(bri BidReviewsInput) ([]*BidReview, int, error)
}

type BidDBRepository struct {
	dtb *sql.DB
}

func NewDBRepo(sdb *sql.DB) *BidDBRepository {
	return &BidDBRepository{dtb: sdb}
}
