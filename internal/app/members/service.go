package members

import (
	"simpus/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMembers(page, limit int, search string) ([]models.Member, int, error) {
	return s.repo.FindAll(page, limit, search)
}

func (s *Service) GetMember(id int) (*models.Member, error) {
	return s.repo.FindByID(id)
}

func (s *Service) CreateMember(data *models.MemberCreate) (int64, error) {
	// Generate member code
	memberCode, err := s.repo.GenerateMemberCode(data.MemberType)
	if err != nil {
		return 0, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(data, string(hashedPassword), memberCode)
}

func (s *Service) UpdateMember(id int, data *models.MemberUpdate) error {
	if data.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data.Password = string(hashedPassword)
	}
	return s.repo.Update(id, data)
}

func (s *Service) DeleteMember(id int) error {
	return s.repo.Delete(id)
}

func (s *Service) GetMemberCount() (int, error) {
	return s.repo.Count()
}
