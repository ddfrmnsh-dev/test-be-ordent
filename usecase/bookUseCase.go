package usecase

import (
	"fmt"
	"strconv"
	"strings"
	"test-be-ordent/model"
	"test-be-ordent/repository"
)

type BookUseCase interface {
	FindAllBook() ([]model.Book, error)
	CreateBook(input model.Book) (model.Book, error)
	FindBookById(id int) (model.Book, error)
	UpdateBook(id model.GetBookDetailInput, user model.Book) (model.Book, error)
	DeleteBookById(id int) (model.Book, error)
}

type bookUseCaseImpl struct {
	bookRepository repository.BookRepository
}

func NewBookUseCase(bookRepository repository.BookRepository) BookUseCase {
	return &bookUseCaseImpl{bookRepository: bookRepository}
}

func (bc *bookUseCaseImpl) FindAllBook() ([]model.Book, error) {
	books, err := bc.bookRepository.FindAll()
	if err != nil {
		return books, err
	}

	return books, nil
}

func (bc *bookUseCaseImpl) CreateBook(input model.Book) (model.Book, error) {
	book := model.Book{}

	if input.Stock < 0 {
		return book, fmt.Errorf("stock tidak boleh kurang dari 0")
	}

	if input.Year < 1000 || input.Year > 2025 {
		return book, fmt.Errorf("year tidak boleh kurang dari 1000 atau lebih dari 2025")
	}

	var missingFields []string

	if strings.TrimSpace(input.Title) == "" {
		missingFields = append(missingFields, "Title")
	}

	if strings.TrimSpace(input.Description) == "" {
		missingFields = append(missingFields, "Description")
	}

	if strings.TrimSpace(input.Author) == "" {
		missingFields = append(missingFields, "Author")
	}

	if len(missingFields) > 0 {
		return book, fmt.Errorf("inputan %s tidak boleh string kosong", strings.Join(missingFields, ", "))
	}

	book.Title = input.Title
	book.Description = input.Description
	book.Author = input.Author
	book.Year = input.Year
	book.Stock = input.Stock

	saveBook, err := bc.bookRepository.Save(book)

	if err != nil {
		return book, err
	}

	return saveBook, nil
}

func (bc *bookUseCaseImpl) FindBookById(id int) (model.Book, error) {
	book, err := bc.bookRepository.FindById(id)
	if err != nil {
		return book, err
	}
	return book, nil
}

func (bc *bookUseCaseImpl) UpdateBook(input model.GetBookDetailInput, book model.Book) (model.Book, error) {
	updateId, _ := strconv.Atoi(input.Id)
	checkBook, err := bc.bookRepository.FindById(updateId)
	if err != nil {
		return book, err
	}

	if strings.TrimSpace(book.Title) != "" {
		checkBook.Title = book.Title
	}

	if strings.TrimSpace(book.Description) != "" {
		checkBook.Description = book.Description
	}

	if strings.TrimSpace(book.Author) != "" {
		checkBook.Author = book.Author
	}

	if book.Year != 0 {
		checkBook.Year = book.Year
	}

	if book.Stock != 0 {
		if book.Year < 1000 || book.Year > 2025 {
			return book, fmt.Errorf("year tidak boleh kurang dari 1000 atau lebih dari 2025")
		}
		checkBook.Stock = book.Stock
	}

	updateBook, err := bc.bookRepository.Update(checkBook)
	if err != nil {
		return updateBook, err
	}

	return updateBook, nil
}

func (bc *bookUseCaseImpl) DeleteBookById(id int) (model.Book, error) {
	book, err := bc.bookRepository.FindById(id)
	if err != nil {
		return book, err
	}
	book, err = bc.bookRepository.Delete(id)
	if err != nil {
		return book, err
	}
	return book, nil
}
