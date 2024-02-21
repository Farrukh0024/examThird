package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"kitab/api/models"
	"kitab/storage"
)

type bookRepo struct {
	db *pgxpool.Pool
}

func NewBookRepo(db *pgxpool.Pool) storage.IBookStorage {
	return &bookRepo{
		db: db,
	}
}

func (b *bookRepo) Create(ctx context.Context, createBook models.Create) (string, error) {

	uid := uuid.New()

	if _, err := b.db.Exec(ctx, `insert into
		books values ($1, $2, $3, $4)
		`,
		uid,
		createBook.Name,
		createBook.Author_name,
		createBook.Page_namber,
	); err != nil {
		fmt.Println("error while inserting data", err)
		return "", err
	}

	return uid.String(), nil
}

func (b *bookRepo) GetByID(ctx context.Context, pKey models.PrimaryKey) (models.Book, error) {
	book := models.Book{}

	query := `
		select id, name, author_name, page_number
			from books where id = $1
			`

	if err := b.db.QueryRow(ctx, query, pKey.ID).Scan(
		&book.ID, //0
		&book.Name,
		&book.Author_name,
		&book.Page_namber,
	); err != nil {
		fmt.Println("error while scanning book", err)
		return models.Book{}, err
	}

	return book, nil
}

func (b *bookRepo) GetList(ctx context.Context, request models.GetListRequest) (models.BookResponse, error) {
	var (
		books             = []models.Book{}
		count             = 0
		countQuery, query string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `
		SELECT count(1) from books `

	if search != "" {
		countQuery += fmt.Sprintf(` and (name ilike '%s' or author_name ilike '%s')`, search, search)
	}

	if err := b.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of books", err.Error())
		return models.BookResponse{}, err
	}

	query = `
		SELECT id, name, author_name, page_number
			FROM books
			    `

	if search != "" {
		query += fmt.Sprintf(` and (name ilike '%s' or author_name ilike '%s') `, search, search)
	}

	query += ` order by name desc LIMIT $1 OFFSET $2`

	rows, err := b.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.BookResponse{}, err
	}

	for rows.Next() {
		book := models.Book{}

		if err = rows.Scan(
			&book.ID,
			&book.Name,
			&book.Author_name,
			&book.Page_namber,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.BookResponse{}, err
		}

		books = append(books, book)
	}

	return models.BookResponse{
		Books: books,
		Count: count,
	}, nil
}

func (b *bookRepo) Update(ctx context.Context, request models.Update) (string, error) {
	query := `
		update books
			set name = $1, author_name = $2, page_number = $3
				where id = $4`

	if _, err := b.db.Exec(ctx, query,
		request.Name,
		request.Author_name,
		request.Page_namber,
		request.ID); err != nil {
		fmt.Println("error while updating book data", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (b *bookRepo) Delete(ctx context.Context, request models.PrimaryKey) error {
	query := `delete from books where id = $1`

	if _, err := b.db.Exec(ctx, query, request.ID); err != nil {
		fmt.Println("error while deleting book by id", err.Error())
		return err
	}

	return nil
}

func (b *bookRepo) UpdatePageNumber(ctx context.Context, request models.UpdatePageNumberRequest) (string, error) {
	query := `
		UPDATE books
		SET page_number = $1
		WHERE id = $2
	`

	if _, err := b.db.Exec(ctx, query, request.Page_Namber, request.ID); err != nil {
		fmt.Println("error while updating page number for book", err.Error())
		return "", err
	}

	return request.ID, nil
}
