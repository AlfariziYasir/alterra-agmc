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

func (c *Controller) CreateBook(ctx echo.Context) error {
	book := model.Book{}
	req := model.BookRequest{}

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

	book.Title = req.Title
	book.Isbn = req.Isbn
	b, err := book.Get(c.DB)
	if err != nil && err != gorm.ErrRecordNotFound {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	} else if b != nil {
		return ctx.JSON(http.StatusConflict, web.ResponError{
			Status:  http.StatusConflict,
			Message: errors.New("book is already created").Error(),
		})
	}

	book.Title = req.Title
	book.Isbn = req.Isbn
	book.Writer = req.Writer
	book.UpdatedAt = time.Now()
	err = book.Create(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	res := model.BookResponse{
		ID:     book.ID,
		Title:  book.Title,
		Isbn:   book.Isbn,
		Writer: book.Writer,
	}

	return ctx.JSON(http.StatusCreated, web.Respons{
		Status:  http.StatusCreated,
		Message: "create book is success",
		Data:    res,
	})
}

func (c *Controller) GetBook(ctx echo.Context) error {
	book := model.Book{}

	id, _ := strconv.Atoi(ctx.Param("id"))

	book.ID = uint(id)
	b, err := book.Get(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	res := model.BookResponse{
		ID:     b.ID,
		Title:  b.Title,
		Isbn:   b.Isbn,
		Writer: b.Writer,
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "get book is success",
		Data:    res,
	})
}

func (c *Controller) GetBooks(ctx echo.Context) error {
	book := model.Book{}
	res := make([]model.BookResponse, 0)

	users, err := book.Gets(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	for _, book := range users {
		res = append(res, model.BookResponse{
			ID:     book.ID,
			Title:  book.Title,
			Isbn:   book.Isbn,
			Writer: book.Writer,
		})
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "get books is success",
		Data:    res,
	})
}

func (c *Controller) UpdateBook(ctx echo.Context) error {
	book := model.Book{}
	req := model.BookRequest{}

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

	book.ID = uint(id)
	b, _ := book.Get(c.DB)

	b.Title = req.Title
	b.Isbn = req.Isbn
	b.Writer = req.Writer
	b.UpdatedAt = time.Now()
	err = b.Update(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	res := model.BookResponse{
		ID:     book.ID,
		Title:  book.Title,
		Isbn:   book.Isbn,
		Writer: book.Writer,
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "update book is success",
		Data:    res,
	})
}

func (c *Controller) DeleteBook(ctx echo.Context) error {
	book := model.Book{}

	id, _ := strconv.Atoi(ctx.Param("id"))

	book.ID = uint(id)
	b, err := book.Get(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	err = b.Delete(c.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, web.ResponError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, web.Respons{
		Status:  http.StatusOK,
		Message: "delete book is success",
		Data:    nil,
	})
}
