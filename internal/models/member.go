package models

import "time"

type Member struct {
	ID         int       `json:"id"`
	MemberCode string    `json:"member_code"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Phone      string    `json:"phone"`
	MemberType string    `json:"member_type"`
	Address    string    `json:"address"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type MemberCreate struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	MemberType string `json:"member_type"`
	Address    string `json:"address"`
}

type MemberUpdate struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	MemberType string `json:"member_type"`
	Address    string `json:"address"`
	IsActive   bool   `json:"is_active"`
}
