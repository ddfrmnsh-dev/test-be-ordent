package usecase

import (
	"errors"
	"test-be-ordent/model"
	"test-be-ordent/service"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationUseCase interface {
	LoginUser(email, password string) (string, model.User, error)
	RegisterUser(input model.InputRegister) (model.User, error)
}
type authUseCaseImpl struct {
	userUseCase UserUseCase
	jwtService  service.JwtService
}

func NewAuthUseCase(uc UserUseCase, jwtService service.JwtService) AuthenticationUseCase {
	return &authUseCaseImpl{
		userUseCase: uc,
		jwtService:  jwtService,
	}
}

func (au *authUseCaseImpl) LoginUser(email, password string) (string, model.User, error) {
	var user model.User
	var err error

	user, err = au.userUseCase.FindUserByEmail(email)
	if err != nil {
		return "", user, err
	}

	hasValidRole := false
	if user.Role == "admin" || user.Role == "member" {
		hasValidRole = true
	}

	if !hasValidRole {
		return "", user, errors.New("you are not authorized to login as admin")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", user, errors.New("invalid credentias")
	}

	token, err := au.jwtService.CreateToken(user)
	if err != nil {
		return "", user, err
	}
	return token, user, nil
}

func (au *authUseCaseImpl) RegisterUser(input model.InputRegister) (model.User, error) {
	var payload model.User

	payload.Email = input.Email
	payload.Password = input.Password
	payload.Role = input.Role
	payload.Name = input.Name

	user, err := au.userUseCase.CreateUser(payload)
	if err != nil {
		return user, err
	}
	return user, nil
}
