package domain

type AuthSession struct {
	UserID       string
	RefreshToken string
	ExpiresAt    int64
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthService interface {
	Login(username, password string) (TokenPair, error)
	ValidateToken(token string) (User, error)
	RefreshToken(refreshToken string) (string, error)
	Logout(userID string) error
	GenerateAccessToken(user User) (string, error)
	Register(username, password string) error
	ValidateUsername(username string) error
}

type AuthRepository interface {
	GetUserByUsername(username string) (User, error)
	GetUserByID(userID string) (User, error)
	UpsertSession(userID string, refreshToken string, expiresAt int64) error
	GetSessionByRefreshToken(token string) (AuthSession, error)
	DeleteSessionByUserID(userID string) error
	CreateUser(user User) error
}
