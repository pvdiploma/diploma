package models

type User struct {
	UUID    int64
	Login   string
	Email   string
	PwdHash []byte
	Role    int32
}
