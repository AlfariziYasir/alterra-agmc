package dto

type (
	UserRequest struct {
		ID         uint   `param:"id"`
		Name       string `json:"name" validate:"required,alphanum,min=4,max=10"`
		Email      string `json:"email" validate:"required;email"`
		Password   string `json:"password" validate:"required,min=8"`
		RePassword string `json:"repassword" validate:"required,max=20,min=8,eqfield=Password"`
	}

	UserResponse struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	UserResponseDetail struct {
		UserResponse
		Password string `json:"password"`
	}

	UserJwtResponse struct {
		UserResponse
		JWTAccess  string `json:"access_token"`
		JWTRefresh string `json:"refresh_token"`
	}
)
