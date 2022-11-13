package book

import (
	"api-mvc/internal/dto"
	"api-mvc/internal/factory"
	"api-mvc/internal/model"
	"api-mvc/internal/repository"
	pkgdto "api-mvc/pkg/dto"
	"api-mvc/pkg/util/response"
	"context"
	"errors"

	"gorm.io/gorm"
)

type service struct {
	BookRepo repository.Book
}

type Service interface {
	Create(ctx context.Context, payload *dto.BookRequest) (*dto.BookResponse, error)
	Get(ctx context.Context, payload *dto.BookRequest) (*dto.BookResponse, error)
	Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.BookResponse], error)
	Update(ctx context.Context, payload *dto.BookRequest) (*dto.BookResponse, error)
	Delete(ctx context.Context, payload *dto.BookRequest) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		BookRepo: f.Book,
	}
}

func (s *service) Create(ctx context.Context, payload *dto.BookRequest) (*dto.BookResponse, error) {
	isExist, err := s.BookRepo.Get(ctx, &model.Book{Title: payload.Title, Isbn: payload.Isbn})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if isExist != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("book is already exist"))
	}

	book := &model.Book{
		Title:  payload.Title,
		Isbn:   payload.Isbn,
		Writer: payload.Writer,
		Base: model.Base{
			CreatedBy: ctx.Value("username").(string),
			UpdatedBy: ctx.Value("username").(string),
		},
	}
	book, err = s.BookRepo.Create(ctx, book)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	return &dto.BookResponse{
		ID:     book.ID,
		Title:  book.Title,
		Isbn:   book.Isbn,
		Writer: book.Writer,
	}, nil
}

func (s *service) Get(ctx context.Context, payload *dto.BookRequest) (*dto.BookResponse, error) {
	book, err := s.BookRepo.Get(ctx, &model.Book{
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

	return &dto.BookResponse{
		ID:     book.ID,
		Title:  book.Title,
		Isbn:   book.Isbn,
		Writer: book.Writer,
	}, nil
}

func (s *service) Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.BookResponse], error) {
	books, info, err := s.BookRepo.FindAll(ctx, payload, &payload.Pagination)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	data := make([]dto.BookResponse, 0)

	for _, book := range books {
		data = append(data, dto.BookResponse{
			ID:     book.ID,
			Title:  book.Title,
			Isbn:   book.Isbn,
			Writer: book.Writer,
		})
	}

	return &pkgdto.SearchGetResponse[dto.BookResponse]{
		Data:           data,
		PaginationInfo: *info,
	}, nil
}

func (s *service) Update(ctx context.Context, payload *dto.BookRequest) (*dto.BookResponse, error) {
	book, err := s.BookRepo.Get(ctx, &model.Book{Base: model.Base{ID: payload.ID}})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	if book.ID != payload.ID {
		return nil, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("email is already exist"))
	}

	book.Title = payload.Title
	book.Isbn = payload.Isbn
	book.Writer = payload.Writer
	book.UpdatedBy = ctx.Value("username").(string)
	_, err = s.BookRepo.Update(ctx, book)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	return &dto.BookResponse{
		ID:     book.ID,
		Title:  book.Title,
		Isbn:   book.Isbn,
		Writer: book.Writer,
	}, nil
}

func (s *service) Delete(ctx context.Context, payload *dto.BookRequest) error {
	book, err := s.BookRepo.Get(ctx, &model.Book{Base: model.Base{ID: payload.ID}})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	book.DeletedBy = ctx.Value("username").(string)
	_, err = s.BookRepo.Delete(ctx, book)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	return nil
}
