package jwt

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/EwvwGeN/medods_assignment/internal/domain/models"
	"github.com/golang-jwt/jwt"
)

type jwtManager struct {
	secretKey string
}

func NewJwtManager(secretKey string) *jwtManager {
	return &jwtManager{
		secretKey: secretKey,
	}
}

func (jm *jwtManager) CreateJwt(user *models.User, ttl time.Duration) (token string, err error) {
	if user.Email == "" {
		return "", ErrEmptyValue
	}
	if user.UUID == "" {
		return "", ErrEmptyValue
	}
	tokenObject := jwt.New(jwt.SigningMethodHS512)
	claims := tokenObject.Claims.(jwt.MapClaims)
	claims["uuid"] = user.UUID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	token, err = tokenObject.SignedString([]byte(jm.secretKey))
	if err != nil {
		return "", err
	}
	return
}

func (jm *jwtManager) CreateRefresh() (refresh string, err error) {
	buffer := make([]byte, 32)
	gen := rand.New(rand.NewSource(time.Now().Unix()))
	if _, err = gen.Read(buffer); err != nil {
		return "", ErrRefreshGenerate
	}
	return fmt.Sprintf("%x", buffer), nil
}