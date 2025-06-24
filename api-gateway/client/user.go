package client

import (
	"context"
	"time"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(conn *grpc.ClientConn) *UserClient {
	return &UserClient{
		client: userpb.NewUserServiceClient(conn),
	}
}

func (uc *UserClient) RegisterUser(ctx context.Context, name, email, password string) (*userpb.UserResponse, error) {
	req := &userpb.RegisterUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.RegisterUser(ctx, req)
}

func (uc *UserClient) LoginUser(ctx context.Context, email, password string) (*userpb.LoginUserResponse, error) {
	req := &userpb.LoginUserRequest{
		Email:    email,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.LoginUser(ctx, req)
}

func (uc *UserClient) LogoutUser(ctx context.Context, userID string) (*emptypb.Empty, error) {
	req := &userpb.UserIDRequest{
		UserId: userID,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.LogoutUser(ctx, req)
}

func (uc *UserClient) GetUserByID(ctx context.Context, userID string) (*userpb.UserResponse, error) {
	req := &userpb.UserIDRequest{
		UserId: userID,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.GetUserByID(ctx, req)
}

func (uc *UserClient) GetUserByEmail(ctx context.Context, email string) (*userpb.UserResponse, error) {
	req := &userpb.GetUserByEmailRequest{
		Email: email,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.GetUserByEmail(ctx, req)
}

func (uc *UserClient) GetAllUsers(ctx context.Context, page, page_size int32, filter string) (userpb.UserService_GetAllUsersClient, error) {
	req := &userpb.GetAllUsersRequest{
		Page:     page,
		PageSize: page_size,
		Filter:   filter,
	}

	return uc.client.GetAllUsers(ctx, req)
}

func (uc *UserClient) UpdateUser(ctx context.Context, user_id, name, email string) (*userpb.UserResponse, error) {
	req := &userpb.UpdateUserRequest{
		UserId: user_id,
		Name:   name,
		Email:  email,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.UpdateUser(ctx, req)
}

func (uc *UserClient) DeleteUser(ctx context.Context, userID string) (*emptypb.Empty, error) {
	req := &userpb.UserIDRequest{
		UserId: userID,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.DeleteUser(ctx, req)
}

func (uc *UserClient) GetUserIntersectionIDs(ctx context.Context, userID string) (userpb.UserService_GetUserIntersectionIDsClient, error) {
	req := &userpb.UserIDRequest{
		UserId: userID,
	}

	return uc.client.GetUserIntersectionIDs(ctx, req)
}

func (uc *UserClient) AddIntersectionID(ctx context.Context, userID string, intersection_id int32) (*emptypb.Empty, error) {
	req := &userpb.AddIntersectionIDRequest{
		UserId:         userID,
		IntersectionId: intersection_id,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.AddIntersectionID(ctx, req)
}

// Remove a single intersection ID
func (uc *UserClient) RemoveIntersectionID(ctx context.Context, userID string, intersectionID int32) (*emptypb.Empty, error) {
	return uc.RemoveIntersectionIDs(ctx, userID, []int32{intersectionID})
}

// Remove multiple intersection IDs
func (uc *UserClient) RemoveIntersectionIDs(ctx context.Context, userID string, intersectionIDs []int32) (*emptypb.Empty, error) {
	req := &userpb.RemoveIntersectionIDRequest{
		UserId:         userID,
		IntersectionId: intersectionIDs,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return uc.client.RemoveIntersectionID(ctx, req)
}

func (uc *UserClient) ChangePassword(ctx context.Context, user_id, current_password, new_password string) (*emptypb.Empty, error) {
	req := &userpb.ChangePasswordRequest{
		UserId:          user_id,
		CurrentPassword: current_password,
		NewPassword:     new_password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.ChangePassword(ctx, req)
}

func (uc *UserClient) ResetPassword(ctx context.Context, email string) (*emptypb.Empty, error) {
	req := &userpb.ResetPasswordRequest{
		Email: email,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.ResetPassword(ctx, req)
}

func (uc *UserClient) MakeAdmin(ctx context.Context, user_id, admin_user_id string) (*emptypb.Empty, error) {
	req := &userpb.MakeAdminRequest{
		UserId:      user_id,
		AdminUserId: admin_user_id,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.MakeAdmin(ctx, req)
}

func (uc *UserClient) RemoveAdmin(ctx context.Context, user_id, admin_user_id string) (*emptypb.Empty, error) {
	req := &userpb.RemoveAdminRequest{
		UserId:      user_id,
		AdminUserId: admin_user_id,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return uc.client.RemoveAdmin(ctx, req)
}

// NOTE: Creates stub for testing
type UserClientInterface interface {
	RegisterUser(ctx context.Context, name, email, password string) (*userpb.UserResponse, error)
	LoginUser(ctx context.Context, email, password string) (*userpb.LoginUserResponse, error)
	LogoutUser(ctx context.Context, userID string) (*emptypb.Empty, error)
	GetUserByID(ctx context.Context, userID string) (*userpb.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*userpb.UserResponse, error)
	GetAllUsers(ctx context.Context, page, page_size int32, filter string) (userpb.UserService_GetAllUsersClient, error)
	UpdateUser(ctx context.Context, user_id, name, email string) (*userpb.UserResponse, error)
	DeleteUser(ctx context.Context, userID string) (*emptypb.Empty, error)
	GetUserIntersectionIDs(ctx context.Context, userID string) (userpb.UserService_GetUserIntersectionIDsClient, error)
	RemoveIntersectionID(ctx context.Context, userID string, intersectionID int32) (*emptypb.Empty, error)
	RemoveIntersectionIDs(ctx context.Context, userID string, intersectionIDs []int32) (*emptypb.Empty, error)
	ChangePassword(ctx context.Context, userID, current_password, new_password string) (*emptypb.Empty, error)
	ResetPassword(ctx context.Context, email string) (*emptypb.Empty, error)
	MakeAdmin(ctx context.Context, user_id, admin_user_id string) (*emptypb.Empty, error)
	RemoveAdmin(ctx context.Context, user_id, admin_user_id string) (*emptypb.Empty, error)
}

// NOTE: Asserts Interface Implementation
var _ UserClientInterface = (*UserClient)(nil)
