package repository

import (
	"database/sql"
	"simpus/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) FindByMember(memberID int, limit int) ([]models.Notification, error) {
	query := `SELECT id, borrowing_id, member_id, type, title, message, is_read, created_at
			  FROM notifications WHERE member_id = ? ORDER BY created_at DESC LIMIT ?`

	rows, err := r.db.Query(query, memberID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		var borrowingID sql.NullInt64

		err := rows.Scan(&n.ID, &borrowingID, &n.MemberID, &n.Type, &n.Title, &n.Message, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		if borrowingID.Valid {
			id := int(borrowingID.Int64)
			n.BorrowingID = &id
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *NotificationRepository) Create(n *models.NotificationCreate) (int64, error) {
	query := `INSERT INTO notifications (borrowing_id, member_id, type, title, message) VALUES (?, ?, ?, ?, ?)`

	var borrowingID interface{}
	if n.BorrowingID > 0 {
		borrowingID = n.BorrowingID
	}

	result, err := r.db.Exec(query, borrowingID, n.MemberID, n.Type, n.Title, n.Message)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *NotificationRepository) MarkAsRead(id int) error {
	_, err := r.db.Exec(`UPDATE notifications SET is_read = TRUE WHERE id = ?`, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(memberID int) error {
	_, err := r.db.Exec(`UPDATE notifications SET is_read = TRUE WHERE member_id = ?`, memberID)
	return err
}

func (r *NotificationRepository) CountUnread(memberID int) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE member_id = ? AND is_read = FALSE`, memberID).Scan(&count)
	return count, err
}

func (r *NotificationRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM notifications WHERE id = ?`, id)
	return err
}
