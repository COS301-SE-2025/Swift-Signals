package user

import (
	"context"
	"errors"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestService_RegisterUser(t *testing.T) {
	tests := []struct {
		name           string
		inputName      string
		inputEmail     string
		inputPassword  string
		setupMock      func(*db.MockRepository)
		expectedError  error
		validateResult func(*testing.T, *models.User)
	}{
		{
			name:          "successful registration",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock: func(repo *db.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(nil, models.ErrUserNotFound)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(&models.User{
						ID:    "test-id",
						Name:  "John Doe",
						Email: "john.doe@example.com",
					}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, user *models.User) {
				assert.NotNil(t, user)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john.doe@example.com", user.Email)
			},
		},
		{
			name:          "empty name",
			inputName:     "",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock:     func(repo *db.MockRepository) {},
			expectedError: ErrInvalidName,
		},
		{
			name:          "whitespace only name",
			inputName:     "   ",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock:     func(repo *db.MockRepository) {},
			expectedError: ErrInvalidName,
		},
		{
			name:          "invalid email format",
			inputName:     "John Doe",
			inputEmail:    "invalid-email",
			inputPassword: "password123",
			setupMock:     func(repo *db.MockRepository) {},
			expectedError: ErrInvalidEmail,
		},
		{
			name:          "empty email",
			inputName:     "John Doe",
			inputEmail:    "",
			inputPassword: "password123",
			setupMock:     func(repo *db.MockRepository) {},
			expectedError: ErrInvalidEmail,
		},
		{
			name:          "short password",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "short",
			setupMock:     func(repo *db.MockRepository) {},
			expectedError: ErrInvalidPassword,
		},
		{
			name:          "user already exists",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock: func(repo *db.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(&models.User{Email: "john.doe@example.com"}, nil)
			},
			expectedError: ErrUserExists,
		},
		{
			name:          "repository error on user check",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock: func(repo *db.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("failed to check existing user: database error"),
		},
		{
			name:          "repository error on user creation",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			setupMock: func(repo *db.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(nil, models.ErrUserNotFound)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("failed to create user: database error"),
		},
		{
			name:          "email normalization",
			inputName:     "  John Doe  ",
			inputEmail:    "  JOHN.DOE@EXAMPLE.COM  ",
			inputPassword: "password123",
			setupMock: func(repo *db.MockRepository) {
				repo.On("GetUserByEmail", mock.Anything, "john.doe@example.com").
					Return(nil, models.ErrUserNotFound)
				repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *models.User) bool {
					return user.Name == "John Doe" && user.Email == "john.doe@example.com"
				})).Return(&models.User{
					ID:    "test-id",
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, user *models.User) {
				assert.NotNil(t, user)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john.doe@example.com", user.Email)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(db.MockRepository)
			tt.setupMock(mockRepo)

			service := &Service{repo: mockRepo}
			ctx := context.Background()

			result, err := service.RegisterUser(ctx, tt.inputName, tt.inputEmail, tt.inputPassword)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
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
	mockRepo := new(db.MockRepository)
	mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
		Return(nil, models.ErrUserNotFound)

	var capturedUser *models.User
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*models.User)
		}).
		Return(&models.User{ID: "test-id"}, nil)

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
}

func TestService_validateUserInput(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name          string
		inputName     string
		inputEmail    string
		inputPassword string
		expectedError error
	}{
		{
			name:          "valid input",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			expectedError: nil,
		},
		{
			name:          "empty name",
			inputName:     "",
			inputEmail:    "john.doe@example.com",
			inputPassword: "password123",
			expectedError: ErrInvalidName,
		},
		{
			name:          "invalid email",
			inputName:     "John Doe",
			inputEmail:    "invalid-email",
			inputPassword: "password123",
			expectedError: ErrInvalidEmail,
		},
		{
			name:          "short password",
			inputName:     "John Doe",
			inputEmail:    "john.doe@example.com",
			inputPassword: "short",
			expectedError: ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateUserInput(tt.inputName, tt.inputEmail, tt.inputPassword)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
