package model

import "time"

type TransactionBook struct {
	Id         int        `json:"id" gorm:"primaryKey"`
	UserId     int        `json:"userId" gorm:"not null;size:255"`
	User       User       `gorm:"foreignKey:UserId"`
	BookId     int        `json:"bookId" gorm:"not null;size:255"`
	Book       Book       `gorm:"foreignKey:BookId"`
	BorrowDate time.Time  `json:"borrowDate" gorm:"not null"`
	ReturnDate *time.Time `json:"returnDate" gorm:"default:null"`
	Status     int        `json:"status" gorm:"not null;size:255"`
	CreatedAt  string     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  string     `json:"updatedAt" gorm:"autoUpdateTime"`
}
