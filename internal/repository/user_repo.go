package repository

import (
	"database/sql"
	"simpus/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, name, role, is_active, created_at, updated_at 
			  FROM users WHERE username = ?`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Name, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, name, role, is_active, created_at, updated_at 
			  FROM users WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Name, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(user *models.UserCreate, hashedPassword string) (int64, error) {
	query := `INSERT INTO users (username, email, password, name, role) VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, user.Username, user.Email, hashedPassword, user.Name, user.Role)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	query := `SELECT id, username, email, password, name, role, is_active, created_at, updated_at 
			  FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Password,
			&user.Name, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
