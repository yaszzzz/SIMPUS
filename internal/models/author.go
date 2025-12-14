package models

import "time"

type Author struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	BookCount int       `json:"book_count,omitempty"`
}

type AuthorCreate struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}
