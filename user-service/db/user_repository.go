package db

import (
	"context"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/models"
	"sync"
	"time"
)

// inMemoryUserRepo is an in-memory implementation of UserRepository
// This is useful for testing or development environments
type inMemoryUserRepo struct {
	mu     sync.RWMutex            // Protects concurrent access to users map
	users  map[string]*models.User // Key: user ID, Value: user
	emails map[string]string       // Key: email, Value: user ID (for email lookups)
}

// NewUserRepository creates a new in-memory user repository
func NewUserRepository() models.UserRepository {
	return &inMemoryUserRepo{
		users:  make(map[string]*models.User),
		emails: make(map[string]string),
	}
}

// CreateUser creates a new user in the repository
func (r *inMemoryUserRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user with this email already exists
	if _, exists := r.emails[user.Email]; exists {
		return nil, models.ErrUserExists
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Store user
	r.users[user.ID] = user
	r.emails[user.Email] = user.ID

	// Return a copy to avoid external modifications
	return r.copyUser(user), nil
}

// GetUserByID retrieves a user by their ID
func (r *inMemoryUserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, models.ErrUserNotFound
	}

	return r.copyUser(user), nil
}

// GetUserByEmail retrieves a user by their email address
func (r *inMemoryUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.emails[email]
	if !exists {
		return nil, models.ErrUserNotFound
	}

	user, exists := r.users[userID]
	if !exists {
		// This should never happen if data is consistent
		return nil, models.ErrUserNotFound
	}

	return r.copyUser(user), nil
}

// UpdateUser updates an existing user
func (r *inMemoryUserRepo) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingUser, exists := r.users[user.ID]
	if !exists {
		return nil, models.ErrUserNotFound
	}

	// Check if email is being changed and if new email already exists
	if existingUser.Email != user.Email {
		if _, emailExists := r.emails[user.Email]; emailExists {
			return nil, models.ErrUserExists
		}

		// Remove old email mapping
		delete(r.emails, existingUser.Email)
		// Add new email mapping
		r.emails[user.Email] = user.ID
	}

	// Update timestamp
	user.UpdatedAt = time.Now()
	user.CreatedAt = existingUser.CreatedAt // Preserve original creation time

	// Update user
	r.users[user.ID] = user

	return r.copyUser(user), nil
}

// DeleteUser removes a user from the repository
func (r *inMemoryUserRepo) DeleteUser(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return models.ErrUserNotFound
	}

	// Remove from both maps
	delete(r.users, id)
	delete(r.emails, user.Email)

	return nil
}

// ListUsers returns a paginated list of users
func (r *inMemoryUserRepo) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	allUsers := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		allUsers = append(allUsers, r.copyUser(user))
	}

	// Simple pagination
	start := offset
	if start > len(allUsers) {
		return []*models.User{}, nil
	}

	end := start + limit
	if end > len(allUsers) {
		end = len(allUsers)
	}

	return allUsers[start:end], nil
}

// copyUser creates a deep copy of a user to prevent external modifications
func (r *inMemoryUserRepo) copyUser(user *models.User) *models.User {
	return &models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Additional helper methods for testing

// GetUserCount returns the total number of users (useful for testing)
func (r *inMemoryUserRepo) GetUserCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.users)
}

// Clear removes all users from the repository (useful for testing)
func (r *inMemoryUserRepo) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = make(map[string]*models.User)
	r.emails = make(map[string]string)
}
