package models

type RegisterRequest struct {
	Email string `json:"email"`
}

type RegisterResponse struct {
	UUID string `json:"uuid"`
}