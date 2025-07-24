package handler

import (
	"context"

	userpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler interface {
	RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.UserResponse, error)
	LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.LoginUserResponse, error)
	LogoutUser(ctx context.Context, req *userpb.UserIDRequest) (*emptypb.Empty, error)
	GetUserByID(ctx context.Context, req *userpb.UserIDRequest) (*userpb.UserResponse, error)
	GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.UserResponse, error)
	GetAllUsers(req *userpb.GetAllUsersRequest, stream userpb.UserService_GetAllUsersServer) error
	UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error)
	DeleteUser(ctx context.Context, req *userpb.UserIDRequest) (*emptypb.Empty, error)
	GetUserIntersectionIDs(req *userpb.UserIDRequest, stream userpb.UserService_GetUserIntersectionIDsServer) error
	AddIntersectionID(ctx context.Context, req *userpb.AddIntersectionIDRequest) (*emptypb.Empty, error)
	RemoveIntersectionIDs(ctx context.Context, req *userpb.RemoveIntersectionIDRequest) (*emptypb.Empty, error)
	ChangePassword(ctx context.Context, req *userpb.ChangePasswordRequest) (*emptypb.Empty, error)
	ResetPassword(ctx context.Context, req *userpb.ResetPasswordRequest) (*emptypb.Empty, error)
	MakeAdmin(ctx context.Context, req *userpb.AdminRequest) (*emptypb.Empty, error)
	RemoveAdmin(ctx context.Context, req *userpb.AdminRequest) (*emptypb.Empty, error)
}
