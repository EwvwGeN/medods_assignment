package service

import "errors"

var (
	ErrCreateUUID = errors.New("error while creating uuid")
	ErrSaveUser = errors.New("error while saving user")
	ErrGetUserUUID = errors.New("cant get user by uuid")
	ErrCreateJWT = errors.New("error while creating jwt")
	ErrValidAccess = errors.New("access token doesnt valid")
	ErrValidRefresh = errors.New("refresh token doesnt valid")
	ErrCreateRefresh = errors.New("error while creating refresh token")
)