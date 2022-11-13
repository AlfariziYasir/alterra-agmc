package repository

import (
	"api-mvc/internal/model"
	pkgdto "api-mvc/pkg/dto"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Book interface {
	Create(ctx context.Context, book *model.Book) (*model.Book, error)
	Get(ctx context.Context, book *model.Book) (*model.Book, error)
	FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, p *pkgdto.Pagination) ([]*model.Book, *pkgdto.PaginationInfo, error)
	Update(ctx context.Context, book *model.Book) (*model.Book, error)
	Delete(ctx context.Context, book *model.Book) (*model.Book, error)
}

type book struct {
	DB *gorm.DB
}

func NewBookRepository(DB *gorm.DB) Book {
	return &book{DB}
}

func (r *book) Create(ctx context.Context, book *model.Book) (*model.Book, error) {
	err := r.DB.Create(&book).First(&book).Error
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *book) Get(ctx context.Context, book *model.Book) (*model.Book, error) {
	err := r.DB.Where(&book).First(&book).Error
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *book) FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]*model.Book, *pkgdto.PaginationInfo, error) {
	users := make([]*model.Book, 0)
	var count int64

	query := r.DB.Model(&model.Book{})

	if payload.Search != "" {
		search := fmt.Sprintf("%%%s%%", strings.ToLower(payload.Search))
		query = query.Where("lower(title) LIKE ? or lower(isbn) LIKE ? or lower(writer) LIKE ?", search, search)
	}

	countquery := query
	if err := countquery.Count(&count).Error; err != nil {
		return nil, nil, err
	}

	limit, offset := pkgdto.GetLimitOffset(pagination)

	err := query.Limit(limit).Offset(offset).Find(&users).Error

	return users, pkgdto.CheckInfoPagination(pagination, count), err
}

func (r *book) Update(ctx context.Context, book *model.Book) (*model.Book, error) {
	err := r.DB.Model(&model.Book{}).Where("id = ?", book.ID).Updates(book).Error
	if err != nil {
		return nil, err
	}

	book, err = r.Get(ctx, &model.Book{Base: model.Base{ID: book.ID}})
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *book) Delete(ctx context.Context, book *model.Book) (*model.Book, error) {
	err := r.DB.Delete(&book).Error
	if err != nil {
		return nil, err
	}

	return book, nil
}
