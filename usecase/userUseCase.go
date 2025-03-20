package usecase

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"test-be-ordent/model"
	"test-be-ordent/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	FindAllUser() ([]model.User, error)
	CreateUser(input model.User) (model.User, error)
	FindUserById(id int) (model.User, error)
	FindUserByEmail(email string) (model.User, error)
	UpdateUser(id model.GetCustomerDetailInput, user model.User) (model.User, error)
	DeleteUserById(id int) (model.User, error)
}

type userUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) UserUseCase {
	return &userUseCaseImpl{userRepository: userRepository}
}

func (uc *userUseCaseImpl) FindAllUser() ([]model.User, error) {
	users, err := uc.userRepository.FindAll()
	if err != nil {
		return users, err
	}

	return users, nil
}

func (uc *userUseCaseImpl) CreateUser(input model.User) (model.User, error) {
	user := model.User{}

	_, err := uc.userRepository.FindBySingle("email", input.Email)

	if err == nil {
		return user, errors.New("email sudah digunakan")
	}

	if !strings.Contains(input.Email, "@") {
		return user, fmt.Errorf("email tidak valid")
	}

	if len(input.Password) < 8 {
		return user, fmt.Errorf("password minimal 8 karakter")
	}

	var missingFields []string

	if strings.TrimSpace(input.Role) == "" {
		missingFields = append(missingFields, "Role")
	}
	if strings.TrimSpace(input.Name) == "" {
		missingFields = append(missingFields, "Name")
	}

	if strings.TrimSpace(input.Password) == "" {
		missingFields = append(missingFields, "Password")
	}

	if strings.TrimSpace(input.Email) == "" {
		missingFields = append(missingFields, "Email")
	}

	if len(missingFields) > 0 {
		return user, fmt.Errorf("inputan %s tidak boleh string kosong", strings.Join(missingFields, ", "))
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user.Name = input.Name
	user.Email = strings.ToLower(input.Email)
	user.Password = string(hashPassword)
	user.Role = input.Role

	saveUser, err := uc.userRepository.Save(user)

	if err != nil {
		return user, err
	}

	return saveUser, nil
}

func (uc *userUseCaseImpl) FindUserById(id int) (model.User, error) {
	user, err := uc.userRepository.FindById(id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (uc *userUseCaseImpl) FindUserByEmail(email string) (model.User, error) {
	user, err := uc.userRepository.FindBySingle("email", email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (uc *userUseCaseImpl) UpdateUser(inputId model.GetCustomerDetailInput, user model.User) (model.User, error) {
	updateId, _ := strconv.Atoi(inputId.Id)

	checkUser, err := uc.userRepository.FindById(updateId)
	if err != nil {
		return checkUser, err
	}

	if strings.TrimSpace(user.Email) != "" {
		if !strings.Contains(user.Email, "@") {
			return user, fmt.Errorf("email tidak valid")
		}

		checkUser.Email = user.Email
	}

	if strings.TrimSpace(user.Password) != "" {
		if len(user.Password) < 8 {
			return user, fmt.Errorf("password minimal 8 karakter")
		}

		if strings.TrimSpace(user.Password) == "" {
			return user, fmt.Errorf("password tidak boleh string kosong")
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return user, err
		}

		checkUser.Password = string(hashPassword)
	}

	if strings.TrimSpace(user.Role) != "" {
		checkUser.Role = user.Role
	}

	if strings.TrimSpace(user.Name) != "" {
		checkUser.Name = user.Name
	}

	checkUser.UpdatedAt = time.Now()
	updatedUser, err := uc.userRepository.Update(checkUser)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (uc *userUseCaseImpl) DeleteUserById(id int) (model.User, error) {
	checkUser, err := uc.userRepository.FindById(id)
	if err != nil {
		return checkUser, fmt.Errorf("user tidak ditemukan")
	}

	deleteUser, err := uc.userRepository.Delete(id)
	if err != nil {
		return deleteUser, err
	}

	return deleteUser, nil
}
