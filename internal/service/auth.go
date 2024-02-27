package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/EwvwGeN/medods_assignment/internal/domain/models"
	guuid "github.com/google/uuid"
)

type Auth struct {
	log *slog.Logger
	userRepo UserRepo
	jwtManager JwtManager
	tokenTTL time.Duration

}

type UserRepo interface {
	SaveUser(ctx context.Context, email, uuid string) (err error)
	GetUserByRefresh(ctx context.Context, refresh string) (user *models.User, err error)
	SaveRefreshToken(ctx context.Context, uuid, refresh string) (err error)
	GetUserByUUID(ctx context.Context, uuid string) (user *models.User, err error)
}

type JwtManager interface {
	CreateJwt(user *models.User, ttl time.Duration) (token string, err error)
	CreateRefresh() (refresh string, err error)
}

func NewAuth(ctx context.Context, log *slog.Logger, userRepo UserRepo, jwtManager JwtManager, tokenttl time.Duration) *Auth {
	return &Auth{
		log: log,
		userRepo: userRepo,
		jwtManager: jwtManager,
		tokenTTL: tokenttl,
	}
}

func (a *Auth) RegisterUser(email string) (uuid string, err error) {
	log := a.log.With(slog.String("auth.method", "register_user"))
	uuid = guuid.New().String()
	if uuid == "" {
		log.Error(ErrCreateUUID.Error())
		return "", fmt.Errorf("cant register user: %w", ErrCreateUUID)
	}
	err = a.userRepo.SaveUser(context.Background(), email, uuid)
	if err != nil {
		log.Warn("cant register user", slog.String("error", err.Error()))
		return "", ErrSaveUser
	}
	return
}

func (a *Auth) CreateTokenPair(uuid string) (token, refresh string, err error) {
	log := a.log.With(slog.String("auth.method", "create_token_pair"))
	user, err:= a.userRepo.GetUserByUUID(uuid)
	if err != nil {
		log.Warn(ErrGetUserUUID.Error(), slog.String("error", err.Error()))
		return "", "", ErrGetUserUUID
	}
	token, err = a.jwtManager.CreateJwt(user, a.tokenTTL)
	if err != nil {
		log.Error(ErrCreateJWT.Error(), slog.String("error", err.Error()))
		return "", "", ErrCreateJWT
	}
	refresh, err = a.jwtManager.CreateRefresh()
	if err != nil {
		log.Error(ErrCreateRefresh.Error(), slog.String("error", err.Error()))
		return "", "", ErrCreateRefresh
	}
	return
}

func (a *Auth) RefreshToken(oldRefresh string) (newToken, newRefresh string, err error) {
	log := a.log.With(slog.String("auth.method", "refresh_token"))
	user, err := a.userRepo.GetUserByRefresh(oldRefresh)
	if err != nil {
		log.Warn(ErrGetUserRefresh.Error(), slog.String("error", err.Error()))
		return "", "", ErrGetUserRefresh
	}
	newToken, err = a.jwtManager.CreateJwt(user, a.tokenTTL)
	if err != nil {
		log.Error(ErrCreateJWT.Error(), slog.String("error", err.Error()))
		return "", "", ErrCreateJWT
	}
	newRefresh, err = a.jwtManager.CreateRefresh()
	if err != nil {
		log.Error(ErrCreateRefresh.Error(), slog.String("error", err.Error()))
		return "", "", ErrCreateRefresh
	}
	return
}