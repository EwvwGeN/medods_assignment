package models

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	TokenPair TokenPair `json:"token_pair"`
}