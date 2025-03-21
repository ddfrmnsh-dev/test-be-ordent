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

type BorrowBookHandler struct {
	borrowBookUseCase usecase.BorrowBookUseCase
	rg                *gin.RouterGroup
	authMiddleware    middleware.AuthMiddleware
}

func NewBorrowBookHandler(borrowBookUseCase usecase.BorrowBookUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *BorrowBookHandler {
	return &BorrowBookHandler{
		borrowBookUseCase: borrowBookUseCase,
		rg:                rg,
		authMiddleware:    authMiddleware,
	}
}

func (bb *BorrowBookHandler) Route() {
	bb.rg.GET("/borrowBooks", bb.authMiddleware.RequireToken("admin"), bb.getAllBorrowBook)
	bb.rg.POST("/borrowBooks", bb.authMiddleware.RequireToken("admin", "member"), bb.createBorrowBook)
	bb.rg.GET("/borrowBook/:id", bb.authMiddleware.RequireToken("admin"), bb.getBorrowBookById)
	bb.rg.PUT("/borrowBook/:id/returned", bb.authMiddleware.RequireToken("admin"), bb.updateBorrowBook)
}

func (bb *BorrowBookHandler) getAllBorrowBook(c *gin.Context) {
	borrowBooks, err := bb.borrowBookUseCase.FindAllBorrowBook()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse("Failed to retrieve data borrow books"))
		return
	}
	if len(borrowBooks) > 0 {
		c.JSON(http.StatusOK, helper.APIResponse("Success to get data borrow books", borrowBooks))
		return
	}
	c.JSON(http.StatusOK, helper.APIErrorResponse("List borrow books is empty"))
}

func (bb *BorrowBookHandler) createBorrowBook(c *gin.Context) {
	// Ambil data user dari context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helper.APIErrorResponse("User not found"))
		return
	}

	// Pastikan user adalah tipe struct User
	user, ok := userInterface.(model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse("Invalid user format"))
		return
	}

	var borrowBook model.TransactionBook
	if err := c.ShouldBindJSON(&borrowBook); err != nil {
		c.JSON(http.StatusBadRequest, helper.APIErrorResponse("Failed to bind data"))
		return
	}
	borrowBook, err := bb.borrowBookUseCase.CreateBorrowBook(user, borrowBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.APIResponse("Success to create borrow book", borrowBook))
}

func (bb *BorrowBookHandler) getBorrowBookById(c *gin.Context) {
	idBorrowBook := c.Param("id")
	id, err := strconv.Atoi(idBorrowBook)

	borrowBook, err := bb.borrowBookUseCase.FindBorrowBookById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse("Failed to retrieve data borrow book"))
		return
	}
	if borrowBook.Id == 0 {
		c.JSON(http.StatusNotFound, helper.APIErrorResponse("Borrow book not found"))
		return
	}
	c.JSON(http.StatusOK, helper.APIResponse("Success to get data borrow book", borrowBook))
}

func (bb *BorrowBookHandler) updateBorrowBook(c *gin.Context) {
	var inputId model.GetBorrowBookDetailInput
	if err := c.ShouldBindUri(&inputId); err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	borrowBook, err := bb.borrowBookUseCase.UpdateBorrowBook(inputId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse("Failed to update borrow book"))
		return
	}
	c.JSON(http.StatusOK, helper.APIResponse("Success to update borrow book", borrowBook))
}
