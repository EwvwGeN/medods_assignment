package models

type RefreshRequest struct {
	TokenPair TokenPair `json:"token_pair"`
}

type RefreshResponse struct {
	TokenPair TokenPair `json:"token_pair"`
}