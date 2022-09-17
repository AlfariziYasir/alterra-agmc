package controller

import (
	"api-mvc/lib"
	"api-mvc/lib/token"
	"api-mvc/model"
	"api-mvc/web"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (c *Controller) Login(ctx echo.Context) error {
	user := model.User{}
	auth := model.Auth{}
	tk := token.NewToken()

	err := ctx.Bind(&auth)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, web.ResponError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	err = lib.Struct(auth)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, web.ResponError{
			Status:  http.StatusConflict,
			Message: err.Error(),
		})
	}

	user.Email = auth.Email
	u, err := user.Get(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, web.ResponError{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
	}

	if u.TokenUuid != " " {
		refreshUuid := fmt.Sprintf("%s++%v%s", u.TokenUuid, u.ID, u.Name)
		//delete access token
		_, err = c.Redis.Del(context.Background(), u.TokenUuid).Result()
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, web.ResponError{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
		}

		_, err = c.Redis.Del(context.Background(), refreshUuid).Result()
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, web.ResponError{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(auth.Password))
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, web.ResponError{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	claims := map[string]interface{}{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	}

	ts, err := tk.CreateToken(claims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	at := time.Unix(ts.AtExpires, 0)
	rt := time.Unix(ts.RtExpires, 0)
	now := time.Now()
	b, _ := json.Marshal(claims)
	atCreated, err := c.Redis.Set(context.Background(), ts.TokenUuid, b, at.Sub(now)).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	rtCreated, err := c.Redis.Set(context.Background(), ts.RefreshUuid, b, rt.Sub(now)).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if atCreated == "0" || rtCreated == "0" {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: errors.New("no record inserted").Error(),
		})
	}

	u.TokenUuid = ts.TokenUuid
	err = u.Update(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: errors.New("no record inserted").Error(),
		})
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "login success",
		Data: model.AuthResponse{
			AccessToken:  ts.AccessToken,
			RefreshToken: ts.RefreshToken,
		},
	})
}

func (c *Controller) Refresh(ctx echo.Context) error {
	tk := token.NewToken()
	user := model.User{}

	t, err := tk.ExtractTokenMetadata(ctx.Request())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, web.ResponError{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	_, err = c.Redis.Get(context.Background(), t.TokenUuid).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	refreshUuid := fmt.Sprintf("%s++%v%s", t.TokenUuid, t.UserId, t.Username)
	//delete access token
	_, err = c.Redis.Del(context.Background(), t.TokenUuid).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	_, err = c.Redis.Del(context.Background(), refreshUuid).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	claims := map[string]interface{}{
		"user_id": t.UserId,
		"name":    t.Username,
		"email":   t.Email,
	}

	ts, err := tk.CreateToken(claims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	at := time.Unix(ts.AtExpires, 0)
	rt := time.Unix(ts.RtExpires, 0)
	now := time.Now()
	b, _ := json.Marshal(claims)
	atCreated, err := c.Redis.Set(context.Background(), ts.TokenUuid, b, at.Sub(now)).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	rtCreated, err := c.Redis.Set(context.Background(), ts.RefreshUuid, b, rt.Sub(now)).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if atCreated == "0" || rtCreated == "0" {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: errors.New("no record inserted").Error(),
		})
	}

	user.ID = t.UserId
	u, err := user.Get(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, web.ResponError{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
	}

	u.TokenUuid = ts.TokenUuid
	err = u.Update(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: errors.New("no record inserted").Error(),
		})
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "refresh success",
		Data: model.AuthResponse{
			AccessToken:  ts.AccessToken,
			RefreshToken: ts.RefreshToken,
		},
	})
}

func (c *Controller) Logout(ctx echo.Context) error {
	tk := token.NewToken()
	user := model.User{}

	t, err := tk.ExtractTokenMetadata(ctx.Request())
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, web.ResponError{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	refreshUuid := fmt.Sprintf("%s++%v%s", t.TokenUuid, t.UserId, t.Username)
	//delete access token
	_, err = c.Redis.Del(context.Background(), t.TokenUuid).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	_, err = c.Redis.Del(context.Background(), refreshUuid).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	user.ID = t.UserId
	u, err := user.Get(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, web.ResponError{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
	}

	u.TokenUuid = " "
	err = u.Update(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: errors.New("no record inserted").Error(),
		})
	}

	ctx.Response().Writer.Header().Del("username")
	ctx.Response().Writer.Header().Del("user_id")

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "logout success",
		Data:    nil,
	})
}
