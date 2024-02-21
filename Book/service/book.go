package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"kitab/api/models"
	"kitab/storage"
)

type bookService struct {
	storage storage.IStorage
}

func NewBookService(storage storage.IStorage) bookService {
	return bookService{
		storage: storage,
	}
}

func (b bookService) Create(ctx context.Context, createBook models.Create) (models.Book, error) {
	fmt.Println("Book create service layer", "createBook", createBook)

	pKey, err := b.storage.Book().Create(ctx, createBook)
	if err != nil {
		fmt.Println("error while creating book", err)
		return models.Book{}, err
	}

	book, err := b.storage.Book().GetByID(ctx, models.PrimaryKey{
		ID: pKey,
	})

	return book, nil
}

func (b bookService) GetBook(ctx context.Context, pKey models.PrimaryKey) (models.Book, error) {
	book, err := b.storage.Book().GetByID(ctx, pKey)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("ERROR in service layer while getting book by id", err.Error())
			return models.Book{}, err
		}
	}

	return book, nil
}

func (b bookService) GetBooks(ctx context.Context, request models.GetListRequest) (models.BookResponse, error) {
	fmt.Println("Get book list service layer", "request", request)
	booksResponse, err := b.storage.Book().GetList(ctx, request)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("error while getting books list", err)
			return models.BookResponse{}, err
		}
	}

	return booksResponse, err
}

func (b bookService) Update(ctx context.Context, updateUser models.Update) (models.Book, error) {
	pKey, err := b.storage.Book().Update(ctx, updateUser)
	if err != nil {
		fmt.Println("ERROR in service layer while updating updateBook", err)
		return models.Book{}, err
	}

	book, err := b.storage.Book().GetByID(ctx, models.PrimaryKey{
		ID: pKey,
	})
	if err != nil {
		fmt.Println("ERROR in service layer while getting books after update", err)
		return models.Book{}, err
	}

	return book, nil
}

func (b bookService) Delete(ctx context.Context, key models.PrimaryKey) error {
	err := b.storage.Book().Delete(ctx, key)
	return err
}

func (b bookService) UpdatePageNumber(ctx context.Context, updatePageNumber models.UpdatePageNumberRequest) (models.Book, error) {
	pKey, err := b.storage.Book().UpdatePageNumber(ctx, updatePageNumber)
	if err != nil {
		fmt.Printf("Error in service layer while patching book: %v\n", err)
		return models.Book{}, err
	}

	updatedBookPageNumber, err := b.storage.Book().GetByID(ctx, models.PrimaryKey{
		ID: pKey,
	})
	if err != nil {
		fmt.Printf("Error in service layer while getting book after patch: %v\n", err)
		return models.Book{}, err
	}

	return updatedBookPageNumber, nil
}
