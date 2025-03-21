package model

import "time"

type TransactionBook struct {
	Id         int        `json:"id" gorm:"primaryKey"`
	UserId     int        `json:"userId" gorm:"not null;size:255"`
	User       User       `json:"-" gorm:"foreignKey:UserId"`
	BookId     int        `json:"bookId" gorm:"not null;size:255"`
	Book       Book       `json:"-" gorm:"foreignKey:BookId"`
	BorrowDate time.Time  `json:"borrowDate" gorm:"not null"`
	ReturnDate *time.Time `json:"returnDate" gorm:"default:null"`
	Status     string     `json:"status" gorm:"not null;size:255"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

type GetBorrowBookDetailInput struct {
	Id string `uri:"id" binding:"required,numeric"`
}
