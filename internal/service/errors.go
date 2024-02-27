package service

import "errors"

var (
	ErrCreateUUID = errors.New("error while creating uuid")
	ErrSaveUser = errors.New("error while saving user")
	ErrGetUserUUID = errors.New("cant get user by uuid")
	ErrGetUserRefresh = errors.New("cant get user by refresh token")
	ErrCreateJWT = errors.New("error while creating jwt")
	ErrCreateRefresh = errors.New("error while creating refresh token")
)