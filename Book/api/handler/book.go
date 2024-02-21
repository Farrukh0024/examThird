package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kitab/api/models"
)

// CreateBook godoc
// @Router       /book [POST]
// @Summary      Creates a new book
// @Description  create a new book
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book body models.Create false "book"
// @Success      201  {object}  models.Book
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateBook(c *gin.Context) {
	createBook := models.Create{}

	if err := c.ShouldBindJSON(&createBook); err != nil {
		handleResponse(c, "error while reading book body from client", http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := h.services.Book().Create(ctx, createBook)
	if err != nil {
		handleResponse(c, "error while creating book", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, resp)
}

// GetBook godoc
// @Router       /book/{id} [GET]
// @Summary      Gets book
// @Description  get book by ID
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        id path string true "book"
// @Success      200  {object}  models.Book
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBook(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "invalid uuid type", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	book, err := h.services.Book().GetBook(ctx, models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, "error while getting book by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, book)
}

// GetBookList godoc
// @Router       /books [GET]
// @Summary      Get book list
// @Description  get book list
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.BookResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBookList(c *gin.Context) {
	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error while parsing page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := h.services.Book().GetBooks(ctx, models.GetListRequest{
		page,
		limit,
		search,
	})
	if err != nil {
		handleResponse(c, "error while getting books", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "success!", http.StatusOK, resp)
}

// UpdateBook godoc
// @Router       /book/{id} [PUT]
// @Summary      Update book
// @Description  update book
// @Tags         book
// @Accept       json
// @Produce      json
// @Param 		 id path string true "book_id"
// @Param        user body models.Update true "book"
// @Success      200  {object}  models.Book
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateBook(c *gin.Context) {
	updateBook := models.Update{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateBook.ID = uid

	if err := c.ShouldBindJSON(&updateBook); err != nil {
		handleResponse(c, "error while reading book body", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := h.services.Book().Update(ctx, updateBook)
	if err != nil {
		handleResponse(c, "error while updating book", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, resp)
}

// DeleteBook godoc
// @Router       /book/{id} [DELETE]
// @Summary      Delete book
// @Description  delete book
// @Tags         book
// @Accept       json
// @Produce      json
// @Param 		 id path string true "book_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteBook(c *gin.Context) {
	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = h.services.Book().Delete(ctx, models.PrimaryKey{
		ID: id.String(),
	}); err != nil {
		handleResponse(c, "error while deleting book by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data successfully deleted")
}

// UpdatePageNumber godoc
// @Router       /book/{id} [PATCH]
// @Summary      Update book page numbers
// @Description  update book page numbers
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        id path string true "book_id"
// @Param        book body models.UpdatePageNumberRequest true "book"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdatePageNumber(c *gin.Context) {
	updatePageNumber := models.UpdatePageNumberRequest{}

	if err := c.ShouldBindJSON(&updatePageNumber); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	bookID := c.Param("id")
	uid, err := uuid.Parse(bookID)
	if err != nil {
		handleResponse(c, "error while parsing UUID", http.StatusBadRequest, err.Error())
		return
	}

	updatePageNumber.ID = uid.String()

	if _, err := h.services.Book().UpdatePageNumber(context.Background(), updatePageNumber); err != nil {
		handleResponse(c, "error while updating book page numbers", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "page numbers successfully updated")
}
