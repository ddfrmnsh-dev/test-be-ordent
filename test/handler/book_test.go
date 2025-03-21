package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

type MockBookUseCase struct {
	mock.Mock
}

func (m *MockBookUseCase) FindAllBook() ([]model.Book, error) {
	args := m.Called()
	return args.Get(0).([]model.Book), args.Error(1)
}

func (m *MockBookUseCase) CreateBook(input model.Book) (model.Book, error) {
	args := m.Called(input)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *MockBookUseCase) FindBookById(id int) (model.Book, error) {
	args := m.Called(id)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *MockBookUseCase) UpdateBook(id model.GetBookDetailInput, user model.Book) (model.Book, error) {
	args := m.Called(id, user)
	return args.Get(0).(model.Book), args.Error(1)
}

func (m *MockBookUseCase) DeleteBookById(id int) (model.Book, error) {
	args := m.Called(id)
	return args.Get(0).(model.Book), args.Error(1)
}

func TestCreateBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	mockBookUseCase := new(MockBookUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	bookHandler := handler.NewBookHandler(mockBookUseCase, rg, mockAuthMiddleware)

	payload := model.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Year:   2023,
		Stock:  10,
	}

	mockBook := model.Book{
		Id:     1,
		Title:  "Test Book",
		Author: "Test Author",
		Year:   2023,
		Stock:  10,
	}

	mockBookUseCase.On("CreateBook", payload).Return(mockBook, nil)
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	bookHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBookUseCase.AssertExpectations(t)

	var response helper.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	var actualBook model.Book
	if response.Data != nil {
		actualBookBytes, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(actualBookBytes, &actualBook)
		assert.NoError(t, err)
	}

	expectedBook := mockBook
	expectedBook.CreatedAt = time.Time{}
	expectedBook.UpdatedAt = time.Time{}
	actualBook.CreatedAt = time.Time{}
	actualBook.UpdatedAt = time.Time{}

	assert.True(t, response.Status)
	assert.Equal(t, "Success to create book", response.Message)
	assert.Equal(t, expectedBook, actualBook)
}

func TestCreateBook_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	mockBookUseCase := new(MockBookUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	bookHandler := handler.NewBookHandler(mockBookUseCase, rg, mockAuthMiddleware)

	invalidPayload := `{ "title": "","author":"", year:0, stock:0}`
	req, _ := http.NewRequest("POST", "/books", bytes.NewBufferString(invalidPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	bookHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBook_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	mockBookUseCase := new(MockBookUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	bookHandler := handler.NewBookHandler(mockBookUseCase, rg, mockAuthMiddleware)

	payload := model.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Year:   2023,
		Stock:  10,
	}

	mockBookUseCase.On("CreateBook", payload).Return(model.Book{}, errors.New("database error"))

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	bookHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockBookUseCase.AssertExpectations(t)
}

func TestFindBookById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	mockBookUseCase := new(MockBookUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	bookHandler := handler.NewBookHandler(mockBookUseCase, rg, mockAuthMiddleware)

	bookID := 1
	mockBook := model.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Year:   2023,
		Stock:  10,
	}

	mockBookUseCase.On("FindBookById", bookID).Return(mockBook, nil)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/book/%d", bookID), nil)
	w := httptest.NewRecorder()

	bookHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBookUseCase.AssertExpectations(t)

	var response helper.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	var actualBook model.Book
	if response.Data != nil {
		actualBookBytes, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(actualBookBytes, &actualBook)
		assert.NoError(t, err)
	}

	assert.True(t, response.Status)
	assert.Equal(t, "Success to get data book", response.Message)
	assert.Equal(t, mockBook, actualBook)
}

func TestFindBookById_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	mockBookUseCase := new(MockBookUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	bookHandler := handler.NewBookHandler(mockBookUseCase, rg, mockAuthMiddleware)

	bookID := 99
	mockBookUseCase.On("FindBookById", bookID).Return(model.Book{}, fmt.Errorf("Book not found"))

	req, _ := http.NewRequest("GET", fmt.Sprintf("/book/%d", bookID), nil)
	w := httptest.NewRecorder()

	bookHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockBookUseCase.AssertExpectations(t)

	var response helper.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response.Status)
	assert.Equal(t, "Book not found", response.Message)
}
