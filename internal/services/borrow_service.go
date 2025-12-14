package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"simpus/internal/models"
	"simpus/internal/repository"
)

type BorrowService struct {
	borrowRepo *repository.BorrowingRepository
	bookRepo   *repository.BookRepository
	memberRepo *repository.MemberRepository
	notifRepo  *repository.NotificationRepository
}

func NewBorrowService(
	borrowRepo *repository.BorrowingRepository,
	bookRepo *repository.BookRepository,
	memberRepo *repository.MemberRepository,
	notifRepo *repository.NotificationRepository,
) *BorrowService {
	return &BorrowService{
		borrowRepo: borrowRepo,
		bookRepo:   bookRepo,
		memberRepo: memberRepo,
		notifRepo:  notifRepo,
	}
}

func (s *BorrowService) GetBorrowings(filter models.BorrowingFilter) ([]models.Borrowing, int, error) {
	return s.borrowRepo.FindAll(filter)
}

func (s *BorrowService) GetBorrowing(id int) (*models.Borrowing, error) {
	return s.borrowRepo.FindByID(id)
}

func (s *BorrowService) CreateBorrowing(data *models.BorrowingCreate, userID int) (int64, error) {
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
	id, err := s.borrowRepo.Create(data, userID, dueDate)
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

func (s *BorrowService) ReturnBook(id int) (*models.Borrowing, error) {
	borrowing, err := s.borrowRepo.FindByID(id)
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

	err = s.borrowRepo.Return(id, returnData)
	if err != nil {
		return nil, err
	}

	// Update book availability
	err = s.bookRepo.UpdateAvailable(borrowing.BookID, 1)
	if err != nil {
		return nil, err
	}

	// Get updated borrowing
	borrowing, _ = s.borrowRepo.FindByID(id)
	return borrowing, nil
}

func (s *BorrowService) GetActiveCount() (int, error) {
	return s.borrowRepo.CountActive()
}

func (s *BorrowService) GetOverdueCount() (int, error) {
	return s.borrowRepo.CountOverdue()
}

func (s *BorrowService) GetMemberBorrowings(memberID int) ([]models.Borrowing, error) {
	return s.borrowRepo.GetMemberBorrowings(memberID)
}

func (s *BorrowService) CheckAndCreateOverdueNotifications() (int, error) {
	overdue, err := s.borrowRepo.FindOverdue()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, br := range overdue {
		days := int(math.Ceil(time.Now().Sub(br.DueDate).Hours() / 24))
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
