package token

import (
	"api-mvc/config"
	"api-mvc/internal/dto"

	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

func CreateJWTClaims(data map[string]interface{}) []dto.JWTClaims {
	tokenUuid, _ := uuid.NewV4()
	tokenUuids := tokenUuid.String()
	RefreshUuid := fmt.Sprintf("%s++%v%v", tokenUuids, data["user_id"], data["email"])

	return []dto.JWTClaims{
		{
			TokenUuid: tokenUuids,
			UserID:    data["user_id"].(uint),
			Username:  data["username"].(string),
			Email:     data["email"].(string),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			},
		},
		{
			TokenUuid: RefreshUuid,
			UserID:    data["user_id"].(uint),
			Username:  data["username"].(string),
			Email:     data["email"].(string),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			},
		},
	}
}

func CreateToken(data []dto.JWTClaims) (*dto.TokenDetails, error) {
	var err error
	//Creating Access Token
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, data[0])
	accessToken, err := at.SignedString([]byte(config.Cfg().JwtSecretKey))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, data[1])
	refreshToken, err := rt.SignedString([]byte(config.Cfg().JwtRefreshKey))
	if err != nil {
		return nil, err
	}
	return &dto.TokenDetails{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenUuid:    data[0].TokenUuid,
		RefreshUuid:  data[1].TokenUuid,
		AtExpires:    data[0].RegisteredClaims.ExpiresAt.Unix(),
		RtExpires:    data[1].RegisteredClaims.ExpiresAt.Unix(),
	}, nil
}

func TokenValid(r *http.Request) (*dto.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("token invalid")
	}

	data, err := extract(token)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Cfg().JwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//get the token from the request body
func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func extract(token *jwt.Token) (*dto.AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		userId, userOk := claims["user_id"].(float64)
		username, usernameOk := claims["username"].(string)

		if !ok && !userOk && !usernameOk {
			return nil, errors.New("unauthorized")
		} else {
			return &dto.AccessDetails{
				TokenUuid: accessUuid,
				UserId:    uint(userId),
				Username:  username,
			}, nil
		}
	}
	return nil, errors.New("something went wrong")
}

func ExtractTokenMetadata(r *http.Request) (*dto.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	acc, err := extract(token)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
