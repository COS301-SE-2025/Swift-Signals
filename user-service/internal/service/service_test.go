package service

import (
	"context"
	"errors"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	mocks "github.com/COS301-SE-2025/Swift-Signals/user-service/internal/db/mock"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestService_RegisterUser(t *testing.T) {
	tests := []struct {
		name            string
		inputName       string
		inputEmail      string
		inputPassword   string
		setupMock       func(*mocks.MockRepository)
		expectedError   error
		expectedErrCode errs.ErrorCode
		validateResult  func(*testing.T, *model.User)
	}{
		{
			name:          "successful registration",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock: func(repo *mocks.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(nil, nil)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(&model.User{
						ID:    "test-id",
						Name:  "John Doe",
						Email: "john.doe@example.com",
					}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, user *model.User) {
				assert.NotNil(t, user)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john.doe@example.com", user.Email)
			},
		},
		{
			name:            "empty name",
			inputName:       "",
			inputEmail:      "john.doe@example.com",
			inputPassword:   "password123",
			setupMock:       func(repo *mocks.MockRepository) {},
			expectedErrCode: errs.ErrValidation,
		},
		{
			name:            "whitespace only name",
			inputName:       "   ",
			inputEmail:      "john.doe@example.com",
			inputPassword:   "password123",
			setupMock:       func(repo *mocks.MockRepository) {},
			expectedErrCode: errs.ErrValidation,
		},
		{
			name:            "invalid email format",
			inputName:       "John Doe",
			inputEmail:      "invalid-email",
			inputPassword:   "password123",
			setupMock:       func(repo *mocks.MockRepository) {},
			expectedErrCode: errs.ErrValidation,
		},
		{
			name:            "empty email",
			inputName:       "John Doe",
			inputEmail:      "",
			inputPassword:   "password123",
			setupMock:       func(repo *mocks.MockRepository) {},
			expectedErrCode: errs.ErrValidation,
		},
		{
			name:            "short password (less than 8 characters)",
			inputName:       "John Doe",
			inputEmail:      "john.doe@example.com",
			inputPassword:   "short",
			setupMock:       func(repo *mocks.MockRepository) {},
			expectedErrCode: errs.ErrValidation,
		},
		{
			name:          "user already exists",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock: func(repo *mocks.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(&model.User{Email: "john.doe@example.com"}, nil)
			},
			expectedErrCode: errs.ErrAlreadyExists,
		},
		{
			name:          "repository error on user check",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock: func(repo *mocks.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(nil, errors.New("database error"))
			},
			expectedErrCode: errs.ErrInternal,
		},
		{
			name:          "repository error on user creation",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock: func(repo *mocks.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(nil, nil)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(nil, errors.New("database error"))
			},
			expectedErrCode: errs.ErrInternal,
		},
		{
			name:          "repository service error on user creation",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock: func(repo *mocks.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(nil, nil)
				serviceErr := errs.NewDatabaseError("database constraint violation", errors.New("unique constraint"), nil)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(nil, serviceErr)
			},
			expectedErrCode: errs.ErrDatabase,
		},
		{
			name:          "email normalization",
			inputName:     "  John Doe  ",
			inputEmail:    "  JOHN.DOE@EXAMPLE.COM  ",
			inputPassword: "password123",
			setupMock: func(repo *mocks.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(nil, nil)
				repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *model.User) bool {
					return user.Name == "John Doe" && user.Email == "john.doe@example.com"
				})).Return(&model.User{
					ID:    "test-id",
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, user *model.User) {
				assert.NotNil(t, user)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john.doe@example.com", user.Email)
			},
		},
		{
			name:            "multiple validation errors",
			inputName:       "",
			inputEmail:      "invalid-email",
			inputPassword:   "short",
			setupMock:       func(repo *mocks.MockRepository) {},
			expectedErrCode: errs.ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockRepository)
			tt.setupMock(mockRepo)

			service := &Service{repo: mockRepo}
			ctx := context.Background()

			result, err := service.RegisterUser(ctx, tt.inputName, tt.inputEmail, tt.inputPassword)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else if tt.expectedErrCode != "" {
				assert.Error(t, err)
				var serviceErr *errs.ServiceError
				assert.True(t, errors.As(err, &serviceErr))
				assert.Equal(t, tt.expectedErrCode, serviceErr.Code)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_RegisterUser_PasswordHashing(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
		Return(nil, nil)

	var capturedUser *model.User
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*model.User)
		}).
		Return(&model.User{ID: "test-id"}, nil)

	service := &Service{repo: mockRepo}
	ctx := context.Background()

	password := "testpassword123"
	_, err := service.RegisterUser(ctx, "Test User", "test@example.com", password)

	assert.NoError(t, err)
	assert.NotNil(t, capturedUser)

	// Verify password was hashed
	assert.NotEqual(t, password, capturedUser.Password)

	// Verify password can be validated
	err = bcrypt.CompareHashAndPassword([]byte(capturedUser.Password), []byte(password))
	assert.NoError(t, err)

	// Verify UUID was generated
	assert.NotEmpty(t, capturedUser.ID)
	_, err = uuid.Parse(capturedUser.ID)
	assert.NoError(t, err)
}

func TestService_validateUserInput(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name            string
		inputName       string
		inputEmail      string
		inputPassword   string
		expectedError   bool
		expectedErrCode errs.ErrorCode
		expectedMessage string
	}{
		{
			name:          "valid input",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			expectedError: false,
		},
		{
			name:            "empty name",
			inputName:       "",
			inputEmail:      "john.doe@example.com",
			inputPassword:   "password123",
			expectedError:   true,
			expectedErrCode: errs.ErrValidation,
			expectedMessage: "name is required",
		},
		{
			name:            "whitespace only name",
			inputName:       "   ",
			inputEmail:      "john.doe@example.com",
			inputPassword:   "password123",
			expectedError:   true,
			expectedErrCode: errs.ErrValidation,
			expectedMessage: "name is required",
		},
		{
			name:            "invalid email",
			inputName:       "John Doe",
			inputEmail:      "invalid-email",
			inputPassword:   "password123",
			expectedError:   true,
			expectedErrCode: errs.ErrValidation,
			expectedMessage: "email is invalid",
		},
		{
			name:            "empty email",
			inputName:       "John Doe",
			inputEmail:      "",
			inputPassword:   "password123",
			expectedError:   true,
			expectedErrCode: errs.ErrValidation,
			expectedMessage: "email is invalid",
		},
		{
			name:            "short password (less than 8 characters)",
			inputName:       "John Doe",
			inputEmail:      "john.doe@example.com",
			inputPassword:   "short",
			expectedError:   true,
			expectedErrCode: errs.ErrValidation,
			expectedMessage: "password is too short",
		},
		{
			name:          "password exactly 8 characters",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password",
			expectedError: false,
		},
		{
			name:            "multiple validation errors",
			inputName:       "",
			inputEmail:      "invalid-email",
			inputPassword:   "short",
			expectedError:   true,
			expectedErrCode: errs.ErrValidation,
			expectedMessage: "name is required; email is invalid; password is too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateUserInput(tt.inputName, tt.inputEmail, tt.inputPassword)

			if tt.expectedError {
				assert.Error(t, err)
				var serviceErr *errs.ServiceError
				assert.True(t, errors.As(err, &serviceErr))
				assert.Equal(t, tt.expectedErrCode, serviceErr.Code)
				assert.Contains(t, serviceErr.Message, tt.expectedMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_normalizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase email",
			input:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "uppercase email",
			input:    "TEST@EXAMPLE.COM",
			expected: "test@example.com",
		},
		{
			name:     "mixed case email",
			input:    "TeSt@ExAmPlE.CoM",
			expected: "test@example.com",
		},
		{
			name:     "email with whitespace",
			input:    "  test@example.com  ",
			expected: "test@example.com",
		},
		{
			name:     "email with whitespace and mixed case",
			input:    "  TeSt@ExAmPlE.CoM  ",
			expected: "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeEmail(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
