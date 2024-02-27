package jwt

import "errors"

var (
	ErrEmptyValue = errors.New("empty value")
	ErrRefreshGenerate = errors.New("cant generate refresh token")
)