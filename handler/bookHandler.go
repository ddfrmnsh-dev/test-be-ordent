package handler

import (
	"net/http"
	"strconv"
	"test-be-ordent/helper"
	"test-be-ordent/middleware"
	"test-be-ordent/model"
	"test-be-ordent/usecase"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookUseCase    usecase.BookUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewBookHandler(bookUseCase usecase.BookUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *BookHandler {
	return &BookHandler{
		bookUseCase:    bookUseCase,
		rg:             rg,
		authMiddleware: authMiddleware,
	}
}

func (bh *BookHandler) Route() {
	bh.rg.GET("/books", bh.authMiddleware.RequireToken("admin", "guest"), bh.getAllBook)
	bh.rg.POST("/books", bh.authMiddleware.RequireToken("admin", "guest"), bh.createBook)
	bh.rg.GET("/book/:id", bh.authMiddleware.RequireToken("admin", "guest"), bh.getBookById)
	bh.rg.PUT("/book/:id", bh.authMiddleware.RequireToken("admin", "guest"), bh.updateBook)
	bh.rg.DELETE("/book/:id", bh.authMiddleware.RequireToken("admin", "guest"), bh.deleteBook)
}

func (bh *BookHandler) getAllBook(c *gin.Context) {
	books, err := bh.bookUseCase.FindAllBook()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse("Failed to retrieve data books"))
		return
	}

	if len(books) > 0 {
		c.JSON(http.StatusOK, helper.APIResponse("Success to get data books", books))
		return
	}

	c.JSON(http.StatusOK, helper.APIErrorResponse("List books is empty"))
}

func (bh *BookHandler) createBook(c *gin.Context) {
	var payload model.Book
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helper.APIErrorResponse(err.Error()))
		return
	}

	book, err := bh.bookUseCase.CreateBook(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIResponse("Success to create book", book))
}

func (bh *BookHandler) getBookById(c *gin.Context) {
	idBook := c.Param("id")
	id, err := strconv.Atoi(idBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIErrorResponse("Invalid user ID format"))
		return
	}

	book, err := bh.bookUseCase.FindBookById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.APIErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIResponse("Success to get data book", book))
}

func (bh *BookHandler) updateBook(c *gin.Context) {
	var inputId model.GetBookDetailInput
	if err := c.ShouldBindUri(&inputId); err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	var payload model.Book
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	book, err := bh.bookUseCase.UpdateBook(inputId, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIResponse("Success to update book", book))
}

func (bh *BookHandler) deleteBook(c *gin.Context) {
	var inputId model.GetBookDetailInput
	if err := c.ShouldBindUri(&inputId); err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	newId, _ := strconv.Atoi(inputId.Id)

	deleteBook, err := bh.bookUseCase.DeleteBookById(newId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIResponse("Success to delete book", deleteBook))
}
