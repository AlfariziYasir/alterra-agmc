package dto

import "github.com/golang-jwt/jwt/v4"

type (
	LoginRequest struct {
		Email    string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	JWTClaims struct {
		TokenUuid string
		UserID    uint   `json:"user_id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		jwt.RegisteredClaims
	}

	TokenDetails struct {
		AccessToken  string
		RefreshToken string
		TokenUuid    string
		RefreshUuid  string
		AtExpires    int64
		RtExpires    int64
	}

	AccessDetails struct {
		TokenUuid string
		UserId    uint
		Username  string
		Email     string
	}
)
