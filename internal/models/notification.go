package models

import "time"

type Notification struct {
	ID          int       `json:"id"`
	BorrowingID *int      `json:"borrowing_id"`
	MemberID    int       `json:"member_id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	Member    *Member    `json:"member,omitempty"`
	Borrowing *Borrowing `json:"borrowing,omitempty"`
}

type NotificationCreate struct {
	BorrowingID int    `json:"borrowing_id"`
	MemberID    int    `json:"member_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Message     string `json:"message"`
}
