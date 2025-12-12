package services

import (
	"errors"
	"os"
	"project-management-backend/models"
	"project-management-backend/repositories"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(user *models.User) error
	Login(email, password string) (string, *models.User, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserFromToken(tokenString string) (*models.User, error)
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var (
	ErrEmailExists  = errors.New("email sudah terdaftar")
	ErrHashPassword = errors.New("gagal mengenkripsi password")
	ErrSystemError  = errors.New("terjadi kesalahan sistem")
)

func (s *authService) Register(user *models.User) error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	existingUser, err := s.userRepo.GetByEmail(user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else {
			return ErrSystemError
		}
	} else if existingUser != nil {
		return ErrEmailExists
	}

	hashedPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return ErrHashPassword
	}
	user.Password = hashedPassword

	if user.Role == "" {
		user.Role = "member"
	}

	if err := s.userRepo.Create(user); err != nil {
		return ErrSystemError
	}

	return nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("email atau password salah")
	}

	if !s.CheckPassword(password, user.Password) {
		return "", nil, errors.New("email atau password salah")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, errors.New("gagal generate token")
	}

	return token, user, nil
}

func (s *authService) generateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(8 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(jwtSecret))
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	return jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
}

func (s *authService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("token tidak valid")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		user, err := s.userRepo.GetByID(claims.UserID)
		if err != nil {
			return nil, errors.New("user tidak ditemukan")
		}
		return user, nil
	}

	return nil, errors.New("token tidak valid")
}

func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *authService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
