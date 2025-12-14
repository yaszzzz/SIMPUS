package books

import (
	"simpus/internal/models"
)

type Service struct {
	bookRepo     *BookRepository
	categoryRepo *CategoryRepository
	authorRepo   *AuthorRepository
}

func NewService(bookRepo *BookRepository, categoryRepo *CategoryRepository, authorRepo *AuthorRepository) *Service {
	return &Service{
		bookRepo:     bookRepo,
		categoryRepo: categoryRepo,
		authorRepo:   authorRepo,
	}
}

func (s *Service) GetBooks(filter models.BookFilter) ([]models.Book, int, error) {
	return s.bookRepo.FindAll(filter)
}

func (s *Service) GetBook(id int) (*models.Book, error) {
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

func (s *Service) CreateBook(data *models.BookCreate) (int64, error) {
	return s.bookRepo.Create(data)
}

func (s *Service) UpdateBook(id int, data *models.BookUpdate) error {
	return s.bookRepo.Update(id, data)
}

func (s *Service) DeleteBook(id int) error {
	return s.bookRepo.Delete(id)
}

func (s *Service) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

func (s *Service) GetCategory(id int) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *Service) CreateCategory(data *models.CategoryCreate) (int64, error) {
	return s.categoryRepo.Create(data)
}

func (s *Service) UpdateCategory(id int, data *models.CategoryCreate) error {
	return s.categoryRepo.Update(id, data)
}

func (s *Service) DeleteCategory(id int) error {
	return s.categoryRepo.Delete(id)
}

func (s *Service) GetAuthors() ([]models.Author, error) {
	return s.authorRepo.FindAll()
}

func (s *Service) GetAuthor(id int) (*models.Author, error) {
	return s.authorRepo.FindByID(id)
}

func (s *Service) CreateAuthor(data *models.AuthorCreate) (int64, error) {
	return s.authorRepo.Create(data)
}

func (s *Service) UpdateAuthor(id int, data *models.AuthorCreate) error {
	return s.authorRepo.Update(id, data)
}

func (s *Service) DeleteAuthor(id int) error {
	return s.authorRepo.Delete(id)
}

func (s *Service) GetStats() (totalBooks int, availableBooks int, err error) {
	totalBooks, err = s.bookRepo.Count()
	if err != nil {
		return
	}
	availableBooks, err = s.bookRepo.CountAvailable()
	return
}

func (s *Service) GetBookCount() (int, error) {
	return s.bookRepo.Count()
}
