package usecase

import (
	"auth-comparison/internal/domain"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type sessionService struct {
	store *session.Store
}

func NewSessionService(store *session.Store) *sessionService {
	return &sessionService{store: store}
}

func (s *sessionService) Login(c *fiber.Ctx, username, password string) error {
	if username != "admin" || password != "password" {
		return errors.New("invalid credentials")
	}
	sess, err := s.store.Get(c)
	if err != nil {
		return err
	}
	sess.Set("username", username)
	return sess.Save()
}

func (s *sessionService) Validate(c *fiber.Ctx) (domain.User, error) {
	sess, err := s.store.Get(c)
	if err != nil {
		return domain.User{}, err
	}

	username := sess.Get("username")
	if usernameStr, ok := username.(string); ok {
		return domain.User{Username: usernameStr}, nil
	}
	return domain.User{}, errors.New("no session")
}
