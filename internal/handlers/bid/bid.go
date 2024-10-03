package bid

import (
	"sync"
	bids "tendermanagement/internal/bid"
)

type BidHandler struct {
	BidRepo      bids.BidRepo
	DecisionOnce sync.Once
}

// BidCreationReq - тело запроса для создания предложения для существующего тендера
type BidCreationReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderID    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorId    string `json:"authorId"`
}
