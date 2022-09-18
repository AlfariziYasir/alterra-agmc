package web

import "github.com/labstack/echo/v4"

type Respons struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func MarshalPayload(c echo.Context, code int, message string, payload interface{}) {
	res := Respons{
		Status:  code,
		Message: message,
		Data:    payload,
	}

	c.JSON(code, res)
}

func MarshalError(c echo.Context, code int, message string) {
	res := ResponError{
		Status:  code,
		Message: message,
	}

	c.JSON(code, res)
}
