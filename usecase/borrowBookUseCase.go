package usecase

import (
	"fmt"
	"strconv"
	"test-be-ordent/model"
	"test-be-ordent/repository"
	"time"
)

type BorrowBookUseCase interface {
	FindAllBorrowBook() ([]model.TransactionBook, error)
	CreateBorrowBook(userId model.User, input model.TransactionBook) (model.TransactionBook, error)
	FindBorrowBookById(id int) (model.TransactionBook, error)
	UpdateBorrowBook(id model.GetBorrowBookDetailInput) (model.TransactionBook, error)
}

type borrowBookUseCaseImpl struct {
	borrowBookRepository repository.BorrowRepository
	bookRepository       repository.BookRepository
}

func NewBorrowBookUseCase(borrowBookRepository repository.BorrowRepository, bookRepository repository.BookRepository) BorrowBookUseCase {
	return &borrowBookUseCaseImpl{borrowBookRepository: borrowBookRepository, bookRepository: bookRepository}
}

func (bb *borrowBookUseCaseImpl) FindAllBorrowBook() ([]model.TransactionBook, error) {
	borrowBooks, err := bb.borrowBookRepository.FindAll()
	if err != nil {
		return borrowBooks, err
	}
	return borrowBooks, nil
}

func (bb *borrowBookUseCaseImpl) CreateBorrowBook(userId model.User, input model.TransactionBook) (model.TransactionBook, error) {
	borrowBookExists, _ := bb.borrowBookRepository.FindBorrowBookExists(userId.Id, input.BookId)

	if borrowBookExists {
		return model.TransactionBook{}, fmt.Errorf("Book already borrowed")
	}

	book, err := bb.bookRepository.FindById(input.BookId)
	if err != nil {
		return model.TransactionBook{}, err
	}

	if book.Stock < 1 {
		return model.TransactionBook{}, fmt.Errorf("Book is out of stock")
	}

	book.Stock -= 1

	_, err = bb.bookRepository.Update(book)

	if err != nil {
		return model.TransactionBook{}, err
	}

	borrowBook := model.TransactionBook{}
	borrowBook.BookId = input.BookId
	borrowBook.UserId = userId.Id
	borrowBook.Status = "borrowed"
	borrowBook.BorrowDate = time.Now()

	saveBorrowBook, err := bb.borrowBookRepository.Save(borrowBook)

	if err != nil {
		return saveBorrowBook, err
	}

	return saveBorrowBook, nil
}

func (bb *borrowBookUseCaseImpl) FindBorrowBookById(id int) (model.TransactionBook, error) {
	borrowBook, err := bb.borrowBookRepository.FindById(id)
	if err != nil {
		return borrowBook, err
	}
	return borrowBook, nil
}

func (bb *borrowBookUseCaseImpl) UpdateBorrowBook(input model.GetBorrowBookDetailInput) (model.TransactionBook, error) {

	newId, err := strconv.Atoi(input.Id)
	fmt.Println("newId", newId)
	checkBorrowBook, err := bb.borrowBookRepository.FindById(newId)
	if err != nil {
		return checkBorrowBook, fmt.Errorf("transaction book not found")
	}

	book, _ := bb.bookRepository.FindById(checkBorrowBook.BookId)
	if err != nil {
		return checkBorrowBook, fmt.Errorf("book not found")
	}

	book.Stock += 1

	_, err = bb.bookRepository.Update(book)

	nowDate := time.Now()

	checkBorrowBook.Status = "returned"
	checkBorrowBook.ReturnDate = &nowDate

	updatedBorrowBook, err := bb.borrowBookRepository.Update(checkBorrowBook)
	if err != nil {
		return updatedBorrowBook, err
	}
	return updatedBorrowBook, nil
}
