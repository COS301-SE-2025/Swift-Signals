package service

import (
	"context"
	"errors"
	"maps"
	"slices"
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

func (s *Service) AddIntersectionID(
	ctx context.Context,
	userID string,
	intersectionID string,
) error {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input")
	req := AddIntersectionIDRequest{
		UserID:         strings.TrimSpace(userID),
		IntersectionID: strings.TrimSpace(intersectionID),
	}
	if err := s.validator.Struct(req); err != nil {
		return handleValidationError(err)
	}

	logger.Debug("getting user")
	_, err := s.repo.GetUserByID(ctx, req.UserID)
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

	logger.Debug("checking if intersection ID already exists")
	intIDs, err := s.repo.GetIntersectionsByUserID(ctx, req.UserID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return err
		}
		return errs.NewInternalError(
			"failed to check existing intersection IDs",
			err,
			map[string]any{"userID": userID},
		)
	}
	if slices.Contains(intIDs, req.IntersectionID) {
		logger.Warn("intersection ID already exists in user's list")
		return nil
	}

	logger.Debug("Add intersection ID to user's list")
	err = s.repo.AddIntersectionID(ctx, userID, req.IntersectionID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return err
		}
		return errs.NewInternalError(
			"failed to add intersection ID",
			err,
			map[string]any{"userID": userID},
		)
	}

	return nil
}

func (s *Service) RemoveIntersectionIDs(
	ctx context.Context,
	userID string,
	intersectionID []string,
) error {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input parameters")
	req := RemoveIntersectionIDsRequest{
		UserID:          strings.TrimSpace(userID),
		IntersectionIDs: intersectionID,
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

	logger.Debug("fetching current intersections")
	currentIntersectionIDs, err := s.repo.GetIntersectionsByUserID(ctx, req.UserID)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return err
		}
		return errs.NewInternalError(
			"failed to fetch current intersections",
			err,
			map[string]any{"userID": userID},
		)
	}

	logger.Debug("removing requested intersection ids")
	currentSet := make(map[string]struct{}, len(currentIntersectionIDs))
	for _, id := range currentIntersectionIDs {
		currentSet[id] = struct{}{}
	}
	invalidSet := make(map[string]struct{})
	for _, id := range req.IntersectionIDs {
		if _, exists := currentSet[id]; !exists {
			invalidSet[id] = struct{}{}
		}
	}
	if len(invalidSet) > 0 {
		logger.Warn(
			"attempting to remove intersection IDs that are not in the current user's list",
			"invalidIDs", maps.Keys(invalidSet),
		)
	}
	requestedSet := make(map[string]struct{}, len(req.IntersectionIDs))
	for _, id := range req.IntersectionIDs {
		if _, isInvalid := invalidSet[id]; isInvalid {
			continue
		}
		requestedSet[id] = struct{}{}
	}
	var updatedIntersections []string
	for _, id := range currentIntersectionIDs {
		if _, remove := requestedSet[id]; !remove {
			updatedIntersections = append(updatedIntersections, id)
		}
	}

	logger.Debug("updating user in database")
	updatedUser := &model.User{
		ID:              existingUser.ID,
		Name:            existingUser.Name,
		Email:           existingUser.Email,
		Password:        existingUser.Password,
		IsAdmin:         existingUser.IsAdmin,
		IntersectionIDs: updatedIntersections,
		CreatedAt:       existingUser.CreatedAt,
		UpdatedAt:       time.Now(),
	}
	_, err = s.repo.UpdateUser(ctx, updatedUser)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return err
		}
		return errs.NewInternalError(
			"failed to update user",
			err,
			map[string]any{
				"userID": req.UserID,
			},
		)
	}

	return nil
}

func (s *Service) ChangePassword(
	ctx context.Context,
	userID, currentPassword, newPassword string,
) error {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("validating input parameters")
	req := ChangePasswordRequest{
		UserID:          strings.TrimSpace(userID),
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
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

	logger.Debug("authorising current password")
	err = checkPassword(currentPassword, existingUser.Password)
	if err != nil {
		return errs.NewUnauthorizedError(
			"current password is incorrect",
			map[string]any{"user": existingUser.PublicUser()},
		)
	}

	logger.Debug("hashing new password")
	newPasswordHashed, err := hashPassword(newPassword)
	if err != nil {
		return errs.NewExternalError("failed to hash password", err, nil)
	}

	logger.Debug("updating user in database")
	updatedUserData := &model.User{
		ID:              existingUser.ID,
		Name:            existingUser.Name,
		Email:           existingUser.Email,
		Password:        string(newPasswordHashed),
		IsAdmin:         existingUser.IsAdmin,
		IntersectionIDs: existingUser.IntersectionIDs,
		CreatedAt:       existingUser.CreatedAt,
		UpdatedAt:       time.Now(),
	}
	_, err = s.repo.UpdateUser(ctx, updatedUserData)
	if err != nil {
		var svcErr *errs.ServiceError
		if errors.As(err, &svcErr) {
			return err
		}
		return errs.NewInternalError(
			"failed to update user",
			err,
			map[string]any{
				"userID": req.UserID,
			},
		)
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, email string) error {
	// TODO: Implement password reset
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
