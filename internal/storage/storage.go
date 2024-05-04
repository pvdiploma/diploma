package storage

import "errors"

var (
	ErrUserExists    = errors.New("user already exist")
	ErrUserNotFound  = errors.New("user not found")
	ErrAppNotFound   = errors.New("app not found")
	ErrEventNotFound = errors.New("event not found")
	ErrEventExists   = errors.New("event already exist")
)

// implement that in proto files
const (
	DefaultEmptyInt = -1
	DefaultEmptyStr = "None"
)
