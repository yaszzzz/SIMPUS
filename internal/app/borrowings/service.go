package borrowings

import (
	"errors"
	"fmt"
	"math"
	"time"

	"simpus/internal/app/books"
	"simpus/internal/app/members"
	"simpus/internal/models"
)

// NotificationRepository defines the interface for notification repository
// to avoid import cycle with legacy repository package if needed,
// though currently we might still import it if no cycle exists.
type NotificationRepository interface {
	Create(notif *models.NotificationCreate) (int64, error)
}

type Service struct {
	repo       *Repository
	bookRepo   *books.BookRepository
	memberRepo *members.Repository
	notifRepo  NotificationRepository
}

func NewService(
	repo *Repository,
	bookRepo *books.BookRepository,
	memberRepo *members.Repository,
	notifRepo NotificationRepository,
) *Service {
	return &Service{
		repo:       repo,
		bookRepo:   bookRepo,
		memberRepo: memberRepo,
		notifRepo:  notifRepo,
	}
}

func (s *Service) GetBorrowings(filter models.BorrowingFilter) ([]models.Borrowing, int, error) {
	return s.repo.FindAll(filter)
}

func (s *Service) GetBorrowing(id int) (*models.Borrowing, error) {
	return s.repo.FindByID(id)
}

func (s *Service) CreateBorrowing(data *models.BorrowingCreate, userID int) (int64, error) {
	// Check if book is available
	book, err := s.bookRepo.FindByID(data.BookID)
	if err != nil {
		return 0, errors.New("buku tidak ditemukan")
	}
	if book.Available <= 0 {
		return 0, errors.New("buku tidak tersedia")
	}

	// Check if member exists and is active
	member, err := s.memberRepo.FindByID(data.MemberID)
	if err != nil {
		return 0, errors.New("anggota tidak ditemukan")
	}
	if !member.IsActive {
		return 0, errors.New("anggota tidak aktif")
	}

	// Calculate due date (default 7 days)
	borrowDays := data.BorrowDays
	if borrowDays <= 0 {
		borrowDays = 7
	}
	dueDate := time.Now().AddDate(0, 0, borrowDays)

	// Create borrowing
	id, err := s.repo.Create(data, userID, dueDate)
	if err != nil {
		return 0, err
	}

	// Update book availability
	err = s.bookRepo.UpdateAvailable(data.BookID, -1)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) ReturnBook(id int) (*models.Borrowing, error) {
	borrowing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("peminjaman tidak ditemukan")
	}

	if borrowing.Status != "dipinjam" {
		return nil, errors.New("buku sudah dikembalikan")
	}

	// Calculate fine if overdue
	now := time.Now()
	var fine float64
	if now.After(borrowing.DueDate) {
		days := int(math.Ceil(now.Sub(borrowing.DueDate).Hours() / 24))
		fine = float64(days) * models.FinePerDay
	}

	returnData := &models.BorrowingReturn{
		ReturnDate: now,
		Fine:       fine,
	}

	err = s.repo.Return(id, returnData)
	if err != nil {
		return nil, err
	}

	// Update book availability
	err = s.bookRepo.UpdateAvailable(borrowing.BookID, 1)
	if err != nil {
		return nil, err
	}

	// Get updated borrowing
	borrowing, _ = s.repo.FindByID(id)
	return borrowing, nil
}

func (s *Service) GetActiveCount() (int, error) {
	return s.repo.CountActive()
}

func (s *Service) GetOverdueCount() (int, error) {
	return s.repo.CountOverdue()
}

func (s *Service) GetMemberBorrowings(memberID int) ([]models.Borrowing, error) {
	return s.repo.GetMemberBorrowings(memberID)
}

func (s *Service) CheckAndCreateOverdueNotifications() (int, error) {
	overdue, err := s.repo.FindOverdue()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, br := range overdue {
		days := int(math.Ceil(time.Since(br.DueDate).Hours() / 24))
		fine := float64(days) * models.FinePerDay

		notif := &models.NotificationCreate{
			BorrowingID: br.ID,
			MemberID:    br.MemberID,
			Type:        "keterlambatan",
			Title:       "Buku Terlambat Dikembalikan",
			Message:     fmt.Sprintf("Buku '%s' terlambat %d hari. Denda: Rp %.0f", br.Book.Title, days, fine),
		}

		_, err := s.notifRepo.Create(notif)
		if err == nil {
			count++
		}
	}

	return count, nil
}
