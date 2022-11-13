package user

import (
	"api-mvc/internal/dto"
	"api-mvc/internal/factory"
	pkgdto "api-mvc/pkg/dto"
	"api-mvc/pkg/util"
	"api-mvc/pkg/util/response"
	"errors"
	"strconv"

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

func (h *handler) Create(c echo.Context) error {
	payload := new(dto.UserRequest)

	err := c.Bind(&payload)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	err = util.Struct(payload)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	res, err := h.service.Create(c.Request().Context(), payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(res).Send(c)
}

func (h *handler) Get(c echo.Context) error {
	if c.Param("id") == "" {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, errors.New("parameter id not found")).Send(c)
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 0)

	res, err := h.service.Get(c.Request().Context(), &dto.UserRequest{ID: uint(id)})
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(res).Send(c)
}

func (h *handler) Find(c echo.Context) error {
	payload := new(pkgdto.SearchGetRequest)

	err := c.Bind(&payload)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	err = util.Struct(payload)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	res, err := h.service.Find(c.Request().Context(), payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.CustomSuccessBuilder(response.SuccessConstant.OK.Code, res.Data, "get users success", &res.PaginationInfo).Send(c)
}

func (h *handler) Update(c echo.Context) error {
	payload := new(dto.UserRequest)

	err := c.Bind(&payload)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	err = util.Struct(c)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	res, err := h.service.Update(c.Request().Context(), payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(res).Send(c)
}

func (h *handler) Delete(c echo.Context) error {
	if c.Param("id") == "" {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, errors.New("parameter id not found")).Send(c)
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 0)

	err := h.service.Delete(c.Request().Context(), &dto.UserRequest{ID: uint(id)})
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse("success delete user").Send(c)
}
