package auth

import (
	"api-mvc/internal/dto"
	"api-mvc/internal/factory"
	"api-mvc/internal/pkg/token"
	"api-mvc/pkg/util"
	"api-mvc/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) Login(c echo.Context) error {
	payload := new(dto.LoginRequest)

	err := c.Bind(&payload)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	err = util.Struct(payload)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	res, err := h.service.Login(c.Request().Context(), payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(res).Send(c)
}

func (h *handler) Refresh(c echo.Context) error {
	t, err := token.ExtractTokenMetadata(c.Request())
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, err).Send(c)
	}

	res, err := h.service.Refresh(c.Request().Context(), t)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err).Send(c)
	}

	return response.SuccessResponse(res).Send(c)
}

func (h *handler) Logout(c echo.Context) error {
	t, err := token.ExtractTokenMetadata(c.Request())
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, err).Send(c)
	}

	err = h.service.Logout(c.Request().Context(), t)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err).Send(c)
	}

	return response.CustomSuccessBuilder(response.SuccessConstant.OK.Code, nil, "logout success", nil).Send(c)
}
