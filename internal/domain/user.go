package domain

type User struct {
	ID       int
	Username string
	Password string
}

type UserService interface {
	Login(username, password string) (string, error)
	ValidateToken(token string) (User, error)
}
