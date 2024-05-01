package models

// ROLE : 0 - admin, 1 - orginiser, 2 - distributor, 3 - buyer
type User struct {
	Id      int64
	Login   string
	Email   string
	PwdHash []byte
	Role    int32
}
