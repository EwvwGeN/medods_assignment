package storage

import "errors"

var (
	ErrDbNotExist   = errors.New("database does not exist")
	ErrCollNotExist = errors.New("collection does not exist")
	ErrUserExist    = errors.New("user already exist")
	ErrUserNotFound = errors.New("user not found")
	ErrUpdate = errors.New("error while update")
)