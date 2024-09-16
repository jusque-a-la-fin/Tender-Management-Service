package bid

import (
	bids "tendermanagement/internal/bid"
)

type BidHandler struct {
	BidRepo bids.BidRepo
}

// BidCreationReq - тело запроса для создания предложения для существующего тендера
type BidCreationReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderId    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorId    string `json:"authorId"`
}
