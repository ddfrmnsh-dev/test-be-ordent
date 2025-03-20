package repository

import (
	"fmt"
	"test-be-ordent/model"

	"gorm.io/gorm"
)

type BookRepository interface {
	FindAll() ([]model.Book, error)
	Save(book model.Book) (model.Book, error)
	FindById(id int) (model.Book, error)
	Delete(id int) (model.Book, error)
	Update(book model.Book) (model.Book, error)
}

type BookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepositoryImpl {
	return &BookRepositoryImpl{db: db}
}

func (b *BookRepositoryImpl) FindAll() ([]model.Book, error) {
	var books []model.Book

	res := b.db.Find(&books)

	if res.Error != nil {
		return books, res.Error
	}

	if res.RowsAffected == 0 {
		fmt.Println("No book found")
		return books, nil
	}

	return books, nil
}

func (b *BookRepositoryImpl) Save(book model.Book) (model.Book, error) {
	res := b.db.Create(&book)
	if res.Error != nil {
		fmt.Println("Err db:", res.Error)
		return book, res.Error
	}
	return book, nil
}

func (b *BookRepositoryImpl) FindById(id int) (model.Book, error) {
	var book model.Book
	res := b.db.First(&book, id)
	if res.Error != nil {
		fmt.Println("Err db:", res.Error)
		return book, res.Error
	}
	if res.RowsAffected == 0 {
		fmt.Println("No book found")
		return book, nil
	}
	return book, nil
}

func (b *BookRepositoryImpl) Delete(id int) (model.Book, error) {
	var book model.Book
	res := b.db.Delete(&book, id)
	if res.Error != nil {
		return book, res.Error
	}
	if res.RowsAffected == 0 {
		return book, gorm.ErrRecordNotFound
	}
	return book, nil
}

func (b *BookRepositoryImpl) Update(book model.Book) (model.Book, error) {
	res := b.db.Save(&book)
	if res.Error != nil {
		return book, res.Error
	}
	return book, nil
}
