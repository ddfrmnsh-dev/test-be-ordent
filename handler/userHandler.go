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

type UserHandler struct {
	userUseCase    usecase.UserUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewUserHandler(userUseCase usecase.UserUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *UserHandler {
	return &UserHandler{
		userUseCase:    userUseCase,
		rg:             rg,
		authMiddleware: authMiddleware,
	}
}
func (uh *UserHandler) Route() {
	uh.rg.GET("/users", uh.authMiddleware.RequireToken("admin", "guest"), uh.getAllUser)
	uh.rg.POST("/users", uh.authMiddleware.RequireToken("admin", "guest"), uh.createUser)
	uh.rg.GET("/user/:id", uh.authMiddleware.RequireToken("admin", "guest"), uh.getUserById)
	uh.rg.PUT("/user/:id", uh.authMiddleware.RequireToken("admin", "guest"), uh.updateUser)
	uh.rg.DELETE("/user/:id", uh.authMiddleware.RequireToken("admin", "guest"), uh.deleteUser)
}

func (uh *UserHandler) getAllUser(c *gin.Context) {

	users, err := uh.userUseCase.FindAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse("Failed to retrieve data users"))
		return
	}

	if len(users) > 0 {
		c.JSON(http.StatusOK, helper.APIResponse("Success to get data users", users))
		return
	}

	c.JSON(http.StatusOK, helper.APIErrorResponse("List users is empty"))
}

func (uh *UserHandler) createUser(c *gin.Context) {
	var payload model.User

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helper.APIErrorResponse(err.Error()))
		return
	}

	user, err := uh.userUseCase.CreateUser(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIResponse("Success to create user", user))
}

func (uh *UserHandler) getUserById(c *gin.Context) {
	idUser := c.Param("id")

	id, err := strconv.Atoi(idUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIErrorResponse("Invalid user ID format"))
		return
	}

	user, err := uh.userUseCase.FindUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIResponse("Success to get user", user))
}

func (uh *UserHandler) updateUser(c *gin.Context) {
	var inputId model.GetCustomerDetailInput
	err := c.ShouldBindUri(&inputId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIErrorResponse("Invalid user ID format"))
		return
	}

	var input model.User
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIErrorResponse(err.Error()))
		return
	}

	updateUser, err := uh.userUseCase.UpdateUser(inputId, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIResponse("Success to update user", updateUser))
}

func (uh *UserHandler) deleteUser(c *gin.Context) {
	var inputId model.GetCustomerDetailInput
	err := c.ShouldBindUri(&inputId)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIErrorResponse("Invalid user ID format"))
		return
	}

	newId, _ := strconv.Atoi(inputId.Id)

	deleteUser, err := uh.userUseCase.DeleteUserById(newId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.APIResponse("Success to delete user", deleteUser))
}
