package model

type Auth struct {
	Email    string `json:"email" validate:"required,min=10,email"`
	Password string `json:"password" validate:"required,min=6,contains"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
