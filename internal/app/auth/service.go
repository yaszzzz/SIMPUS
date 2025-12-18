package auth

import (
	"errors"
	"time"

	"simpus/config"
	"simpus/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type MemberRepository interface {
	FindByEmail(email string) (*models.Member, error)
	GenerateMemberCode(memberType string) (string, error)
	Create(m *models.MemberCreate, hashedPassword, memberCode string) (int64, error)
}

type Service struct {
	userRepo   *Repository
	memberRepo MemberRepository
	config     *config.Config
}

func NewService(userRepo *Repository, memberRepo MemberRepository, cfg *config.Config) *Service {
	return &Service{
		userRepo:   userRepo,
		memberRepo: memberRepo,
		config:     cfg,
	}
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Type     string `json:"type"` // "admin" or "member"
	jwt.RegisteredClaims
}

func (s *Service) LoginAdmin(username, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("username atau password salah")
	}

	if !user.IsActive {
		return nil, "", errors.New("akun tidak aktif")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("username atau password salah")
	}

	token, err := s.generateToken(user.ID, user.Username, user.Role, "admin")
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) LoginMember(email, password string) (*models.Member, string, error) {
	member, err := s.memberRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("email atau password salah")
	}

	if !member.IsActive {
		return nil, "", errors.New("akun tidak aktif")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(password)); err != nil {
		return nil, "", errors.New("email atau password salah")
	}

	token, err := s.generateToken(member.ID, member.Email, member.MemberType, "member")
	if err != nil {
		return nil, "", err
	}

	return member, token, nil
}

func (s *Service) RegisterMember(data *models.MemberCreate) (int64, error) {
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

func (s *Service) generateToken(userID int, username, role, userType string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Type:     userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.Expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *Service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
