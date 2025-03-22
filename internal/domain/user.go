package domain

type User struct {
	ID       string
	Username string
	Password string
}

type AuthService interface {
	Login(username, password string) (string, error)
	ValidateToken(token string) (User, error)
}
