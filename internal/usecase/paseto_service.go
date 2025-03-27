package usecase

import (
	"auth-comparison/internal/domain"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type pasetoService struct {
	paseto     *paseto.V2
	secretKey  []byte
	issuer     string
	expiration time.Duration
}

func NewPasetoService(secretKey []byte) domain.UserService {
	return &pasetoService{
		paseto:     paseto.NewV2(),
		secretKey:  secretKey,
		issuer:     "auth-comparison",
		expiration: time.Hour,
	}
}

type customPayload struct {
	Username string    `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	Exp      time.Time `json:"exp"`
	ID       string    `json:"id"`
}

func (s *pasetoService) Login(username, password string) (string, error) {
	if username != "admin" || password != "password" {
		return "", errors.New("invalid credentials")
	}

	payload := customPayload{
		Username: username,
		IssuedAt: time.Now(),
		Exp:      time.Now().Add(s.expiration),
		ID:       uuid.NewString(),
	}

	token, err := s.paseto.Encrypt(s.secretKey, payload, nil)
	return token, err
}

func (s *pasetoService) ValidateToken(token string) (domain.User, error) {
	var payload customPayload
	err := s.paseto.Decrypt(token, s.secretKey, &payload, nil)
	if err != nil {
		return domain.User{}, err
	}

	if time.Now().After(payload.Exp) {
		return domain.User{}, errors.New("token expired")
	}

	return domain.User{
		Username: payload.Username,
	}, nil
}
