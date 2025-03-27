package usecase

import (
	"auth-comparison/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo     domain.AuthRepository
	secret   string
	expireIn time.Duration
}

func NewAuthService(repo domain.AuthRepository, secret string, expireIn time.Duration) domain.AuthService {
	return &authService{
		repo:     repo,
		secret:   secret,
		expireIn: expireIn,
	}
}

func (s *authService) Login(username, password string) (domain.TokenPair, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return domain.TokenPair{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return domain.TokenPair{}, errors.New("invalid credentials")
	}

	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refreshToken := uuid.NewString()
	exp := time.Now().Add(7 * 24 * time.Hour).Unix()

	err = s.repo.UpsertSession(user.ID, refreshToken, exp)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) ValidateToken(tokenStr string) (domain.User, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return domain.User{}, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.User{}, errors.New("invalid token claims")
	}
	return domain.User{
		ID:       claims["user_id"].(int),
		Username: claims["username"].(string),
	}, nil
}

func (s *authService) GenerateAccessToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(s.expireIn).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *authService) RefreshToken(refreshToken string) (string, error) {
	session, err := s.repo.GetSessionByRefreshToken(refreshToken)
	if err != nil || session.ExpiresAt < time.Now().Unix() {
		return "", errors.New("invalid or expired refresh token")
	}

	user, err := s.repo.GetUserByID(session.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	return s.GenerateAccessToken(user)
}

func (s *authService) Logout(userID int) error {
	return s.repo.DeleteSessionByUserID(userID)
}

func (s *authService) Register(username, password string) error {
	_, err := s.repo.GetUserByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return s.repo.CreateUser(newUser)
}

func (s *authService) ValidateUsername(username string) error {
	_, err := s.repo.GetUserByUsername(username)
	return err
}
