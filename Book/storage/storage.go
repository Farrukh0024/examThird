package storage

import (
	"context"

	"kitab/api/models"
)

type IStorage interface {
	Close()
	Book() IBookStorage
}

type IBookStorage interface {
	Create(context.Context, models.Create) (string, error)
	GetByID(context.Context, models.PrimaryKey) (models.Book, error)
	GetList(context.Context, models.GetListRequest) (models.BookResponse, error)
	Update(context.Context, models.Update) (string, error)
	Delete(context.Context, models.PrimaryKey) error
	UpdatePageNumber(context.Context, models.UpdatePageNumberRequest) (string, error)
}
