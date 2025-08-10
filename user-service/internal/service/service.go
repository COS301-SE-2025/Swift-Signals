package service

import (
	"context"
	"errors"
	"strings"
	"time"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/shared/jwt"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service struct {
	repo      db.UserRepository
	validator *validator.Validate
}

func NewUserService(r db.UserRepository) UserService {
	return &Service{
		repo:      r,
		validator: validator.New(),
	}
}

func (s *Service) RegisterUser(
	ctx context.Context,
	name, email, password string,
) (*model.User, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := RegisterUserRequest{
		Name:     strings.TrimSpace(name),
		Email:    strings.TrimSpace(email),
		Password: password,
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("checking if email already exists")
	existingUser, err := s.repo.GetUserByEmail(ctx, normalizeEmail(email))
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

	logger.Debug("hashing password")
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, errs.NewExternalError("failed to hash password", err, nil)
	}

	logger.Debug("creating user")
	id := uuid.New().String()
	user := &model.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
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

func (s *Service) LoginUser(
	ctx context.Context,
	email, password string,
) (string, time.Time, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := LoginUserRequest{
		Email:    strings.TrimSpace(email),
		Password: password,
	}
	if err := s.validator.Struct(req); err != nil {
		return "", time.Time{}, handleValidationError(err)
	}

	logger.Debug("checking if email already exists")
	user, err := s.repo.GetUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return "", time.Time{}, errs.NewInternalError(
			"failed to check existing user",
			err,
			map[string]any{"email": email})
	}
	if user == nil {
		return "", time.Time{}, errs.NewInternalError(
			"user does not exist",
			err,
			map[string]any{"email": email},
		)
	}

	logger.Debug("checking if password is correct")
	err = checkPassword(password, user.Password)
	if err != nil {
		return "", time.Time{}, errs.NewUnauthorizedError(
			"password is incorrect",
			map[string]any{"user": user.PublicUser()},
		)
	}

	logger.Debug("generating token")
	role := "regular"
	if user.IsAdmin {
		role = "admin"
	}
	expiryDate := time.Now().Add(time.Hour * 72)
	token, err := jwt.GenerateToken(user.ID, role, time.Hour*72)
	if err != nil {
		return "", time.Time{}, errs.NewInternalError(
			"failed to generated token",
			err,
			map[string]any{"user": user.PublicUser()},
		)
	}

	return token, expiryDate, nil
}

func (s *Service) LogoutUser(ctx context.Context, userID string) error {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := LogoutUserRequest{
		UserID: userID,
	}
	if err := s.validator.Struct(req); err != nil {
		return handleValidationError(err)
	}

	logger.Debug("checking if user exists")
	_, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return err
		}
		return errs.NewInternalError("failed to find user", err, map[string]any{})
	}

	return nil
}

func (s *Service) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := GetUserByIDRequest{
		UserID: userID,
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("query database for user")
	user, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError(
			"failed to find user",
			err,
			map[string]any{"userID": userID},
		)
	}

	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := GetUserByEmailRequest{
		Email: email,
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("query database for user")
	user, err := s.repo.GetUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError(
			"failed to find user",
			err,
			map[string]any{"email": email},
		)
	}

	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (s *Service) GetAllUsers(
	ctx context.Context,
	page, pageSize int32,
	filter string,
) ([]*model.User, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating pagination parameters")
	req := GetAllUsersRequest{
		Page:     page,
		PageSize: pageSize,
		Filter:   filter,
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("preparing filter parameters")
	normalizedFilter := strings.TrimSpace(filter)
	offset := (page - 1) * pageSize
	limit := pageSize

	logger.Debug("querying database with pagination",
		"page", page,
		"pageSize", pageSize,
		"offset", offset,
		"filter", normalizedFilter,
	)
	// NOTE: ListUsers does not support filtering at the moment
	users, err := s.repo.ListUsers(ctx, int(limit), int(offset))
	if err != nil {
		return nil, errs.NewInternalError(
			"failed to retrieve users",
			err,
			map[string]any{
				"page":     page,
				"pageSize": pageSize,
				"filter":   normalizedFilter,
			},
		)
	}

	return users, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID, name, email string) (*model.User, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input parameters")
	req := UpdateUserRequest{
		UserID: strings.TrimSpace(userID),
		Name:   strings.TrimSpace(name),
		Email:  strings.TrimSpace(email),
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("checking if user exists")
	existingUser, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError(
			"failed to find user",
			err,
			map[string]any{"userID": req.UserID},
		)
	}

	logger.Debug("checking if email is being changed")
	normalizedEmail := normalizeEmail(req.Email)
	if existingUser.Email != normalizedEmail {
		logger.Debug("checking if new email is already taken")
		userWithEmail, err := s.repo.GetUserByEmail(ctx, normalizedEmail)
		if err != nil {
			var svcErr *errs.ServiceError
			if errors.As(err, &svcErr) {
				return nil, err
			}
			return nil, errs.NewInternalError(
				"failed to find user",
				err,
				map[string]any{"email": email},
			)
		}

		if userWithEmail != nil && userWithEmail.ID != req.UserID {
			return nil, errs.NewAlreadyExistsError(
				"email already taken by another user",
				map[string]any{"email": req.Email},
			)
		}
	}

	logger.Debug("updating user in database")
	updatedUserData := &model.User{
		ID:              existingUser.ID,
		Name:            req.Name,
		Email:           req.Email,
		Password:        existingUser.Password,
		IsAdmin:         existingUser.IsAdmin,
		IntersectionIDs: existingUser.IntersectionIDs,
		CreatedAt:       existingUser.CreatedAt,
		UpdatedAt:       time.Now(),
	}
	updatedUser, err := s.repo.UpdateUser(ctx, updatedUserData)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError(
			"failed to update user",
			err,
			map[string]any{
				"userID": req.UserID,
				"name":   req.Name,
				"email":  req.Email,
			},
		)
	}

	return updatedUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input parameters")
	req := DeleteUserRequest{
		UserID: strings.TrimSpace(userID),
	}
	if err := s.validator.Struct(req); err != nil {
		return handleValidationError(err)
	}

	logger.Debug("checking if user exists")
	existingUser, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return err
		}
		return errs.NewInternalError(
			"failed to find user",
			err,
			map[string]any{"userID": userID},
		)
	}

	if existingUser.IsAdmin {
		logger.Warn("deleting admin user")
	}

	logger.Debug("deleting user in database")
	err = s.repo.DeleteUser(ctx, userID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return err
		}
		return errs.NewInternalError(
			"failed to delete user",
			err,
			map[string]any{
				"userID": req.UserID,
			},
		)
	}

	return nil
}

func (s *Service) GetUserIntersectionIDs(ctx context.Context, userID string) ([]string, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := GetUserIntersectionIDsRequest{
		UserID: userID,
	}
	if err := s.validator.Struct(req); err != nil {
		return nil, handleValidationError(err)
	}

	logger.Debug("getting user")
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError(
			"failed to find user",
			err,
			map[string]any{"userID": userID},
		)
	}

	logger.Debug("getting user's intersection IDs")
	intIDs, err := s.repo.GetIntersectionsByUserID(ctx, userID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return nil, err
		}
		return nil, errs.NewInternalError(
			"failed to fetch user's intersection IDs",
			err,
			map[string]any{"userID": userID},
		)
	}

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
