package models

import "time"

type Borrowing struct {
	ID         int        `json:"id"`
	MemberID   int        `json:"member_id"`
	BookID     int        `json:"book_id"`
	UserID     *int       `json:"user_id"`
	BorrowDate time.Time  `json:"borrow_date"`
	DueDate    time.Time  `json:"due_date"`
	ReturnDate *time.Time `json:"return_date"`
	Status     string     `json:"status"`
	Fine       float64    `json:"fine"`
	Notes      string     `json:"notes"`
	CreatedAt  time.Time  `json:"created_at"`

	// Relations
	Member *Member `json:"member,omitempty"`
	Book   *Book   `json:"book,omitempty"`
	User   *User   `json:"user,omitempty"`
}

type BorrowingCreate struct {
	MemberID   int    `json:"member_id"`
	BookID     int    `json:"book_id"`
	BorrowDays int    `json:"borrow_days"`
	Notes      string `json:"notes"`
}

type BorrowingReturn struct {
	ReturnDate time.Time `json:"return_date"`
	Fine       float64   `json:"fine"`
	Notes      string    `json:"notes"`
}

type BorrowingFilter struct {
	MemberID int
	BookID   int
	Status   string
	FromDate time.Time
	ToDate   time.Time
	Page     int
	Limit    int
}

// Fine calculation: Rp 1000 per day
const FinePerDay = 1000.0
