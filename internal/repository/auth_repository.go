package repository

import (
	"auth-comparison/internal/domain"
	"database/sql"
	"time"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetUserByID(userID string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, username, password FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *authRepository) DeleteSessionByUserID(userID int) error {
	_, err := r.db.Exec("DELETE FROM auth_sessions WHERE user_id = $1", userID)
	return err
}

func (r *authRepository) GetSessionByRefreshToken(token string) (domain.AuthSession, error) {
	var session domain.AuthSession
	err := r.db.QueryRow("SELECT user_id, refresh_token, expires_at FROM auth_sessions WHERE refresh_token = $1", token).
		Scan(&session.UserID, &session.RefreshToken, &session.ExpiresAt)
	if err != nil {
		return domain.AuthSession{}, err
	}
	return session, nil
}

func (r *authRepository) GetUserByUsername(username string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *authRepository) UpsertSession(userID int, refreshToken string, expiresAt int64) error {
	_, err := r.db.Exec(`
		INSERT INTO auth_sessions (user_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET refresh_token = EXCLUDED.refresh_token, expires_at = EXCLUDED.expires_at
	`, userID, refreshToken, time.Unix(expiresAt, 0))
	return err
}

func (r *authRepository) CreateUser(user domain.User) error {
	_, err := r.db.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, user.Username, user.Password)
	return err
}
