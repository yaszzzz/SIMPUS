package repository

import (
	"database/sql"
	"fmt"
	"simpus/internal/models"
)

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) FindAll(page, limit int, search string) ([]models.Member, int, error) {
	offset := (page - 1) * limit

	// Count total
	countQuery := `SELECT COUNT(*) FROM members WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		countQuery += ` AND (name LIKE ? OR email LIKE ? OR member_code LIKE ?)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := `SELECT id, member_code, name, email, password, phone, member_type, address, is_active, created_at, updated_at 
			  FROM members WHERE 1=1`

	if search != "" {
		query += ` AND (name LIKE ? OR email LIKE ? OR member_code LIKE ?)`
	}
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var m models.Member
		var phone, address sql.NullString
		err := rows.Scan(
			&m.ID, &m.MemberCode, &m.Name, &m.Email, &m.Password,
			&phone, &m.MemberType, &address, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		m.Phone = phone.String
		m.Address = address.String
		members = append(members, m)
	}
	return members, total, nil
}

func (r *MemberRepository) FindByID(id int) (*models.Member, error) {
	m := &models.Member{}
	var phone, address sql.NullString
	query := `SELECT id, member_code, name, email, password, phone, member_type, address, is_active, created_at, updated_at 
			  FROM members WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&m.ID, &m.MemberCode, &m.Name, &m.Email, &m.Password,
		&phone, &m.MemberType, &address, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	m.Phone = phone.String
	m.Address = address.String
	return m, nil
}

func (r *MemberRepository) FindByEmail(email string) (*models.Member, error) {
	m := &models.Member{}
	var phone, address sql.NullString
	query := `SELECT id, member_code, name, email, password, phone, member_type, address, is_active, created_at, updated_at 
			  FROM members WHERE email = ?`

	err := r.db.QueryRow(query, email).Scan(
		&m.ID, &m.MemberCode, &m.Name, &m.Email, &m.Password,
		&phone, &m.MemberType, &address, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	m.Phone = phone.String
	m.Address = address.String
	return m, nil
}

func (r *MemberRepository) GenerateMemberCode(memberType string) (string, error) {
	prefix := "MBR"
	switch memberType {
	case "mahasiswa":
		prefix = "MHS"
	case "guru":
		prefix = "GRU"
	case "karyawan":
		prefix = "KRY"
	}

	var count int
	query := `SELECT COUNT(*) FROM members WHERE member_type = ?`
	err := r.db.QueryRow(query, memberType).Scan(&count)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%03d", prefix, count+1), nil
}

func (r *MemberRepository) Create(m *models.MemberCreate, hashedPassword, memberCode string) (int64, error) {
	query := `INSERT INTO members (member_code, name, email, password, phone, member_type, address) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, memberCode, m.Name, m.Email, hashedPassword, m.Phone, m.MemberType, m.Address)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MemberRepository) Update(id int, m *models.MemberUpdate) error {
	query := `UPDATE members SET name = ?, email = ?, phone = ?, member_type = ?, address = ?, is_active = ? WHERE id = ?`
	_, err := r.db.Exec(query, m.Name, m.Email, m.Phone, m.MemberType, m.Address, m.IsActive, id)
	return err
}

func (r *MemberRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM members WHERE id = ?`, id)
	return err
}

func (r *MemberRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM members WHERE is_active = TRUE`).Scan(&count)
	return count, err
}
