package services

import (
	"simpus/internal/models"
	"simpus/internal/repository"
)

type NotificationService struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationService(notifRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo}
}

func (s *NotificationService) GetMemberNotifications(memberID int, limit int) ([]models.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.notifRepo.FindByMember(memberID, limit)
}

func (s *NotificationService) CreateNotification(data *models.NotificationCreate) (int64, error) {
	return s.notifRepo.Create(data)
}

func (s *NotificationService) MarkAsRead(id int) error {
	return s.notifRepo.MarkAsRead(id)
}

func (s *NotificationService) MarkAllAsRead(memberID int) error {
	return s.notifRepo.MarkAllAsRead(memberID)
}

func (s *NotificationService) GetUnreadCount(memberID int) (int, error) {
	return s.notifRepo.CountUnread(memberID)
}

func (s *NotificationService) DeleteNotification(id int) error {
	return s.notifRepo.Delete(id)
}
