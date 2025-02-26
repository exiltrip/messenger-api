package models

type User struct {
	ID           int    `db:"id"`
	Email        string `db:"email"`
	FullName     string `db:"fullname"`
	PasswordHash string `db:"password_hash"`
	PeerID       int    `db:"peer_id"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
