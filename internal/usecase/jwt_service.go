package usecase

import (
	"auth-comparison/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey string
}

func NewJWTService(secret string) domain.AuthService {
	return &jwtService{secretKey: secret}
}

func (s *jwtService) Login(username, password string) (string, error) {
	if username != "admin" || password != "password" {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenStr string) (domain.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return domain.User{
			Username: claims["username"].(string),
		}, nil
	} else {
		return domain.User{}, err
	}
}
