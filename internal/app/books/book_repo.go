package books

import (
	"database/sql"
	"simpus/internal/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) FindAll(filter models.BookFilter) ([]models.Book, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	// Build query
	baseQuery := `FROM books b 
				  LEFT JOIN categories c ON b.category_id = c.id
				  LEFT JOIN authors a ON b.author_id = a.id
				  WHERE 1=1`
	args := []interface{}{}

	if filter.Search != "" {
		baseQuery += ` AND (b.title LIKE ? OR b.isbn LIKE ? OR a.name LIKE ?)`
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}
	if filter.CategoryID > 0 {
		baseQuery += ` AND b.category_id = ?`
		args = append(args, filter.CategoryID)
	}
	if filter.AuthorID > 0 {
		baseQuery += ` AND b.author_id = ?`
		args = append(args, filter.AuthorID)
	}
	if filter.Available {
		baseQuery += ` AND b.available > 0`
	}

	// Count
	var total int
	countQuery := `SELECT COUNT(*) ` + baseQuery
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := `SELECT b.id, b.isbn, b.title, b.category_id, b.author_id, b.publisher, 
			  b.publish_year, b.stock, b.available, b.cover_image, b.description, 
			  b.created_at, b.updated_at,
			  c.id, c.name, a.id, a.name ` + baseQuery + ` ORDER BY b.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		var categoryID, authorID, catID, authID sql.NullInt64
		var isbn, publisher, cover, desc, catName, authName sql.NullString

		err := rows.Scan(
			&b.ID, &isbn, &b.Title, &categoryID, &authorID, &publisher,
			&b.PublishYear, &b.Stock, &b.Available, &cover, &desc,
			&b.CreatedAt, &b.UpdatedAt,
			&catID, &catName, &authID, &authName,
		)
		if err != nil {
			return nil, 0, err
		}

		b.ISBN = isbn.String
		b.Publisher = publisher.String
		b.CoverImage = cover.String
		b.Description = desc.String

		if categoryID.Valid {
			id := int(categoryID.Int64)
			b.CategoryID = &id
			b.Category = &models.Category{ID: int(catID.Int64), Name: catName.String}
		}
		if authorID.Valid {
			id := int(authorID.Int64)
			b.AuthorID = &id
			b.Author = &models.Author{ID: int(authID.Int64), Name: authName.String}
		}

		books = append(books, b)
	}
	return books, total, nil
}

func (r *BookRepository) FindByID(id int) (*models.Book, error) {
	b := &models.Book{}
	var categoryID, authorID sql.NullInt64
	var isbn, publisher, cover, desc sql.NullString

	query := `SELECT id, isbn, title, category_id, author_id, publisher, 
			  publish_year, stock, available, cover_image, description, 
			  created_at, updated_at FROM books WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&b.ID, &isbn, &b.Title, &categoryID, &authorID, &publisher,
		&b.PublishYear, &b.Stock, &b.Available, &cover, &desc,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	b.ISBN = isbn.String
	b.Publisher = publisher.String
	b.CoverImage = cover.String
	b.Description = desc.String
	if categoryID.Valid {
		id := int(categoryID.Int64)
		b.CategoryID = &id
	}
	if authorID.Valid {
		id := int(authorID.Int64)
		b.AuthorID = &id
	}

	return b, nil
}

func (r *BookRepository) Create(b *models.BookCreate) (int64, error) {
	query := `INSERT INTO books (isbn, title, category_id, author_id, publisher, 
			  publish_year, stock, available, cover_image, description) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var catID, authID interface{}
	if b.CategoryID > 0 {
		catID = b.CategoryID
	}
	if b.AuthorID > 0 {
		authID = b.AuthorID
	}

	result, err := r.db.Exec(query, b.ISBN, b.Title, catID, authID, b.Publisher,
		b.PublishYear, b.Stock, b.Stock, b.CoverImage, b.Description)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *BookRepository) Update(id int, b *models.BookUpdate) error {
	query := `UPDATE books SET isbn = ?, title = ?, category_id = ?, author_id = ?, 
			  publisher = ?, publish_year = ?, stock = ?, cover_image = ?, description = ? 
			  WHERE id = ?`

	var catID, authID interface{}
	if b.CategoryID > 0 {
		catID = b.CategoryID
	}
	if b.AuthorID > 0 {
		authID = b.AuthorID
	}

	_, err := r.db.Exec(query, b.ISBN, b.Title, catID, authID, b.Publisher,
		b.PublishYear, b.Stock, b.CoverImage, b.Description, id)
	return err
}

func (r *BookRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM books WHERE id = ?`, id)
	return err
}

func (r *BookRepository) UpdateAvailable(id int, delta int) error {
	query := `UPDATE books SET available = available + ? WHERE id = ? AND available + ? >= 0`
	result, err := r.db.Exec(query, delta, id, delta)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *BookRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM books`).Scan(&count)
	return count, err
}

func (r *BookRepository) CountAvailable() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT SUM(available) FROM books`).Scan(&count)
	return count, err
}
