package types

type DecodedToken struct {
	Username string `json:"username"`
	Exp int64 `json:"exp"`
}

type User struct {
	Username string
	Password string
}