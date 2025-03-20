package model

import "time"

type User struct {
	Id        int       `json:"id" gorm:"primaryKey"`
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
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UserResponse struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	IsActive  *bool     `json:"isActive"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FormatUserResponse(user User) UserResponse {
	return UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
