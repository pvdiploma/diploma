package models

import "time"

type Event struct {
	ID           int64
	OwnerID      int64
	Name         string
	Description  string
	Country      string
	City         string
	Place        string
	Address      string
	Date         time.Time
	TicketAmount uint32 `gorm:"omitempty"`
	Age          string
	Categories   []EventCategory
}

type EventCategory struct {
	ID       int64
	EventID  int64
	Category string
	Price    uint32
	Amount   uint32 `gorm:"omitempty"`
}
