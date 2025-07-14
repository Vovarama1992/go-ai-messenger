package httpadapter

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=320"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=320"`
	Password string `json:"password" validate:"required,min=8"`
}

type TokenResponse struct {
	UserID int64  `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	Token  string `json:"access_token,omitempty"`
}
