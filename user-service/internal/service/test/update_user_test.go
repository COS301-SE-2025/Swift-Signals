package test

import (
	"context"
	"errors"
	"time"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestUpdateUser_Success() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "updated@example.com"

	existingUser := &model.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "hashedpassword",
		IsAdmin:   false,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}

	updatedUser := &model.User{
		ID:        userID,
		Name:      name,
		Email:     email,
		Password:  "hashedpassword",
		IsAdmin:   false,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: time.Now(),
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && u.Name == name && u.Email == email
	})).Return(updatedUser, nil)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Equal(userID, result.ID)
	suite.Equal(name, result.Name)
	suite.Equal(email, result.Email)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_InvalidUserID() {
	userID := "invalid-uuid"
	name := "Updated Name"
	email := "updated@example.com"

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"userid": "UserID must be a valid UUID",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_EmptyUserID() {
	userID := ""
	name := "Updated Name"
	email := "updated@example.com"

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"userid": "UserID is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_InvalidName() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "A"
	email := "updated@example.com"

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"name": "Name must be at least 2 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_NameTooLong() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "This is a very long name that exceeds the maximum allowed length of 100 characters and should trigger validation error when properly implemented with validation logic"
	email := "updated@example.com"

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"name": "Name must be at most 100 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_EmptyName() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := ""
	email := "updated@example.com"

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"name": "Name is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_InvalidEmail() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "invalid-email-format"

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"email": "Invalid email format",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_EmptyEmail() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := ""

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"email": "Email is required",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_EmailTooLong() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "thisisaverylongemailaddressthatexceedsthemaximumallowedlengthof255charactersandshouldtriggeravalidationerrorbutthisisstilltooshortsoneedtostillsaythingstomakeupfortheextracharactersiamnotcountingallofthesecharactersbutihopethatthisissufficient@verylongdomainname.com"

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrValidation, svcError.Code)
	suite.Equal("invalid input", svcError.Message)
	expectedErrors := map[string]string{
		"email": "Email must be at most 255 characters long",
	}
	suite.Equal(map[string]any{"validation errors": expectedErrors}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_UserNotFound() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "updated@example.com"
	userNotFoundError := errs.NewNotFoundError("user not found", map[string]any{})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, userNotFoundError)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_EmailAlreadyTaken() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "existing@example.com"

	existingUser := &model.User{
		ID:    userID,
		Name:  "Old Name",
		Email: "old@example.com",
	}

	userWithEmail := &model.User{
		ID:    "550e8400-e29b-41d4-a716-446655440001",
		Name:  "Other User",
		Email: email,
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(userWithEmail, nil)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrAlreadyExists, svcError.Code)
	suite.Equal("email already taken by another user", svcError.Message)
	suite.Equal(map[string]any{"email": email}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_EmailAlreadyTakenServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "existing@example.com"
	userNotFoundError := errs.NewNotFoundError("user not found", map[string]any{})

	existingUser := &model.User{
		ID:    userID,
		Name:  "Old Name",
		Email: "old@example.com",
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, userNotFoundError)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_EmailAlreadyTakenError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "existing@example.com"
	repoError := errors.New("repo error")

	existingUser := &model.User{
		ID:    userID,
		Name:  "Old Name",
		Email: "old@example.com",
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, repoError)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_SameEmailAsCurrentUser() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "current@example.com"

	existingUser := &model.User{
		ID:    userID,
		Name:  "Old Name",
		Email: email,
	}

	updatedUser := &model.User{
		ID:    userID,
		Name:  name,
		Email: email,
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && u.Name == name && u.Email == email
	})).Return(updatedUser, nil)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Equal(userID, result.ID)
	suite.Equal(name, result.Name)
	suite.Equal(email, result.Email)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_RepositoryGetUserError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "updated@example.com"
	repoError := errors.New("database connection failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, repoError)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to find user", svcError.Message)
	suite.Equal(map[string]any{"userID": userID}, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_RepositoryUpdateError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "updated@example.com"

	existingUser := &model.User{
		ID:    userID,
		Name:  "Old Name",
		Email: "old@example.com",
	}

	repoError := errors.New("database update failed")

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, repoError)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrInternal, svcError.Code)
	suite.Equal("failed to update user", svcError.Message)

	expectedContext := map[string]any{
		"userID": userID,
		"name":   name,
		"email":  email,
	}
	suite.Equal(expectedContext, svcError.Context)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_RepositoryUpdateServiceError() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "updated@example.com"

	existingUser := &model.User{
		ID:    userID,
		Name:  "Old Name",
		Email: "old@example.com",
	}

	updateError := errs.NewUnauthorizedError("failed to update user", map[string]any{})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, updateError)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrUnauthorized, svcError.Code)
	suite.Equal("failed to update user", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_ServiceErrorPropagation() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Updated Name"
	email := "updated@example.com"
	serviceError := errs.NewNotFoundError("user not found", map[string]any{"user_id": userID})

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(nil, serviceError)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Nil(result)

	svcError, ok := err.(*errs.ServiceError)
	suite.True(ok)
	suite.Equal(errs.ErrNotFound, svcError.Code)
	suite.Equal("user not found", svcError.Message)

	suite.repo.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateUser_NameTrimming() {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	name := "  Updated Name  "
	email := "updated@example.com"

	existingUser := &model.User{
		ID:    userID,
		Name:  "Old Name",
		Email: "old@example.com",
	}

	updatedUser := &model.User{
		ID:    userID,
		Name:  "Updated Name",
		Email: email,
	}

	suite.repo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	suite.repo.On("GetUserByEmail", mock.Anything, email).Return(nil, nil)
	suite.repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.ID == userID && u.Name == "Updated Name" && u.Email == email
	})).Return(updatedUser, nil)

	ctx := context.Background()

	result, err := suite.service.UpdateUser(ctx, userID, name, email)

	suite.Require().NoError(err)
	suite.NotNil(result)
	suite.Equal("Updated Name", result.Name)

	suite.repo.AssertExpectations(suite.T())
}
