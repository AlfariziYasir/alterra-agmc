package controller

import (
	"api-mvc/lib"
	"api-mvc/model"
	"api-mvc/web"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (c *Controller) CreateUser(ctx echo.Context) error {
	user := model.User{}
	req := model.UserRequest{}

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, web.ResponError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	err = lib.Struct(req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, web.ResponError{
			Status:  http.StatusConflict,
			Message: err.Error(),
		})
	}

	user.Name = req.Name
	u, err := user.Get(c.DB)
	if err != nil && err != gorm.ErrRecordNotFound {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	} else if u != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: errors.New("user is already created").Error(),
		})
	}

	user.Name = req.Name
	user.Password = req.Password
	user.UpdatedAt = time.Now()
	err = user.Create(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	res := model.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Password: user.Password,
	}

	return ctx.JSON(http.StatusCreated, web.Respons{
		Status:  http.StatusCreated,
		Message: "create user is success",
		Data:    res,
	})
}

func (c *Controller) GetUser(ctx echo.Context) error {
	user := model.User{}

	id, _ := strconv.Atoi(ctx.Param("id"))

	user.ID = uint(id)
	u, err := user.Get(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	res := model.UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Password: u.Password,
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "get user is success",
		Data:    res,
	})
}

func (c *Controller) GetUsers(ctx echo.Context) error {
	user := model.User{}
	res := make([]model.UserResponse, 0)

	users, err := user.Gets(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	for _, user := range users {
		res = append(res, model.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Password: user.Password,
		})
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "get users is success",
		Data:    res,
	})
}

func (c *Controller) UpdateUser(ctx echo.Context) error {
	user := model.User{}
	req := model.UserRequest{}

	id, _ := strconv.Atoi(ctx.Param("id"))

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, web.ResponError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	err = lib.Struct(req)
	if err != nil {
		return ctx.JSON(http.StatusConflict, web.ResponError{
			Status:  http.StatusConflict,
			Message: err.Error(),
		})
	}

	user.ID = uint(id)
	u, _ := user.Get(c.DB)

	u.Name = req.Name
	u.Password = req.Password
	u.UpdatedAt = time.Now()
	err = u.Update(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	res := model.UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Password: u.Password,
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "update user is success",
		Data:    res,
	})
}

func (c *Controller) DeleteUser(ctx echo.Context) error {
	user := model.User{}

	id, _ := strconv.Atoi(ctx.Param("id"))

	user.ID = uint(id)
	u, err := user.Get(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	err = u.Delete(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "delete user is success",
		Data:    nil,
	})
}
