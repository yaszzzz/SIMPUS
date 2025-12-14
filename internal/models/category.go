package models

import "time"

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	BookCount   int       `json:"book_count,omitempty"`
}

type CategoryCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
