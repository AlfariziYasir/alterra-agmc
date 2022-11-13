package repository

import (
	"api-mvc/internal/dto"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Auth interface {
	CreateAuth([]dto.JWTClaims, *dto.TokenDetails) error
	FetchAuth(tokenUuid string) (map[string]interface{}, error)
	FetchUuid(key string) (map[string]interface{}, error)
	DeleteRefresh(string) error
	DeleteTokens(*dto.AccessDetails) error
}

type auth struct {
	Redis *redis.Client
}

func NewAuthRepository(Redis *redis.Client) Auth {
	return &auth{Redis}
}

func (r *auth) CreateAuth(claims []dto.JWTClaims, td *dto.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	b, _ := json.Marshal(claims)
	atCreated, err := r.Redis.Set(td.TokenUuid, b, at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := r.Redis.Set(td.RefreshUuid, b, rt.Sub(now)).Result()
	if err != nil {
		return err
	}

	uuid := map[string]interface{}{
		"token_uuid": td.TokenUuid,
	}
	bu, _ := json.Marshal(uuid)
	uuidCreated, err := r.Redis.Set(fmt.Sprintf("login:%v%v", claims[0].UserID, claims[0].Username), bu, at.Sub(now)).Result()
	if err != nil {
		return err
	}

	if atCreated == "0" || rtCreated == "0" || uuidCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

func (r *auth) FetchUuid(key string) (map[string]interface{}, error) {
	authD, err := r.Redis.Get(key).Result()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}

	_ = json.Unmarshal([]byte(authD), &data)

	return data, nil
}

func (r *auth) FetchAuth(tokenUuid string) (map[string]interface{}, error) {
	authD, err := r.Redis.Get(tokenUuid).Result()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}

	_ = json.Unmarshal([]byte(authD), &data)

	return data, nil
}

func (r *auth) DeleteTokens(authD *dto.AccessDetails) error {
	refreshUuid := fmt.Sprintf("%s++%v%s", authD.TokenUuid, authD.UserId, authD.Username)
	//delete access token
	_, err := r.Redis.Del(authD.TokenUuid).Result()
	if err != nil {
		return err
	}

	_, err = r.Redis.Del(refreshUuid).Result()
	if err != nil {
		return err
	}

	_, err = r.Redis.Del(fmt.Sprintf("login:%v%v", authD.UserId, authD.Username)).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *auth) DeleteRefresh(refreshUuid string) error {
	deleted, err := r.Redis.Del(refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
