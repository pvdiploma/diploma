package storage

import "errors"

var (
	ErrUserExists     = errors.New("user already exist")
	ErrUserNotFound   = errors.New("user not found")
	ErrAppNotFound    = errors.New("app not found")
	ErrEventNotFound  = errors.New("event not found")
	ErrEventExists    = errors.New("event already exist")
	ErrTicketExists   = errors.New("ticket already exists")
	ErrTicketNotFound = errors.New("ticket not found")
)

// implement that in proto files
const (
	DefaultEmptyInt = 0 // change to -1 after fixing proto
	DefaultEmptyStr = ""
)
