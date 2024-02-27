package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/EwvwGeN/medods_assignment/internal/domain/models"
	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log *slog.Logger
	userRepo UserRepo
	jwtManager JwtManager
	tokenTTL time.Duration
	refreshTTL time.Duration

}

type UserRepo interface {
	SaveUser(ctx context.Context, email, uuid string) (err error)
	GetUserByRefresh(ctx context.Context, refresh []byte) (user *models.User, err error)
	SaveRefreshToken(ctx context.Context, uuid string, refresh []byte, refreshTTL time.Duration) (err error)
	GetUserByUUID(ctx context.Context, uuid string) (user *models.User, err error)
}

type JwtManager interface {
	CreateJwt(user *models.User, ttl time.Duration) (token string, err error)
	CreateRefresh() (refresh string, err error)
}

func NewAuth(ctx context.Context, log *slog.Logger, userRepo UserRepo, jwtManager JwtManager, tokenttl, refreshttl time.Duration) *Auth {
	return &Auth{
		log: log,
		userRepo: userRepo,
		jwtManager: jwtManager,
		tokenTTL: tokenttl,
		refreshTTL: refreshttl,
	}
}

func (a *Auth) RegisterUser(email string) (uuid string, err error) {
	log := a.log.With(slog.String("auth.method", "register_user"))
	uuid = guuid.New().String()
	if uuid == "" {
		log.Error(ErrCreateUUID.Error())
		return "", fmt.Errorf("failed register user: %w", ErrCreateUUID)
	}
	err = a.userRepo.SaveUser(context.Background(), email, uuid)
	if err != nil {
		log.Warn("failed register user", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed register user: %w", ErrSaveUser)
	}
	return
}

func (a *Auth) CreateTokenPair(uuid string) (token, refresh string, err error) {
	log := a.log.With(slog.String("auth.method", "create_token_pair"))
	user, err:= a.userRepo.GetUserByUUID(context.Background(), uuid)
	if err != nil {
		log.Warn(ErrGetUserUUID.Error(), slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed create token pair: %w", ErrGetUserUUID)
	}
	token, err = a.jwtManager.CreateJwt(user, a.tokenTTL)
	if err != nil {
		log.Error(ErrCreateJWT.Error(), slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed create token pair: %w", ErrCreateJWT)
	}
	refresh, err = a.jwtManager.CreateRefresh()
	if err != nil {
		log.Error(ErrCreateRefresh.Error(), slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed create token pair: %w", ErrCreateRefresh)
	}
	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate refresh hash", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed create token pair: %w", err)
	}
	err = a.userRepo.SaveRefreshToken(context.Background(), user.UUID, refreshHash, a.refreshTTL)
	if err != nil {
		log.Error("failed to save refresh token")
		return "", "", fmt.Errorf("failed create token pair: %w", err)
	}
	return
}

func (a *Auth) RefreshToken(oldRefresh string) (newToken, newRefresh string, err error) {
	log := a.log.With(slog.String("auth.method", "refresh_token"))

	oldRefreshHash, err := bcrypt.GenerateFromPassword([]byte(oldRefresh), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate old refresh hash", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed refresh token: %w", err)
	}

	user, err := a.userRepo.GetUserByRefresh(context.Background(), oldRefreshHash)
	if err != nil {
		log.Warn(ErrGetUserRefresh.Error(), slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed refresh token: %w", ErrGetUserRefresh)
	}
	newToken, err = a.jwtManager.CreateJwt(user, a.tokenTTL)
	if err != nil {
		log.Error(ErrCreateJWT.Error(), slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed refresh token: %w", ErrCreateJWT)
	}
	newRefresh, err = a.jwtManager.CreateRefresh()
	if err != nil {
		log.Error(ErrCreateRefresh.Error(), slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed refresh token: %w", ErrCreateRefresh)
	}
	newRefreshHash, err := bcrypt.GenerateFromPassword([]byte(newRefresh), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate new refresh hash", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed refresh token: %w", err)
	}
	err = a.userRepo.SaveRefreshToken(context.Background(), user.UUID, newRefreshHash, a.refreshTTL)
	if err != nil {
		log.Error("failed to save refresh token")
		return "", "", fmt.Errorf("failed refresh token: %w", err)
	}
	return
}