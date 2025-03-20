package handler

import (
	"net/http"
	"test-be-ordent/helper"
	"test-be-ordent/model"
	"test-be-ordent/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase usecase.AuthenticationUseCase
	rg          *gin.RouterGroup
}

func NewAuthHandler(authUseCase usecase.AuthenticationUseCase, rg *gin.RouterGroup) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase, rg: rg}
}

func (ah *AuthHandler) Route() {
	ah.rg.POST("/signinAuth", ah.loginUser)
	ah.rg.POST("/registerUser", ah.registerUser)
}

func (ah *AuthHandler) loginUser(c *gin.Context) {
	var payload model.InputLogin

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	token, user, err := ah.authUseCase.LoginUser(payload.Identifier, payload.Password)
	if err != nil {
		if err.Error() == "invalid password" || err.Error() == "record not found" || err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, helper.APIErrorResponse("Invalid credentials"))
			return
		}
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}

	newToken := "Bearer" + " " + token

	formattedUser := model.FormatUserResponse(user)

	c.JSON(
		http.StatusOK,
		helper.APIResponse("Login success", gin.H{"userPrincipal": formattedUser, "token": newToken}),
	)
}

func (ah *AuthHandler) registerUser(c *gin.Context) {
	var payload model.InputRegister
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}
	user, err := ah.authUseCase.RegisterUser(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIErrorResponse(err.Error()))
		return
	}
	formattedUser := model.FormatUserResponse(user)
	c.JSON(http.StatusOK, helper.APIResponse("Register success", formattedUser))
}
