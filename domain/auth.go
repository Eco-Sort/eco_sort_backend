package domain

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenPayload struct {
	UserId uint  `json:"user_id"`
	Role   Role  `json:"role"`
	Exp    int64 `json:"exp"`
}
type AuthService interface {
	AuthenticateAdmin(user *AuthLoginRequest) (User, error)
	ValidateToken(token string) (TokenPayload, error)
	Register(req AuthRegisterRequest) (uint, error)
}
