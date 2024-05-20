package storage

import "errors"

var (
	ErrUserExists         = errors.New("user already exist")
	ErrUserNotFound       = errors.New("user not found")
	ErrAppNotFound        = errors.New("app not found")
	ErrEventNotFound      = errors.New("event not found")
	ErrEventExists        = errors.New("event already exist")
	ErrTicketExists       = errors.New("ticket already exists")
	ErrTicketNotFound     = errors.New("ticket not found")
	ErrDealNotFound       = errors.New("deal not found")
	ErrDealExists         = errors.New("deal already exists")
	ErrDealWidgetNotFound = errors.New("deal not found")
)

const (
	DefaultEmptyInt = 0
	DefaultEmptyStr = ""
)
