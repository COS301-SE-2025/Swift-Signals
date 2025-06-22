package user

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo db.UserRepository
}

func NewService(r db.UserRepository) *Service {
	return &Service{repo: r}
}

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("password must be at least 8 characters long")
	ErrInvalidName     = errors.New("name cannot be empty")
	ErrUserExists      = errors.New("user with this email already exists")
)

// emailRegex is a simple regex for basic email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// RegisterUser creates a new user with proper validation and password hashing
func (s *Service) RegisterUser(ctx context.Context, name, email, password string) (*models.User, error) {

	email = normalizeEmail(email)

	// Validate input
	if err := s.validateUserInput(name, email, password); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	id := uuid.New().String()
	user := &models.User{
		ID:       id,
		Name:     strings.TrimSpace(name),
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Password: string(hashedPassword),
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// validateUserInput validates the input parameters for user registration
func (s *Service) validateUserInput(name, email, password string) error {
	// Validate name
	if strings.TrimSpace(name) == "" {
		return ErrInvalidName
	}

	// Validate email
	email = strings.TrimSpace(email)
	if email == "" || !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	// Validate password
	if len(password) < 8 {
		return ErrInvalidPassword
	}

	return nil
}

// TODO: Implement LoginUser
func (s *Service) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user.Password != password {
		return "", errors.New("invalid credentials")
	}
	return "mock-token-for-" + user.ID, nil
}

// TODO: Implement LogoutUser
func (s *Service) LogoutUser(ctx context.Context, userID string) bool {
	return true
}
