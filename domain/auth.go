package domain

type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthService interface {
	AuthenticateAdmin(user *AuthLoginRequest) (bool, error)
}
