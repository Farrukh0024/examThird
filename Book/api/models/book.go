package models

type Book struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Author_name string `json:"author_name"`
	Page_namber int    `json:"page_namber"`
}

type Create struct {
	Name        string `json:"name"`
	Author_name string `json:"author_name"`
	Page_namber int    `json:"page_namber"`
}

type Update struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Author_name string `json:"author_name"`
	Page_namber int    `json:"page_namber"`
}

type BookResponse struct {
	Books []Book `json:"books"`
	Count int    `json:"count"`
}

type UpdatePageNumberRequest struct {
	ID          string `json:"id"`
	Page_Namber int    `json:"page_number"`
}
