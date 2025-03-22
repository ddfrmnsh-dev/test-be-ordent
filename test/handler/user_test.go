package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"test-be-ordent/handler"
	"test-be-ordent/helper"
	"test-be-ordent/model"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) FindAllUser() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserUseCase) FindUserById(id int) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserUseCase) FindUserByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserUseCase) CreateUser(input model.User) (model.User, error) {
	args := m.Called(input)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(inputId model.GetCustomerDetailInput, user model.User) (model.User, error) {
	args := m.Called(inputId, user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserUseCase) DeleteUserById(id int) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

type MockAuthMiddleware struct {
	mock.Mock
}

func (m *MockAuthMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mockTokenClaim := struct {
			UserId string
			Role   string
		}{
			UserId: "3",
			Role:   "admin",
		}

		validRole := false
		for _, role := range roles {
			if role == mockTokenClaim.Role {
				validRole = true
				break
			}
		}

		if !validRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden Resource", "status": false})
			return
		}

		newId, _ := strconv.Atoi(mockTokenClaim.UserId)
		ctx.Set("user", model.User{Id: newId, Role: mockTokenClaim.Role})

		ctx.Next()
	}
}

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userHandler := handler.NewUserHandler(mockUserUseCase, rg, mockAuthMiddleware)

	payload := model.User{
		Name:     "John Doe",
		Email:    "johndoe@example.com",
		Password: "12345678",
	}

	mockUser := model.User{
		Id:        1,
		Name:      "John Doe",
		Email:     "johndoe@example.com",
		Password:  "hashedpassword123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUserUseCase.On("CreateUser", payload).Return(mockUser, nil)

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	userHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserUseCase.AssertExpectations(t)

	var response helper.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	var actualUser model.User
	if response.Data != nil {
		actualUserBytes, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(actualUserBytes, &actualUser)
		assert.NoError(t, err)
	}

	expectedUser := mockUser

	expectedUser.CreatedAt = time.Time{}
	expectedUser.UpdatedAt = time.Time{}
	actualUser.CreatedAt = time.Time{}
	actualUser.UpdatedAt = time.Time{}

	assert.True(t, response.Status)
	assert.Equal(t, "Success to create user", response.Message)
	assert.Equal(t, expectedUser, actualUser)

}

func TestCreateUser_BadRequest(t *testing.T) {
	router := gin.Default()
	rg := router.Group("/")

	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userHandler := handler.NewUserHandler(mockUserUseCase, rg, mockAuthMiddleware)

	invalidPayload := `{ "username": "","email":"", password:"", role:""}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(invalidPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	userHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_InternalServerError(t *testing.T) {
	router := gin.Default()
	rg := router.Group("/")

	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userHandler := handler.NewUserHandler(mockUserUseCase, rg, mockAuthMiddleware)

	payload := model.User{
		Name:     "John Doe",
		Email:    "johndoe@example.com",
		Password: "12345678",
	}

	mockUserUseCase.On("CreateUser", payload).Return(model.User{}, errors.New("database error"))

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	userHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUserUseCase.AssertExpectations(t)
}
func TestFindUserById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userHandler := handler.NewUserHandler(mockUserUseCase, rg, mockAuthMiddleware)

	userID := 1
	mockUser := model.User{
		Id:    userID,
		Name:  "John Doe",
		Email: "johndoe@example.com",
		// Role:     "Admin",
	}

	mockUserUseCase.On("FindUserById", userID).Return(mockUser, nil)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/user/%d", userID), nil)
	w := httptest.NewRecorder()

	userHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUserUseCase.AssertExpectations(t)

	var response helper.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	var actualUser model.User
	if response.Data != nil {
		actualUserBytes, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(actualUserBytes, &actualUser)
		assert.NoError(t, err)
	}

	assert.True(t, response.Status)
	assert.Equal(t, "Success to get user", response.Message)
	assert.Equal(t, mockUser, actualUser)
}

func TestFindUserById_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userHandler := handler.NewUserHandler(mockUserUseCase, rg, mockAuthMiddleware)

	userID := 99
	mockUserUseCase.On("FindUserById", userID).Return(model.User{}, fmt.Errorf("User not found"))

	req, _ := http.NewRequest("GET", fmt.Sprintf("/user/%d", userID), nil)
	w := httptest.NewRecorder()

	userHandler.Route()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUserUseCase.AssertExpectations(t)

	var response helper.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response.Status)
	assert.Equal(t, "User not found", response.Message)
}
