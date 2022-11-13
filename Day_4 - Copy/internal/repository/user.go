package repository

import (
	"api-mvc/internal/model"
	pkgdto "api-mvc/pkg/dto"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type User interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Get(ctx context.Context, user *model.User) (*model.User, error)
	FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, p *pkgdto.Pagination) ([]*model.User, *pkgdto.PaginationInfo, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, user *model.User) (*model.User, error)
}

type user struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewUserRepository(DB *gorm.DB, Redis *redis.Client) User {
	return &user{DB, Redis}
}

func (r *user) Create(ctx context.Context, user *model.User) (*model.User, error) {
	str, err := r.Redis.Get(fmt.Sprintf("user_email:%v", user.Email)).Result()
	if err == nil {
		json.Unmarshal([]byte(str), &user)
		return user, nil
	}

	err = r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	user, err = r.Get(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *user) Get(ctx context.Context, user *model.User) (*model.User, error) {
	str, err := r.Redis.Get(fmt.Sprintf("user_email:%v", user.Email)).Result()
	if err == nil {
		json.Unmarshal([]byte(str), &user)
		return user, nil
	}

	err = r.DB.Where(&user).First(&user).Error
	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(user)
	_, err = r.Redis.Set(fmt.Sprintf("user_email:%v", user.Email), b, time.Duration(1*time.Hour)).Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *user) FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]*model.User, *pkgdto.PaginationInfo, error) {
	users := make([]*model.User, 0)
	var count int64

	query := r.DB.Model(&model.User{})

	if payload.Search != "" {
		search := fmt.Sprintf("%%%s%%", strings.ToLower(payload.Search))
		query = query.Where("lower(name) LIKE ? or lower(email) LIKE ? ", search, search)
	}

	countquery := query
	if err := countquery.Count(&count).Error; err != nil {
		return nil, nil, err
	}

	limit, offset := pkgdto.GetLimitOffset(pagination)

	err := query.Limit(limit).Offset(offset).Find(&users).Error

	return users, pkgdto.CheckInfoPagination(pagination, count), err
}

func (r *user) Update(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(user).Error
	if err != nil {
		return nil, err
	}

	_, err = r.Redis.Del(fmt.Sprintf("user_email:%v", user.Email)).Result()
	if err != nil {
		return nil, err
	}

	user, err = r.Get(ctx, &model.User{Base: model.Base{ID: user.ID}})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *user) Delete(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.DB.Delete(user).Error
	if err != nil {
		return nil, err
	}

	_, err = r.Redis.Del(fmt.Sprintf("user_email:%v", user.Email)).Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}
