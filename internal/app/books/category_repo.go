package books

import (
	"database/sql"
	"simpus/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	query := `SELECT c.id, c.name, c.description, c.created_at, COUNT(b.id) as book_count
			  FROM categories c
			  LEFT JOIN books b ON c.id = b.category_id
			  GROUP BY c.id
			  ORDER BY c.name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		var desc sql.NullString
		err := rows.Scan(&c.ID, &c.Name, &desc, &c.CreatedAt, &c.BookCount)
		if err != nil {
			return nil, err
		}
		c.Description = desc.String
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) FindByID(id int) (*models.Category, error) {
	c := &models.Category{}
	var desc sql.NullString
	query := `SELECT id, name, description, created_at FROM categories WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &desc, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	c.Description = desc.String
	return c, nil
}

func (r *CategoryRepository) Create(c *models.CategoryCreate) (int64, error) {
	query := `INSERT INTO categories (name, description) VALUES (?, ?)`
	result, err := r.db.Exec(query, c.Name, c.Description)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CategoryRepository) Update(id int, c *models.CategoryCreate) error {
	query := `UPDATE categories SET name = ?, description = ? WHERE id = ?`
	_, err := r.db.Exec(query, c.Name, c.Description, id)
	return err
}

func (r *CategoryRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = ?`, id)
	return err
}
