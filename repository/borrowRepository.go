package repository

import (
	"test-be-ordent/model"

	"gorm.io/gorm"
)

type BorrowRepository interface {
	FindAll() ([]model.TransactionBook, error)
	Save(borrowBook model.TransactionBook) (model.TransactionBook, error)
	FindById(id int) (model.TransactionBook, error)
	Update(borrowBook model.TransactionBook) (model.TransactionBook, error)
	FindBorrowBookExists(userId int, bookId int) (bool, error)
}
type BorrowRepositoryImpl struct {
	db *gorm.DB
}

func NewBorrowRepository(db *gorm.DB) *BorrowRepositoryImpl {
	return &BorrowRepositoryImpl{db: db}
}

func (b *BorrowRepositoryImpl) FindAll() ([]model.TransactionBook, error) {
	var borrowBooks []model.TransactionBook
	res := b.db.Preload("Book").Preload("User").Find(&borrowBooks)
	if res.Error != nil {
		return borrowBooks, res.Error
	}
	if res.RowsAffected == 0 {
		return borrowBooks, nil
	}
	return borrowBooks, nil
}

func (b *BorrowRepositoryImpl) Save(borrowBook model.TransactionBook) (model.TransactionBook, error) {
	res := b.db.Create(&borrowBook)
	if res.Error != nil {
		return borrowBook, res.Error
	}

	var fullBorrowBook model.TransactionBook
	err := b.db.Preload("User").Preload("Book").Where("id = ?", borrowBook.Id).First(&fullBorrowBook).Error
	if err != nil {
		return model.TransactionBook{}, err
	}
	return fullBorrowBook, nil
}

func (b *BorrowRepositoryImpl) FindById(id int) (model.TransactionBook, error) {
	var borrowBook model.TransactionBook
	res := b.db.Preload("Book").Preload("User").First(&borrowBook, id)
	if res.Error != nil {
		return borrowBook, res.Error
	}
	return borrowBook, nil
}

func (b *BorrowRepositoryImpl) Update(borrowBook model.TransactionBook) (model.TransactionBook, error) {
	res := b.db.Save(&borrowBook)
	if res.Error != nil {
		return borrowBook, res.Error
	}
	return borrowBook, nil
}

func (b *BorrowRepositoryImpl) FindBorrowBookExists(userId int, bookId int) (bool, error) {
	var borrowBook model.TransactionBook
	res := b.db.Where("user_id = ? AND book_id = ?", userId, bookId).First(&borrowBook)
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
