package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"test-be-ordent/handler"
	"test-be-ordent/helper"
	"test-be-ordent/model"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBorrowBookUseCase struct {
	mock.Mock
}

func (m *MockBorrowBookUseCase) FindAllBorrowBook() ([]model.TransactionBook, error) {
	args := m.Called()
	return args.Get(0).([]model.TransactionBook), args.Error(1)
}

func (m *MockBorrowBookUseCase) CreateBorrowBook(userId model.User, input model.TransactionBook) (model.TransactionBook, error) {
	args := m.Called(userId, input)
	return args.Get(0).(model.TransactionBook), args.Error(1)
}

func (m *MockBorrowBookUseCase) FindBorrowBookById(id int) (model.TransactionBook, error) {
	args := m.Called(id)
	return args.Get(0).(model.TransactionBook), args.Error(1)
}

func (m *MockBorrowBookUseCase) UpdateBorrowBook(id model.GetBorrowBookDetailInput) (model.TransactionBook, error) {
	args := m.Called(id)
	return args.Get(0).(model.TransactionBook), args.Error(1)
}

func TestCreateBorrowBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	mockBorrowBookUseCase := new(MockBorrowBookUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	borrowBookHandler := handler.NewBorrowBookHandler(mockBorrowBookUseCase, rg, mockAuthMiddleware)

	mockUser := model.User{
		Id:   3,
		Role: "admin",
	}

	router.Use(func(c *gin.Context) {
		fmt.Println("Middleware executed: setting user context")
		c.Set("user", mockUser)
		c.Next()
	})

	borrowBookHandler.Route()

	payload := model.TransactionBook{
		BookId: 7,
	}

	mockBorrowBook := model.TransactionBook{
		Id:         1,
		UserId:     mockUser.Id,
		BookId:     7,
		BorrowDate: time.Now().UTC(),
		ReturnDate: nil,
		Status:     "borrowed",
	}

	mockBorrowBookUseCase.On("CreateBorrowBook", mockUser, payload).Return(mockBorrowBook, nil)

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/borrowBooks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBorrowBookUseCase.AssertExpectations(t)

	var response helper.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	var actualBook model.TransactionBook
	if response.Data != nil {
		actualBookBytes, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(actualBookBytes, &actualBook)
		assert.NoError(t, err)
	}

	mockBorrowBook.BorrowDate = mockBorrowBook.BorrowDate.UTC()
	actualBook.BorrowDate = actualBook.BorrowDate.UTC()

	assert.True(t, response.Status)
	assert.Equal(t, "Success to create borrow book", response.Message)
}
