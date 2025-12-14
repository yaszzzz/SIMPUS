package services

import (
	"simpus/internal/models"
	"simpus/internal/repository"
)

type BookService struct {
	bookRepo     *repository.BookRepository
	categoryRepo *repository.CategoryRepository
	authorRepo   *repository.AuthorRepository
}

func NewBookService(bookRepo *repository.BookRepository, categoryRepo *repository.CategoryRepository, authorRepo *repository.AuthorRepository) *BookService {
	return &BookService{
		bookRepo:     bookRepo,
		categoryRepo: categoryRepo,
		authorRepo:   authorRepo,
	}
}

func (s *BookService) GetBooks(filter models.BookFilter) ([]models.Book, int, error) {
	return s.bookRepo.FindAll(filter)
}

func (s *BookService) GetBook(id int) (*models.Book, error) {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Load relations
	if book.CategoryID != nil {
		category, _ := s.categoryRepo.FindByID(*book.CategoryID)
		book.Category = category
	}
	if book.AuthorID != nil {
		author, _ := s.authorRepo.FindByID(*book.AuthorID)
		book.Author = author
	}

	return book, nil
}

func (s *BookService) CreateBook(data *models.BookCreate) (int64, error) {
	return s.bookRepo.Create(data)
}

func (s *BookService) UpdateBook(id int, data *models.BookUpdate) error {
	return s.bookRepo.Update(id, data)
}

func (s *BookService) DeleteBook(id int) error {
	return s.bookRepo.Delete(id)
}

func (s *BookService) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

func (s *BookService) GetCategory(id int) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *BookService) CreateCategory(data *models.CategoryCreate) (int64, error) {
	return s.categoryRepo.Create(data)
}

func (s *BookService) UpdateCategory(id int, data *models.CategoryCreate) error {
	return s.categoryRepo.Update(id, data)
}

func (s *BookService) DeleteCategory(id int) error {
	return s.categoryRepo.Delete(id)
}

func (s *BookService) GetAuthors() ([]models.Author, error) {
	return s.authorRepo.FindAll()
}

func (s *BookService) GetAuthor(id int) (*models.Author, error) {
	return s.authorRepo.FindByID(id)
}

func (s *BookService) CreateAuthor(data *models.AuthorCreate) (int64, error) {
	return s.authorRepo.Create(data)
}

func (s *BookService) UpdateAuthor(id int, data *models.AuthorCreate) error {
	return s.authorRepo.Update(id, data)
}

func (s *BookService) DeleteAuthor(id int) error {
	return s.authorRepo.Delete(id)
}

func (s *BookService) GetStats() (totalBooks int, availableBooks int, err error) {
	totalBooks, err = s.bookRepo.Count()
	if err != nil {
		return
	}
	availableBooks, err = s.bookRepo.CountAvailable()
	return
}
