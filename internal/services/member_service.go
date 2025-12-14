package services

import (
	"simpus/internal/models"
	"simpus/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type MemberService struct {
	memberRepo *repository.MemberRepository
}

func NewMemberService(memberRepo *repository.MemberRepository) *MemberService {
	return &MemberService{memberRepo: memberRepo}
}

func (s *MemberService) GetMembers(page, limit int, search string) ([]models.Member, int, error) {
	return s.memberRepo.FindAll(page, limit, search)
}

func (s *MemberService) GetMember(id int) (*models.Member, error) {
	return s.memberRepo.FindByID(id)
}

func (s *MemberService) CreateMember(data *models.MemberCreate) (int64, error) {
	// Generate member code
	memberCode, err := s.memberRepo.GenerateMemberCode(data.MemberType)
	if err != nil {
		return 0, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	return s.memberRepo.Create(data, string(hashedPassword), memberCode)
}

func (s *MemberService) UpdateMember(id int, data *models.MemberUpdate) error {
	return s.memberRepo.Update(id, data)
}

func (s *MemberService) DeleteMember(id int) error {
	return s.memberRepo.Delete(id)
}

func (s *MemberService) GetMemberCount() (int, error) {
	return s.memberRepo.Count()
}
