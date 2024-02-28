package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/EwvwGeN/medods_assignment/internal/domain/models"
	"github.com/EwvwGeN/medods_assignment/internal/storage"
	"github.com/golang-jwt/jwt"
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
	SaveRefreshToken(ctx context.Context, uuid, refresh string, refreshTTL time.Duration) (err error)
	GetUserByUUID(ctx context.Context, uuid string) (user *models.User, err error)
}

type JwtManager interface {
	CreateJwt(user *models.User, ttl time.Duration) (token string, err error)
	CreateRefresh() (refresh string, err error)
	ParseTokenClaims(token string) (jwt.MapClaims, error) 
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
		if errors.Is(err, storage.ErrUserExist) {
			log.Info("failed register user", slog.String("error", storage.ErrUserExist.Error()))
			return "", fmt.Errorf("failed register user: %w", storage.ErrUserExist)
		}
		log.Warn("failed register user", slog.String("uuid", uuid), slog.String("error", err.Error()))
		return "", fmt.Errorf("failed register user: %w", ErrSaveUser)
	}
	return
}

func (a *Auth) CreateTokenPair(uuid string) (token, refresh string, err error) {
	log := a.log.With(slog.String("auth.method", "create_token_pair"))
	user, err:= a.userRepo.GetUserByUUID(context.Background(), uuid)
	log.Debug("got user", slog.Any("user", user))
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
	log.Debug("creted refresh hash", slog.String("hash", string(refreshHash)))
	err = a.userRepo.SaveRefreshToken(context.Background(), user.UUID, string(refreshHash), a.refreshTTL)
	if err != nil {
		log.Error("failed to save refresh token")
		return "", "", fmt.Errorf("failed create token pair: %w", err)
	}
	refresh = base64.StdEncoding.EncodeToString([]byte(refresh))
	return
}

func (a *Auth) RefreshToken(accessToken, refreshToken string) (newToken, newRefresh string, err error) {
	log := a.log.With(slog.String("auth.method", "refresh_token"))
	log.Debug("start refreshing", slog.String("access_token", accessToken), slog.String("refresh_token", refreshToken))
	tokenClaims, err := a.jwtManager.ParseTokenClaims(accessToken)
	if err != nil {
		log.Info("not valid access token", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed refresh token: %w", err)
	}
	uuid, ok := tokenClaims["uuid"]
	if !ok {
		log.Info("cant get uuid from token claims")
		return "", "", fmt.Errorf("failed refresh token: %w", ErrValidAccess)
	}
	queryTime := time.Now().Unix()
	user, err := a.userRepo.GetUserByUUID(context.Background(), uuid.(string))
	log.Debug("gor user by uuid", slog.Any("user", user))
	if err != nil {
		log.Warn(ErrGetUserUUID.Error(), slog.String("error", err.Error()))
		return "", "", fmt.Errorf("failed refresh token: %w", ErrGetUserUUID)
	}
	refreshBytes, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return "", "", ErrValidRefresh
	}
	refreshToken = string(refreshBytes)
	err = bcrypt.CompareHashAndPassword([]byte(user.RefreshHash), []byte(refreshToken))
	if err != nil {
		return "", "", fmt.Errorf("failed refresh token: %w", ErrValidRefresh)
	}
	if queryTime > user.ExpiresAt {
		return "", "", fmt.Errorf("failed refresh token: %w", ErrValidRefresh)
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
	err = a.userRepo.SaveRefreshToken(context.Background(), user.UUID, string(newRefreshHash), a.refreshTTL)
	if err != nil {
		log.Error("failed to save refresh token")
		return "", "", fmt.Errorf("failed refresh token: %w", err)
	}
	newRefresh = base64.StdEncoding.EncodeToString([]byte(newRefresh))
	return
}