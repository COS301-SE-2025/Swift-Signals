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

// LoginUser authenticates a user and returns auth token
func (s *Service) LoginUser(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	// TODO: Implement user login
	// - Validate input parameters
	// - Find user by email
	// - Verify password
	// - Generate JWT token
	// - Return auth response with token and user info
	return nil, nil
}

// LogoutUser invalidates the user's session/token
func (s *Service) LogoutUser(ctx context.Context, userID string) error {
	// TODO: Implement user logout
	// - Invalidate user session/token
	// - Clear any cached user data
	// - Log logout event
	return nil
}

// ValidateUser validates a JWT token and returns user info
func (s *Service) ValidateUser(ctx context.Context, token string) (*models.ValidateResponse, error) {
	// TODO: Implement token validation
	// - Parse and validate JWT token
	// - Check token expiration
	// - Extract user ID from token
	// - Verify user still exists and is active
	// - Return validation response
	return nil, nil
}

// GetUserByID retrieves a user by their ID
func (s *Service) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	// TODO: Implement get user by ID
	// - Validate user ID
	// - Query database for user
	// - Return user model or not found error
	return nil, nil
}

// GetUserByEmail retrieves a user by their email address
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// TODO: Implement get user by email
	// - Validate email format
	// - Query database for user by email
	// - Return user model or not found error
	return nil, nil
}

// GetAllUsers retrieves all users with pagination and filtering
func (s *Service) GetAllUsers(ctx context.Context, page, pageSize int32, filter string) ([]*models.User, error) {
	// TODO: Implement get all users
	// - Validate pagination parameters
	// - Apply filters if provided
	// - Query database with pagination
	// - Return slice of user models
	return nil, nil
}

// UpdateUser updates user information
func (s *Service) UpdateUser(ctx context.Context, userID, name, email string) (*models.User, error) {
	// TODO: Implement user update
	// - Validate input parameters
	// - Check if user exists
	// - Check if email is already taken by another user
	// - Update user in database
	// - Return updated user model
	return nil, nil
}

// DeleteUser removes a user from the system
func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	// TODO: Implement user deletion
	// - Validate user ID
	// - Check if user exists
	// - Perform soft delete or hard delete based on business rules
	// - Clean up related data if necessary
	return nil
}

// GetUserIntersectionIDs retrieves all intersection IDs for a user
func (s *Service) GetUserIntersectionIDs(ctx context.Context, userID string) ([]int32, error) {
	// TODO: Implement get user intersection IDs
	// - Validate user ID
	// - Check if user exists
	// - Query database for user's intersection IDs
	// - Return slice of intersection IDs
	return nil, nil
}

// AddIntersectionID adds an intersection ID to a user's list
func (s *Service) AddIntersectionID(ctx context.Context, userID string, intersectionID int32) error {
	// TODO: Implement add intersection ID
	// - Validate user ID and intersection ID
	// - Check if user exists
	// - Check if intersection ID already exists for user
	// - Add intersection ID to user's list
	// - Update database
	return nil
}

// RemoveIntersectionID removes an intersection ID from a user's list
func (s *Service) RemoveIntersectionID(ctx context.Context, userID string, intersectionID []int32) error {
	// TODO: Implement remove intersection ID
	// - Validate user ID and intersection ID
	// - Check if user exists
	// - Remove intersection ID from user's list
	// - Update database
	return nil
}

// ChangePassword updates a user's password
func (s *Service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// TODO: Implement password change
	// - Validate user ID and passwords
	// - Check if user exists
	// - Verify current password
	// - Validate new password strength
	// - Hash new password
	// - Update password in database
	return nil
}

// ResetPassword initiates password reset process
func (s *Service) ResetPassword(ctx context.Context, email string) error {
	// TODO: Implement password reset
	// - Validate email format
	// - Check if user exists with this email
	// - Generate password reset token
	// - Send password reset email
	// - Store reset token with expiration
	return nil
}

// MakeAdmin grants admin privileges to a user
func (s *Service) MakeAdmin(ctx context.Context, userID, adminUserID string) error {
	// TODO: Implement make admin
	// - Validate user IDs
	// - Check if admin user has permission to grant admin rights
	// - Check if target user exists
	// - Update user's admin status in database
	// - Log admin privilege change
	return nil
}

// RemoveAdmin revokes admin privileges from a user
func (s *Service) RemoveAdmin(ctx context.Context, userID, adminUserID string) error {
	// TODO: Implement remove admin
	// - Validate user IDs
	// - Check if admin user has permission to revoke admin rights
	// - Check if target user exists and is currently admin
	// - Update user's admin status in database
	// - Log admin privilege change
	return nil
}
