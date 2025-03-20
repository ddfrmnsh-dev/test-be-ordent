package model

import "time"

type User struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;size:255"`
	Email     string    `json:"email" gorm:"not null;size:255"`
	Password  string    `json:"password" gorm:"not null;size:255"`
	Role      string    `json:"role" gorm:"not null;size:50"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	TransactionBooks []TransactionBook `gorm:"foreignKey:UserId"`
}

type GetCustomerDetailInput struct {
	Id string `uri:"id" binding:"required,numeric"`
}

type InputLogin struct {
	Identifier string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type InputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UserResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsActive  *bool     `json:"isActive"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FormatUserResponse(user User) UserResponse {
	return UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
