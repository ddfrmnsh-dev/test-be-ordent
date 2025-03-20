package model

type Book struct {
	Id          int    `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null;size:255"`
	Description string `json:"description" gorm:"not null;size:255"`
	Author      string `json:"author" gorm:"not null;size:255"`
	Year        int    `json:"year" gorm:"not null;size:255"`
	Stock       int    `json:"stock" gorm:"not null;size:255"`
	CreatedAt   string `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   string `json:"updatedAt" gorm:"autoUpdateTime"`

	TransactionBooks []TransactionBook `gorm:"foreignKey:BookId"`
}

type GetBookDetailInput struct {
	Id string `uri:"id" binding:"required,numeric"`
}
