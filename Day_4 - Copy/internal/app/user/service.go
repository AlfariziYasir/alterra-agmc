package user

import (
	"api-mvc/internal/dto"
	"api-mvc/internal/factory"
	"api-mvc/internal/model"
	"api-mvc/internal/pkg/token"
	"api-mvc/internal/repository"
	pkgdto "api-mvc/pkg/dto"
	"api-mvc/pkg/util"
	"api-mvc/pkg/util/response"
	"context"
	"errors"

	"gorm.io/gorm"
)

type service struct {
	AuthRepo repository.Auth
	UserRepo repository.User
}

type Service interface {
	Create(ctx context.Context, payload *dto.UserRequest) (*dto.UserJwtResponse, error)
	Get(ctx context.Context, payload *dto.UserRequest) (*dto.UserResponse, error)
	Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.UserResponse], error)
	Update(ctx context.Context, payload *dto.UserRequest) (*dto.UserResponse, error)
	UpdatePassword(ctx context.Context, payload *dto.UserRequest) (*dto.UserResponseDetail, error)
	Delete(ctx context.Context, payload *dto.UserRequest) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepo: f.User,
		AuthRepo: f.Auth,
	}
}

func (s *service) Create(ctx context.Context, payload *dto.UserRequest) (*dto.UserJwtResponse, error) {
	isExist, err := s.UserRepo.Get(ctx, &model.User{Email: payload.Email})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if isExist != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("user is already exist"))
	}

	hashPassword, err := util.HashPassword(payload.Password)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	user := &model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashPassword,
		Base: model.Base{
			CreatedBy: ctx.Value("username").(string),
			UpdatedBy: ctx.Value("username").(string),
		},
	}
	user, err = s.UserRepo.Create(ctx, user)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	claims := token.CreateJWTClaims(map[string]interface{}{
		"user_Email": user.ID,
		"username":   user.Name,
		"email":      user.Email,
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

func (s *service) Get(ctx context.Context, payload *dto.UserRequest) (*dto.UserResponse, error) {
	user, err := s.UserRepo.Get(ctx, &model.User{
		Base: model.Base{
			ID: payload.ID,
		},
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *service) Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.UserResponse], error) {
	users, info, err := s.UserRepo.FindAll(ctx, payload, &payload.Pagination)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	data := make([]dto.UserResponse, 0)

	for _, user := range users {
		data = append(data, dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return &pkgdto.SearchGetResponse[dto.UserResponse]{
		Data:           data,
		PaginationInfo: *info,
	}, nil
}

func (s *service) Update(ctx context.Context, payload *dto.UserRequest) (*dto.UserResponse, error) {
	user, err := s.UserRepo.Get(ctx, &model.User{Email: payload.Email})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if user.ID != payload.ID {
		return nil, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("email is already exist"))
	}

	user.Name = payload.Name
	user.Email = payload.Email
	user.UpdatedBy = ctx.Value("username").(string)
	_, err = s.UserRepo.Update(ctx, user)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *service) UpdatePassword(ctx context.Context, payload *dto.UserRequest) (*dto.UserResponseDetail, error) {
	user, err := s.UserRepo.Get(ctx, &model.User{Base: model.Base{ID: payload.ID}})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	hashPassword, err := util.HashPassword(payload.Password)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	user.Password = hashPassword
	user.UpdatedBy = ctx.Value("username").(string)
	_, err = s.UserRepo.Update(ctx, user)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	return &dto.UserResponseDetail{
		UserResponse: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		Password: hashPassword,
	}, nil
}

func (s *service) Delete(ctx context.Context, payload *dto.UserRequest) error {
	user, err := s.UserRepo.Get(ctx, &model.User{Base: model.Base{ID: payload.ID}})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	user.DeletedBy = ctx.Value("username").(string)
	_, err = s.UserRepo.Delete(ctx, user)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	return nil
}
