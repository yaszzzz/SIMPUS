package books

import (
	"database/sql"
	"simpus/internal/models"
)

type AuthorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) FindAll() ([]models.Author, error) {
	query := `SELECT a.id, a.name, a.bio, a.created_at, COUNT(b.id) as book_count
			  FROM authors a
			  LEFT JOIN books b ON a.id = b.author_id
			  GROUP BY a.id
			  ORDER BY a.name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []models.Author
	for rows.Next() {
		var a models.Author
		var bio sql.NullString
		err := rows.Scan(&a.ID, &a.Name, &bio, &a.CreatedAt, &a.BookCount)
		if err != nil {
			return nil, err
		}
		a.Bio = bio.String
		authors = append(authors, a)
	}
	return authors, nil
}

func (r *AuthorRepository) FindByID(id int) (*models.Author, error) {
	a := &models.Author{}
	var bio sql.NullString
	query := `SELECT id, name, bio, created_at FROM authors WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(&a.ID, &a.Name, &bio, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	a.Bio = bio.String
	return a, nil
}

func (r *AuthorRepository) Create(a *models.AuthorCreate) (int64, error) {
	query := `INSERT INTO authors (name, bio) VALUES (?, ?)`
	result, err := r.db.Exec(query, a.Name, a.Bio)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *AuthorRepository) Update(id int, a *models.AuthorCreate) error {
	query := `UPDATE authors SET name = ?, bio = ? WHERE id = ?`
	_, err := r.db.Exec(query, a.Name, a.Bio, id)
	return err
}

func (r *AuthorRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM authors WHERE id = ?`, id)
	return err
}
