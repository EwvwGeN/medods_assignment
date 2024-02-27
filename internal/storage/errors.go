package storage

import "errors"

var (
	ErrDbNotExist   = errors.New("database does not exist")
	ErrCollNotExist = errors.New("collection does not exist")
	ErrUserExist    = errors.New("user already exist")
	ErrUserNotFound = errors.New("user not found")
	ErrRefresh = errors.New("refresh token doesnt valid")
	ErrUpdate = errors.New("error whiel update")
)