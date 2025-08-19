package client

import (
	"context"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(client userpb.UserServiceClient) *UserClient {
	return &UserClient{
		client: client,
	}
}

func NewUserClientFromConn(conn *grpc.ClientConn) *UserClient {
	return NewUserClient(userpb.NewUserServiceClient(conn))
}

func (uc *UserClient) RegisterUser(
	ctx context.Context,
	name, email, password string,
) (*userpb.UserResponse, error) {
	req := &userpb.RegisterUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := uc.client.RegisterUser(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return user, nil
}

func (uc *UserClient) LoginUser(
	ctx context.Context,
	email, password string,
) (*userpb.LoginUserResponse, error) {
	req := &userpb.LoginUserRequest{
		Email:    email,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.LoginUser(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) LogoutUser(ctx context.Context, userID string) (*emptypb.Empty, error) {
	req := &userpb.UserIDRequest{
		UserId: userID,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.LogoutUser(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) GetUserByID(
	ctx context.Context,
	userID string,
) (*userpb.UserResponse, error) {
	req := &userpb.UserIDRequest{
		UserId: userID,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.GetUserByID(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) GetUserByEmail(
	ctx context.Context,
	email string,
) (*userpb.UserResponse, error) {
	req := &userpb.GetUserByEmailRequest{
		Email: email,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.GetUserByEmail(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) GetAllUsers(
	ctx context.Context,
	page, page_size int32,
	filter string,
) (userpb.UserService_GetAllUsersClient, error) {
	req := &userpb.GetAllUsersRequest{
		Page:     page,
		PageSize: page_size,
		Filter:   filter,
	}

	resp, err := uc.client.GetAllUsers(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) UpdateUser(
	ctx context.Context,
	user_id, name, email string,
) (*userpb.UserResponse, error) {
	req := &userpb.UpdateUserRequest{
		UserId: user_id,
		Name:   name,
		Email:  email,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.UpdateUser(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) DeleteUser(ctx context.Context, userID string) (*emptypb.Empty, error) {
	req := &userpb.UserIDRequest{
		UserId: userID,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.DeleteUser(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) GetUserIntersectionIDs(
	ctx context.Context,
	userID string,
) (userpb.UserService_GetUserIntersectionIDsClient, error) {
	req := &userpb.UserIDRequest{
		UserId: userID,
	}

	resp, err := uc.client.GetUserIntersectionIDs(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) AddIntersectionID(
	ctx context.Context,
	userID string,
	intersection_id string,
) (*emptypb.Empty, error) {
	req := &userpb.AddIntersectionIDRequest{
		UserId:         userID,
		IntersectionId: intersection_id,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.AddIntersectionID(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

// Remove a single intersection ID
func (uc *UserClient) RemoveIntersectionID(
	ctx context.Context,
	userID string,
	intersectionID string,
) (*emptypb.Empty, error) {
	resp, err := uc.RemoveIntersectionIDs(ctx, userID, []string{intersectionID})
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

// Remove multiple intersection IDs
func (uc *UserClient) RemoveIntersectionIDs(
	ctx context.Context,
	userID string,
	intersectionIDs []string,
) (*emptypb.Empty, error) {
	req := &userpb.RemoveIntersectionIDRequest{
		UserId:         userID,
		IntersectionId: intersectionIDs,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := uc.client.RemoveIntersectionIDs(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) ChangePassword(
	ctx context.Context,
	user_id, current_password, new_password string,
) (*emptypb.Empty, error) {
	req := &userpb.ChangePasswordRequest{
		UserId:          user_id,
		CurrentPassword: current_password,
		NewPassword:     new_password,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.ChangePassword(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) ResetPassword(ctx context.Context, email string) (*emptypb.Empty, error) {
	req := &userpb.ResetPasswordRequest{
		Email: email,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.ResetPassword(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) MakeAdmin(
	ctx context.Context,
	user_id, admin_user_id string,
) (*emptypb.Empty, error) {
	req := &userpb.AdminRequest{
		UserId:      user_id,
		AdminUserId: admin_user_id,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.MakeAdmin(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

func (uc *UserClient) RemoveAdmin(
	ctx context.Context,
	user_id, admin_user_id string,
) (*emptypb.Empty, error) {
	req := &userpb.AdminRequest{
		UserId:      user_id,
		AdminUserId: admin_user_id,
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := uc.client.RemoveAdmin(ctx, req)
	if err != nil {
		return nil, util.GrpcErrorToErr(err)
	}
	return resp, nil
}

// NOTE: Creates stub for testing
type UserClientInterface interface {
	RegisterUser(ctx context.Context, name, email, password string) (*userpb.UserResponse, error)
	LoginUser(ctx context.Context, email, password string) (*userpb.LoginUserResponse, error)
	LogoutUser(ctx context.Context, userID string) (*emptypb.Empty, error)
	GetUserByID(ctx context.Context, userID string) (*userpb.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*userpb.UserResponse, error)
	GetAllUsers(
		ctx context.Context,
		page, page_size int32,
		filter string,
	) (userpb.UserService_GetAllUsersClient, error)
	UpdateUser(ctx context.Context, user_id, name, email string) (*userpb.UserResponse, error)
	DeleteUser(ctx context.Context, userID string) (*emptypb.Empty, error)
	GetUserIntersectionIDs(
		ctx context.Context,
		userID string,
	) (userpb.UserService_GetUserIntersectionIDsClient, error)
	AddIntersectionID(
		ctx context.Context,
		userID string,
		intersection_id string,
	) (*emptypb.Empty, error)
	RemoveIntersectionID(
		ctx context.Context,
		userID string,
		intersectionID string,
	) (*emptypb.Empty, error)
	RemoveIntersectionIDs(
		ctx context.Context,
		userID string,
		intersectionIDs []string,
	) (*emptypb.Empty, error)
	ChangePassword(
		ctx context.Context,
		userID, current_password, new_password string,
	) (*emptypb.Empty, error)
	ResetPassword(ctx context.Context, email string) (*emptypb.Empty, error)
	MakeAdmin(ctx context.Context, user_id, admin_user_id string) (*emptypb.Empty, error)
	RemoveAdmin(ctx context.Context, user_id, admin_user_id string) (*emptypb.Empty, error)
}

// NOTE: Asserts Interface Implementation
var _ UserClientInterface = (*UserClient)(nil)
