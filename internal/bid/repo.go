package bid

import (
	"database/sql"
)

type BidRepo interface {
	CreateBid(bci BidCreationInput) (*Bid, int, error)
	GetUserBids(username string, limit, offset int32) ([]Bid, int, error)
	GetBidsForTender(gbi GetBidsInput) ([]*Bid, int, error)
	EditBid(bdi BidEditionInput, bidID, username string) (*Bid, int, error)
	GetBidStatus(bidID, username string) (StatusEnum, int, error)
	UpdateBidStatus(bidID, username string, newStatus StatusEnum) (*Bid, int, error)
	AddBidDecisions(bsi BidSubmissionInput) (string, string, int, error)
	MakeFinalDecision(bidID, tenderID, organizationID string) (*Bid, int, error)
	SubmitBidFeedback(bfi BidFeedbackInput) (*Bid, int, error)
	GetBidReviews(bri BidReviewsInput) ([]BidReview, int, error)
	RollbackBid(bri BidRollbackInput) (*Bid, int, error)
}

type BidDBRepository struct {
	dtb *sql.DB
}

func NewDBRepo(sdb *sql.DB) *BidDBRepository {
	return &BidDBRepository{dtb: sdb}
}
