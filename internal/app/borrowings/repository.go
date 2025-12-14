package borrowings

import (
	"database/sql"
	"simpus/internal/models"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(filter models.BorrowingFilter) ([]models.Borrowing, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	baseQuery := `FROM borrowings br
				  LEFT JOIN members m ON br.member_id = m.id
				  LEFT JOIN books b ON br.book_id = b.id
				  LEFT JOIN users u ON br.user_id = u.id
				  WHERE 1=1`
	args := []interface{}{}

	if filter.MemberID > 0 {
		baseQuery += ` AND br.member_id = ?`
		args = append(args, filter.MemberID)
	}
	if filter.BookID > 0 {
		baseQuery += ` AND br.book_id = ?`
		args = append(args, filter.BookID)
	}
	if filter.Status != "" {
		baseQuery += ` AND br.status = ?`
		args = append(args, filter.Status)
	}
	if !filter.FromDate.IsZero() {
		baseQuery += ` AND br.borrow_date >= ?`
		args = append(args, filter.FromDate)
	}
	if !filter.ToDate.IsZero() {
		baseQuery += ` AND br.borrow_date <= ?`
		args = append(args, filter.ToDate)
	}

	// Count
	var total int
	countQuery := `SELECT COUNT(*) ` + baseQuery
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := `SELECT br.id, br.member_id, br.book_id, br.user_id, br.borrow_date, 
			  br.due_date, br.return_date, br.status, br.fine, br.notes, br.created_at,
			  m.id, m.member_code, m.name, m.email, m.member_type,
			  b.id, b.isbn, b.title,
			  u.id, u.name ` + baseQuery + ` ORDER BY br.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var borrowings []models.Borrowing
	for rows.Next() {
		var br models.Borrowing
		var userID sql.NullInt64
		var returnDate sql.NullTime
		var notes sql.NullString

		var memberID int
		var memberCode, memberName, memberEmail, memberType string
		var bookID int
		var bookISBN, bookTitle sql.NullString
		var uID sql.NullInt64
		var uName sql.NullString

		err := rows.Scan(
			&br.ID, &br.MemberID, &br.BookID, &userID, &br.BorrowDate,
			&br.DueDate, &returnDate, &br.Status, &br.Fine, &notes, &br.CreatedAt,
			&memberID, &memberCode, &memberName, &memberEmail, &memberType,
			&bookID, &bookISBN, &bookTitle,
			&uID, &uName,
		)
		if err != nil {
			return nil, 0, err
		}

		if userID.Valid {
			id := int(userID.Int64)
			br.UserID = &id
		}
		if returnDate.Valid {
			br.ReturnDate = &returnDate.Time
		}
		br.Notes = notes.String

		br.Member = &models.Member{
			ID:         memberID,
			MemberCode: memberCode,
			Name:       memberName,
			Email:      memberEmail,
			MemberType: memberType,
		}
		br.Book = &models.Book{
			ID:    bookID,
			ISBN:  bookISBN.String,
			Title: bookTitle.String,
		}
		if uID.Valid {
			br.User = &models.User{
				ID:   int(uID.Int64),
				Name: uName.String,
			}
		}

		borrowings = append(borrowings, br)
	}
	return borrowings, total, nil
}

func (r *Repository) FindByID(id int) (*models.Borrowing, error) {
	br := &models.Borrowing{}
	var userID sql.NullInt64
	var returnDate sql.NullTime
	var notes sql.NullString

	query := `SELECT id, member_id, book_id, user_id, borrow_date, due_date, 
			  return_date, status, fine, notes, created_at FROM borrowings WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&br.ID, &br.MemberID, &br.BookID, &userID, &br.BorrowDate,
		&br.DueDate, &returnDate, &br.Status, &br.Fine, &notes, &br.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if userID.Valid {
		id := int(userID.Int64)
		br.UserID = &id
	}
	if returnDate.Valid {
		br.ReturnDate = &returnDate.Time
	}
	br.Notes = notes.String

	return br, nil
}

func (r *Repository) Create(br *models.BorrowingCreate, userID int, dueDate time.Time) (int64, error) {
	query := `INSERT INTO borrowings (member_id, book_id, user_id, borrow_date, due_date, status, notes) 
			  VALUES (?, ?, ?, CURDATE(), ?, 'dipinjam', ?)`

	result, err := r.db.Exec(query, br.MemberID, br.BookID, userID, dueDate, br.Notes)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) Return(id int, returnData *models.BorrowingReturn) error {
	status := "dikembalikan"
	if returnData.Fine > 0 {
		status = "terlambat"
	}

	query := `UPDATE borrowings SET return_date = ?, status = ?, fine = ?, notes = CONCAT(IFNULL(notes, ''), ?) WHERE id = ?`
	_, err := r.db.Exec(query, returnData.ReturnDate, status, returnData.Fine, returnData.Notes, id)
	return err
}

func (r *Repository) CountActive() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM borrowings WHERE status = 'dipinjam'`).Scan(&count)
	return count, err
}

func (r *Repository) CountOverdue() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM borrowings WHERE status = 'dipinjam' AND due_date < CURDATE()`).Scan(&count)
	return count, err
}

func (r *Repository) FindOverdue() ([]models.Borrowing, error) {
	query := `SELECT br.id, br.member_id, br.book_id, br.borrow_date, br.due_date, br.status,
			  m.id, m.member_code, m.name, m.email,
			  b.id, b.title
			  FROM borrowings br
			  LEFT JOIN members m ON br.member_id = m.id
			  LEFT JOIN books b ON br.book_id = b.id
			  WHERE br.status = 'dipinjam' AND br.due_date < CURDATE()
			  ORDER BY br.due_date`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var borrowings []models.Borrowing
	for rows.Next() {
		var br models.Borrowing
		var memberID int
		var memberCode, memberName, memberEmail string
		var bookID int
		var bookTitle string

		err := rows.Scan(
			&br.ID, &br.MemberID, &br.BookID, &br.BorrowDate, &br.DueDate, &br.Status,
			&memberID, &memberCode, &memberName, &memberEmail,
			&bookID, &bookTitle,
		)
		if err != nil {
			return nil, err
		}

		br.Member = &models.Member{
			ID:         memberID,
			MemberCode: memberCode,
			Name:       memberName,
			Email:      memberEmail,
		}
		br.Book = &models.Book{
			ID:    bookID,
			Title: bookTitle,
		}
		borrowings = append(borrowings, br)
	}
	return borrowings, nil
}

func (r *Repository) GetMemberBorrowings(memberID int) ([]models.Borrowing, error) {
	filter := models.BorrowingFilter{
		MemberID: memberID,
		Page:     1,
		Limit:    100,
	}
	borrowings, _, err := r.FindAll(filter)
	return borrowings, err
}
