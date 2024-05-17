package models

type PurchaseInfo struct {
	WidgetID        int64
	Name            string
	Surname         string
	Patronymic      string
	Email           string
	Phone           string
	PurchaseTickets []PurchaseTickets
}

type PurchaseTickets struct {
	ID     int64
	Amount uint32
}
