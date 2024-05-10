package models

type Ticket struct {
	ID              int64
	EventCategoryID int64
	Name            string
	Surname         string
	Patronymic      string
	Email           string
	Discount        uint32
	Total           uint32
	QRCode          []byte
	IsActivated     bool
	ImageBytes      []byte
	ImagePath       string
}
