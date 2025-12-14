package models

import "time"

type Book struct {
	ID          int       `json:"id"`
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	CategoryID  *int      `json:"category_id"`
	AuthorID    *int      `json:"author_id"`
	Publisher   string    `json:"publisher"`
	PublishYear int       `json:"publish_year"`
	Stock       int       `json:"stock"`
	Available   int       `json:"available"`
	CoverImage  string    `json:"cover_image"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Category *Category `json:"category,omitempty"`
	Author   *Author   `json:"author,omitempty"`
}

type BookCreate struct {
	ISBN        string `json:"isbn"`
	Title       string `json:"title"`
	CategoryID  int    `json:"category_id"`
	AuthorID    int    `json:"author_id"`
	Publisher   string `json:"publisher"`
	PublishYear int    `json:"publish_year"`
	Stock       int    `json:"stock"`
	CoverImage  string `json:"cover_image"`
	Description string `json:"description"`
}

type BookUpdate struct {
	ISBN        string `json:"isbn"`
	Title       string `json:"title"`
	CategoryID  int    `json:"category_id"`
	AuthorID    int    `json:"author_id"`
	Publisher   string `json:"publisher"`
	PublishYear int    `json:"publish_year"`
	Stock       int    `json:"stock"`
	CoverImage  string `json:"cover_image"`
	Description string `json:"description"`
}

type BookFilter struct {
	Search     string
	CategoryID int
	AuthorID   int
	Available  bool
	Page       int
	Limit      int
}
