package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/shared/jwt"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/util"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo db.UserRepository
}

func NewUserService(r db.UserRepository) UserService {
	return &Service{repo: r}
}

// emailRegex is a simple regex for basic email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// RegisterUser creates a new user with proper validation and password hashing
func (s *Service) RegisterUser(
	ctx context.Context,
	name, email, password string,
) (*model.User, error) {
	logger := util.LoggerFromContext(ctx)

	// Validate input before using db resources
	logger.Debug("validaing input")
	email = normalizeEmail(email)
	if err := s.validateUserInput(name, email, password); err != nil {
		return nil, err
	}

	// Check if user already exists
	logger.Debug("checking if email already exists")
	existingUser, err := s.repo.GetUserByEmail(ctx, email)
	// NOTE: Logic is dependent on GetUserByEmail returning nil if user does not exist
	//       If this returns an error instead, we need to handle it differently
	//       This is a limitation of the current implementation
	//       Perhaps we should define EmailExists repository method instead
	if err != nil {
		return nil, errs.NewInternalError(
			"failed to check existing user",
			err,
			map[string]any{"email": email},
		)
	}

	if existingUser != nil {
		return nil, errs.NewAlreadyExistsError(
			"email already exists",
			map[string]any{"user": existingUser, "email": email},
		)
	}

	// Hash password
	logger.Debug("hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.NewInternalError("failed to hash password", err, nil)
	}

	// Create user
	logger.Debug("creating user")
	id := uuid.New().String()
	user := &model.User{
		ID:       id,
		Name:     strings.TrimSpace(name),
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Password: string(hashedPassword),
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError("failed to create user", err, map[string]any{})
	}
	return createdUser, nil
}

// validateUserInput validates the input parameters for user registration
func (s *Service) validateUserInput(name, email, password string) error {
	var validationErrors []string

	if strings.TrimSpace(name) == "" {
		validationErrors = append(validationErrors, "name is required")
	}
	if email == "" || !emailRegex.MatchString(email) {
		validationErrors = append(validationErrors, "email is invalid")
	}
	if len(password) < 8 {
		validationErrors = append(validationErrors, "password is too short")
	}

	if len(validationErrors) > 0 {
		combinedErrors := strings.Join(validationErrors, "; ")
		return errs.NewValidationError(combinedErrors, map[string]any{"email": email})
	}
	return nil
}

func checkPassword(inputPassword, storedHashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(inputPassword))
}

// LoginUser authenticates a user and returns auth token
func (s *Service) LoginUser(
	ctx context.Context,
	email, password string,
) (string, time.Time, error) {
	// TODO: Implement user login
	// - Validate input parameters

	// Find user by email
	user, err := s.repo.GetUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return "", time.Time{}, err
	}
	if user == nil {
		return "", time.Time{}, errors.New("user does not exist")
	}

	// Verify password
	err = checkPassword(password, user.Password)
	if err != nil {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	// Generate JWT token
	role := "regular"
	if user.IsAdmin {
		role = "admin"
	}
	expiryDate := time.Now().Add(time.Hour * 72)
	token, err := jwt.GenerateToken(user.ID, role, time.Hour*72)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiryDate, nil
}

// LogoutUser invalidates the user's session/token
func (s *Service) LogoutUser(ctx context.Context, userID string) error {
	// TODO: Implement user logout
	// - Invalidate user session/token
	// - Clear any cached user data
	// - Log logout event
	return nil
}

// GetUserByID retrieves a user by their ID
func (s *Service) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	// TODO: Validate user ID

	// Query database for user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Return user model or not found error
	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	// TODO: Implement get user by email
	// - Validate email format
	// - Query database for user by email
	// - Return user model or not found error
	return nil, nil
}

// GetAllUsers retrieves all users with pagination and filtering
func (s *Service) GetAllUsers(
	ctx context.Context,
	page, pageSize int32,
	filter string,
) ([]*model.User, error) {
	// TODO: Implement get all users
	// - Validate pagination parameters
	// - Apply filters if provided
	// - Query database with pagination
	// - Return slice of user model
	return nil, nil
}

// UpdateUser updates user information
func (s *Service) UpdateUser(ctx context.Context, userID, name, email string) (*model.User, error) {
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
func (s *Service) GetUserIntersectionIDs(ctx context.Context, userID string) ([]string, error) {
	// TODO: Validate user ID

	// Check if user exists
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Query database for user's intersection IDs
	intIDs, err := s.repo.GetIntersectionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// - Return slice of intersection IDs
	return intIDs, nil
}

// AddIntersectionID adds an intersection ID to a user's list
func (s *Service) AddIntersectionID(
	ctx context.Context,
	userID string,
	intersectionID string,
) error {
	// TODO: Validate user ID and intersection ID

	// Check if user exists
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Check if intersection ID already exists for user
	intIDs, err := s.repo.GetIntersectionsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	for _, intID := range intIDs {
		if intID == intersectionID {
			return errors.New("intersection already exists")
		}
	}

	// Add intersection ID to user's list
	err = s.repo.AddIntersectionID(ctx, userID, intersectionID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveIntersectionID removes an intersection ID from a user's list
func (s *Service) RemoveIntersectionIDs(
	ctx context.Context,
	userID string,
	intersectionID []string,
) error {
	// TODO: Implement remove intersection ID
	// - Validate user ID and intersection ID
	// - Check if user exists
	// - Remove intersection ID from user's list
	// - Update database
	return nil
}

// ChangePassword updates a user's password
func (s *Service) ChangePassword(
	ctx context.Context,
	userID, currentPassword, newPassword string,
) error {
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
