package domain

type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ReEnterPassword string `json:"re_enter_password"`
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
