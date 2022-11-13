package auth

import (
	"api-mvc/internal/dto"
	"api-mvc/internal/factory"
	"api-mvc/internal/model"
	"api-mvc/internal/pkg/token"
	"api-mvc/internal/repository"
	"api-mvc/pkg/util"
	"api-mvc/pkg/util/response"
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type service struct {
	AuthRepo repository.Auth
	UserRepo repository.User
}

type Service interface {
	Login(ctx context.Context, payload *dto.LoginRequest) (*dto.UserJwtResponse, error)
	Refresh(ctx context.Context, payload *dto.AccessDetails) (*dto.UserJwtResponse, error)
	Logout(ctx context.Context, payload *dto.AccessDetails) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		AuthRepo: f.Auth,
		UserRepo: f.User,
	}
}

func (s *service) Login(ctx context.Context, payload *dto.LoginRequest) (*dto.UserJwtResponse, error) {
	user, err := s.UserRepo.Get(ctx, &model.User{Email: payload.Email})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}

		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	key := fmt.Sprintf("login:%v%v", user.ID, user.Name)
	uuid, err := s.AuthRepo.FetchUuid(key)
	if err == nil {
		metaData := &dto.AccessDetails{
			TokenUuid: uuid["token_uuid"].(string),
			Username:  user.Name,
			UserId:    user.ID,
		}

		err = s.AuthRepo.DeleteTokens(metaData)
		if err != nil {
			return nil, response.ErrorBuilder(
				&response.ErrorConstant.InternalServerError,
				errors.New("generating token error"),
			)
		}
	}

	if !(util.CompareHashPassword(payload.Password, user.Password)) {
		return nil, response.ErrorBuilder(
			&response.ErrorConstant.EmailOrPasswordIncorrect,
			errors.New(response.ErrorConstant.EmailOrPasswordIncorrect.Response.Meta.Message),
		)
	}

	claims := token.CreateJWTClaims(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Name,
		"email":    user.Email,
	})
	t, err := token.CreateToken(claims)
	if err != nil {
		return nil, response.ErrorBuilder(
			&response.ErrorConstant.InternalServerError,
			errors.New("generating token error"),
		)
	}

	err = s.AuthRepo.CreateAuth(claims, t)
	if err != nil {
		return nil, response.ErrorBuilder(
			&response.ErrorConstant.InternalServerError,
			errors.New("error create auth"),
		)
	}

	return &dto.UserJwtResponse{
		UserResponse: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		JWTAccess:  t.AccessToken,
		JWTRefresh: t.RefreshToken,
	}, nil
}

func (s *service) Refresh(ctx context.Context, payload *dto.AccessDetails) (*dto.UserJwtResponse, error) {
	_, err := s.AuthRepo.FetchAuth(payload.TokenUuid)
	if err != nil {
		return nil, response.ErrorBuilder(
			&response.ErrorConstant.InternalServerError,
			errors.New("fetching token error"),
		)
	}

	err = s.AuthRepo.DeleteTokens(payload)
	if err != nil {
		return nil, response.ErrorBuilder(
			&response.ErrorConstant.InternalServerError,
			errors.New("deleting token error"),
		)
	}

	claims := token.CreateJWTClaims(map[string]interface{}{
		"user_id":  payload.UserId,
		"username": payload.Username,
		"email":    payload.Email,
	})
	t, err := token.CreateToken(claims)
	if err != nil {
		return nil, response.ErrorBuilder(
			&response.ErrorConstant.InternalServerError,
			errors.New("generating token error"),
		)
	}

	err = s.AuthRepo.CreateAuth(claims, t)
	if err != nil {
		return nil, response.ErrorBuilder(
			&response.ErrorConstant.InternalServerError,
			errors.New("error create auth"),
		)
	}

	return &dto.UserJwtResponse{
		UserResponse: dto.UserResponse{
			ID:    payload.UserId,
			Name:  payload.Username,
			Email: payload.Email,
		},
		JWTAccess:  t.AccessToken,
		JWTRefresh: t.RefreshToken,
	}, nil
}

func (s *service) Logout(ctx context.Context, payload *dto.AccessDetails) error {
	err := s.AuthRepo.DeleteTokens(payload)
	if err != nil {
		return response.ErrorBuilder(
			&response.ErrorConstant.InternalServerError,
			errors.New("deleting token error"),
		)
	}

	return nil
}
