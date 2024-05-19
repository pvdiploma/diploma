package models

type DealStatus string

const (
	Accepted DealStatus = "accepted"
	Rejected DealStatus = "rejected"
	Pending  DealStatus = "pending"
)

type Deal struct {
	ID            int64
	SenderID      int64
	RecipientID   int64
	OrganizerID   int64
	DistributorID int64
	EventID       int64
	Commission    uint32
	Status        DealStatus
}

type Widget struct {
	ID     int64
	DealID int64
	Body   string
	Script string
}
