package service

import (
	"kitab/storage"
)

type IServiceManager interface {
	Book() bookService
}

type Service struct {
	bookService bookService
}

func New(storage storage.IStorage) Service {
	services := Service{}

	services.bookService = NewBookService(storage)

	return services
}

func (s Service) Book() bookService {
	return s.bookService
}
