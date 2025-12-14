package notifications

import (
	"simpus/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMemberNotifications(memberID int, limit int) ([]models.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.FindByMember(memberID, limit)
}

func (s *Service) CreateNotification(data *models.NotificationCreate) (int64, error) {
	return s.repo.Create(data)
}

func (s *Service) MarkAsRead(id int) error {
	return s.repo.MarkAsRead(id)
}

func (s *Service) MarkAllAsRead(memberID int) error {
	return s.repo.MarkAllAsRead(memberID)
}

func (s *Service) GetUnreadCount(memberID int) (int, error) {
	return s.repo.CountUnread(memberID)
}

func (s *Service) DeleteNotification(id int) error {
	return s.repo.Delete(id)
}
