package domain

type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenPayload struct {
	UserId uint  `json:"user_id"`
	Exp    int64 `json:"exp"`
}
type AuthService interface {
	AuthenticateAdmin(user *AuthLoginRequest) (bool, error)
	ValidateToken(token string) (TokenPayload, error)
}
