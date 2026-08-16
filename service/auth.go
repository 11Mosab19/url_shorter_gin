package service

import (
	"os"
	"time"
	"url_shorter_gin/apperrors"
	"url_shorter_gin/dto"
	"url_shorter_gin/models"

	"github.com/golang-jwt/jwt/v5"
)

func GetConfig() string {
	SecretKey := os.Getenv("SECRET_KEY")
	return SecretKey
}

type AuthService struct {
	US UserService
}

func (AuS *AuthService) GenerateToken(Data models.User) (dto.Token, error) {
	claims := models.Claims{
		Role:   Data.Role,
		UserId: Data.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	obj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	MySecretKey := []byte(GetConfig())
	key, err := obj.SignedString(MySecretKey)
	if err != nil {
		return dto.Token{}, err
	}
	return dto.Token{Key: key}, nil
}

func (AS *AuthService) VerifyToken(TokenString string) (*models.Claims, error) {
	MySecretKey := []byte(GetConfig())
	VerifiedClaims := &models.Claims{}
	Returned, err := jwt.ParseWithClaims(TokenString, VerifiedClaims, func(t *jwt.Token) (any, error) { return MySecretKey, nil })
	if err != nil {
		return &models.Claims{}, apperrors.ErrUnauthorized
	}
	if !Returned.Valid {
		return &models.Claims{}, apperrors.ErrUnauthorized
	}
	return VerifiedClaims, nil
}
