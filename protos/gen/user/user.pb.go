// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: user.proto

package user

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Common Messages
type UserIDRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserIDRequest) Reset() {
	*x = UserIDRequest{}
	mi := &file_user_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIDRequest) ProtoMessage() {}

func (x *UserIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIDRequest.ProtoReflect.Descriptor instead.
func (*UserIDRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *UserIDRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

// Authentication Messages
type RegisterUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email         string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterUserRequest) Reset() {
	*x = RegisterUserRequest{}
	mi := &file_user_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserRequest) ProtoMessage() {}

func (x *RegisterUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserRequest.ProtoReflect.Descriptor instead.
func (*RegisterUserRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterUserRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RegisterUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginUserRequest) Reset() {
	*x = LoginUserRequest{}
	mi := &file_user_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginUserRequest) ProtoMessage() {}

func (x *LoginUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginUserRequest.ProtoReflect.Descriptor instead.
func (*LoginUserRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *LoginUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *LoginUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// User CRUD Messages
type GetUserByEmailRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserByEmailRequest) Reset() {
	*x = GetUserByEmailRequest{}
	mi := &file_user_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserByEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByEmailRequest) ProtoMessage() {}

func (x *GetUserByEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByEmailRequest.ProtoReflect.Descriptor instead.
func (*GetUserByEmailRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserByEmailRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type GetAllUsersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int32                  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Filter        string                 `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"` // Optional filter criteria
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAllUsersRequest) Reset() {
	*x = GetAllUsersRequest{}
	mi := &file_user_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAllUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllUsersRequest) ProtoMessage() {}

func (x *GetAllUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllUsersRequest.ProtoReflect.Descriptor instead.
func (*GetAllUsersRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{4}
}

func (x *GetAllUsersRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetAllUsersRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *GetAllUsersRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

type UpdateUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUserRequest) Reset() {
	*x = UpdateUserRequest{}
	mi := &file_user_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserRequest) ProtoMessage() {}

func (x *UpdateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateUserRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UpdateUserRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

// Intersection ID Messages
type IntersectionIDResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	IntersectionId int32                  `protobuf:"varint,1,opt,name=intersection_id,json=intersectionId,proto3" json:"intersection_id,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *IntersectionIDResponse) Reset() {
	*x = IntersectionIDResponse{}
	mi := &file_user_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IntersectionIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntersectionIDResponse) ProtoMessage() {}

func (x *IntersectionIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntersectionIDResponse.ProtoReflect.Descriptor instead.
func (*IntersectionIDResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{6}
}

func (x *IntersectionIDResponse) GetIntersectionId() int32 {
	if x != nil {
		return x.IntersectionId
	}
	return 0
}

type AddIntersectionIDRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	UserId         string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	IntersectionId int32                  `protobuf:"varint,2,opt,name=intersection_id,json=intersectionId,proto3" json:"intersection_id,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *AddIntersectionIDRequest) Reset() {
	*x = AddIntersectionIDRequest{}
	mi := &file_user_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddIntersectionIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddIntersectionIDRequest) ProtoMessage() {}

func (x *AddIntersectionIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddIntersectionIDRequest.ProtoReflect.Descriptor instead.
func (*AddIntersectionIDRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{7}
}

func (x *AddIntersectionIDRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *AddIntersectionIDRequest) GetIntersectionId() int32 {
	if x != nil {
		return x.IntersectionId
	}
	return 0
}

type RemoveIntersectionIDRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	UserId         string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	IntersectionId []int32                `protobuf:"varint,2,rep,packed,name=intersection_id,json=intersectionId,proto3" json:"intersection_id,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *RemoveIntersectionIDRequest) Reset() {
	*x = RemoveIntersectionIDRequest{}
	mi := &file_user_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveIntersectionIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveIntersectionIDRequest) ProtoMessage() {}

func (x *RemoveIntersectionIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveIntersectionIDRequest.ProtoReflect.Descriptor instead.
func (*RemoveIntersectionIDRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{8}
}

func (x *RemoveIntersectionIDRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *RemoveIntersectionIDRequest) GetIntersectionId() []int32 {
	if x != nil {
		return x.IntersectionId
	}
	return nil
}

// Password Management Messages
type ChangePasswordRequest struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	UserId          string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CurrentPassword string                 `protobuf:"bytes,2,opt,name=current_password,json=currentPassword,proto3" json:"current_password,omitempty"`
	NewPassword     string                 `protobuf:"bytes,3,opt,name=new_password,json=newPassword,proto3" json:"new_password,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ChangePasswordRequest) Reset() {
	*x = ChangePasswordRequest{}
	mi := &file_user_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangePasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangePasswordRequest) ProtoMessage() {}

func (x *ChangePasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangePasswordRequest.ProtoReflect.Descriptor instead.
func (*ChangePasswordRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{9}
}

func (x *ChangePasswordRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ChangePasswordRequest) GetCurrentPassword() string {
	if x != nil {
		return x.CurrentPassword
	}
	return ""
}

func (x *ChangePasswordRequest) GetNewPassword() string {
	if x != nil {
		return x.NewPassword
	}
	return ""
}

type ResetPasswordRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResetPasswordRequest) Reset() {
	*x = ResetPasswordRequest{}
	mi := &file_user_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResetPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetPasswordRequest) ProtoMessage() {}

func (x *ResetPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetPasswordRequest.ProtoReflect.Descriptor instead.
func (*ResetPasswordRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{10}
}

func (x *ResetPasswordRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

// Admin Management Messages
type AdminRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AdminUserId   string                 `protobuf:"bytes,2,opt,name=admin_user_id,json=adminUserId,proto3" json:"admin_user_id,omitempty"` // ID of the user making the request
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminRequest) Reset() {
	*x = AdminRequest{}
	mi := &file_user_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminRequest) ProtoMessage() {}

func (x *AdminRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminRequest.ProtoReflect.Descriptor instead.
func (*AdminRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{11}
}

func (x *AdminRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *AdminRequest) GetAdminUserId() string {
	if x != nil {
		return x.AdminUserId
	}
	return ""
}

// Response Messages
type UserResponse struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Id              string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name            string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email           string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	IsAdmin         bool                   `protobuf:"varint,4,opt,name=is_admin,json=isAdmin,proto3" json:"is_admin,omitempty"`
	IntersectionIds []int32                `protobuf:"varint,5,rep,packed,name=intersection_ids,json=intersectionIds,proto3" json:"intersection_ids,omitempty"`
	CreatedAt       *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt       *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *UserResponse) Reset() {
	*x = UserResponse{}
	mi := &file_user_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserResponse) ProtoMessage() {}

func (x *UserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserResponse.ProtoReflect.Descriptor instead.
func (*UserResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{12}
}

func (x *UserResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UserResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserResponse) GetIsAdmin() bool {
	if x != nil {
		return x.IsAdmin
	}
	return false
}

func (x *UserResponse) GetIntersectionIds() []int32 {
	if x != nil {
		return x.IntersectionIds
	}
	return nil
}

func (x *UserResponse) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *UserResponse) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type LoginUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ExpiresAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"` // Token expiration timestamp
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginUserResponse) Reset() {
	*x = LoginUserResponse{}
	mi := &file_user_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginUserResponse) ProtoMessage() {}

func (x *LoginUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginUserResponse.ProtoReflect.Descriptor instead.
func (*LoginUserResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{13}
}

func (x *LoginUserResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *LoginUserResponse) GetExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

var File_user_proto protoreflect.FileDescriptor

const file_user_proto_rawDesc = "" +
	"\n" +
	"\n" +
	"user.proto\x12\x11swiftsignals.user\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"(\n" +
	"\rUserIDRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\"[\n" +
	"\x13RegisterUserRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x14\n" +
	"\x05email\x18\x02 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x03 \x01(\tR\bpassword\"D\n" +
	"\x10LoginUserRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"-\n" +
	"\x15GetUserByEmailRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\"]\n" +
	"\x12GetAllUsersRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x1b\n" +
	"\tpage_size\x18\x02 \x01(\x05R\bpageSize\x12\x16\n" +
	"\x06filter\x18\x03 \x01(\tR\x06filter\"V\n" +
	"\x11UpdateUserRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x14\n" +
	"\x05email\x18\x03 \x01(\tR\x05email\"A\n" +
	"\x16IntersectionIDResponse\x12'\n" +
	"\x0fintersection_id\x18\x01 \x01(\x05R\x0eintersectionId\"\\\n" +
	"\x18AddIntersectionIDRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12'\n" +
	"\x0fintersection_id\x18\x02 \x01(\x05R\x0eintersectionId\"_\n" +
	"\x1bRemoveIntersectionIDRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12'\n" +
	"\x0fintersection_id\x18\x02 \x03(\x05R\x0eintersectionId\"~\n" +
	"\x15ChangePasswordRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12)\n" +
	"\x10current_password\x18\x02 \x01(\tR\x0fcurrentPassword\x12!\n" +
	"\fnew_password\x18\x03 \x01(\tR\vnewPassword\",\n" +
	"\x14ResetPasswordRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\"K\n" +
	"\fAdminRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\"\n" +
	"\radmin_user_id\x18\x02 \x01(\tR\vadminUserId\"\x84\x02\n" +
	"\fUserResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x14\n" +
	"\x05email\x18\x03 \x01(\tR\x05email\x12\x19\n" +
	"\bis_admin\x18\x04 \x01(\bR\aisAdmin\x12)\n" +
	"\x10intersection_ids\x18\x05 \x03(\x05R\x0fintersectionIds\x129\n" +
	"\n" +
	"created_at\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\"d\n" +
	"\x11LoginUserResponse\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x129\n" +
	"\n" +
	"expires_at\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\texpiresAt2\x83\n" +
	"\n" +
	"\vUserService\x12W\n" +
	"\fRegisterUser\x12&.swiftsignals.user.RegisterUserRequest\x1a\x1f.swiftsignals.user.UserResponse\x12V\n" +
	"\tLoginUser\x12#.swiftsignals.user.LoginUserRequest\x1a$.swiftsignals.user.LoginUserResponse\x12F\n" +
	"\n" +
	"LogoutUser\x12 .swiftsignals.user.UserIDRequest\x1a\x16.google.protobuf.Empty\x12P\n" +
	"\vGetUserByID\x12 .swiftsignals.user.UserIDRequest\x1a\x1f.swiftsignals.user.UserResponse\x12[\n" +
	"\x0eGetUserByEmail\x12(.swiftsignals.user.GetUserByEmailRequest\x1a\x1f.swiftsignals.user.UserResponse\x12W\n" +
	"\vGetAllUsers\x12%.swiftsignals.user.GetAllUsersRequest\x1a\x1f.swiftsignals.user.UserResponse0\x01\x12S\n" +
	"\n" +
	"UpdateUser\x12$.swiftsignals.user.UpdateUserRequest\x1a\x1f.swiftsignals.user.UserResponse\x12F\n" +
	"\n" +
	"DeleteUser\x12 .swiftsignals.user.UserIDRequest\x1a\x16.google.protobuf.Empty\x12g\n" +
	"\x16GetUserIntersectionIDs\x12 .swiftsignals.user.UserIDRequest\x1a).swiftsignals.user.IntersectionIDResponse0\x01\x12X\n" +
	"\x11AddIntersectionID\x12+.swiftsignals.user.AddIntersectionIDRequest\x1a\x16.google.protobuf.Empty\x12_\n" +
	"\x15RemoveIntersectionIDs\x12..swiftsignals.user.RemoveIntersectionIDRequest\x1a\x16.google.protobuf.Empty\x12R\n" +
	"\x0eChangePassword\x12(.swiftsignals.user.ChangePasswordRequest\x1a\x16.google.protobuf.Empty\x12P\n" +
	"\rResetPassword\x12'.swiftsignals.user.ResetPasswordRequest\x1a\x16.google.protobuf.Empty\x12D\n" +
	"\tMakeAdmin\x12\x1f.swiftsignals.user.AdminRequest\x1a\x16.google.protobuf.Empty\x12F\n" +
	"\vRemoveAdmin\x12\x1f.swiftsignals.user.AdminRequest\x1a\x16.google.protobuf.EmptyB\x11Z\x0fprotos/gen/userb\x06proto3"

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData []byte
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_user_proto_rawDesc), len(file_user_proto_rawDesc)))
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_user_proto_goTypes = []any{
	(*UserIDRequest)(nil),               // 0: swiftsignals.user.UserIDRequest
	(*RegisterUserRequest)(nil),         // 1: swiftsignals.user.RegisterUserRequest
	(*LoginUserRequest)(nil),            // 2: swiftsignals.user.LoginUserRequest
	(*GetUserByEmailRequest)(nil),       // 3: swiftsignals.user.GetUserByEmailRequest
	(*GetAllUsersRequest)(nil),          // 4: swiftsignals.user.GetAllUsersRequest
	(*UpdateUserRequest)(nil),           // 5: swiftsignals.user.UpdateUserRequest
	(*IntersectionIDResponse)(nil),      // 6: swiftsignals.user.IntersectionIDResponse
	(*AddIntersectionIDRequest)(nil),    // 7: swiftsignals.user.AddIntersectionIDRequest
	(*RemoveIntersectionIDRequest)(nil), // 8: swiftsignals.user.RemoveIntersectionIDRequest
	(*ChangePasswordRequest)(nil),       // 9: swiftsignals.user.ChangePasswordRequest
	(*ResetPasswordRequest)(nil),        // 10: swiftsignals.user.ResetPasswordRequest
	(*AdminRequest)(nil),                // 11: swiftsignals.user.AdminRequest
	(*UserResponse)(nil),                // 12: swiftsignals.user.UserResponse
	(*LoginUserResponse)(nil),           // 13: swiftsignals.user.LoginUserResponse
	(*timestamppb.Timestamp)(nil),       // 14: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),               // 15: google.protobuf.Empty
}
var file_user_proto_depIdxs = []int32{
	14, // 0: swiftsignals.user.UserResponse.created_at:type_name -> google.protobuf.Timestamp
	14, // 1: swiftsignals.user.UserResponse.updated_at:type_name -> google.protobuf.Timestamp
	14, // 2: swiftsignals.user.LoginUserResponse.expires_at:type_name -> google.protobuf.Timestamp
	1,  // 3: swiftsignals.user.UserService.RegisterUser:input_type -> swiftsignals.user.RegisterUserRequest
	2,  // 4: swiftsignals.user.UserService.LoginUser:input_type -> swiftsignals.user.LoginUserRequest
	0,  // 5: swiftsignals.user.UserService.LogoutUser:input_type -> swiftsignals.user.UserIDRequest
	0,  // 6: swiftsignals.user.UserService.GetUserByID:input_type -> swiftsignals.user.UserIDRequest
	3,  // 7: swiftsignals.user.UserService.GetUserByEmail:input_type -> swiftsignals.user.GetUserByEmailRequest
	4,  // 8: swiftsignals.user.UserService.GetAllUsers:input_type -> swiftsignals.user.GetAllUsersRequest
	5,  // 9: swiftsignals.user.UserService.UpdateUser:input_type -> swiftsignals.user.UpdateUserRequest
	0,  // 10: swiftsignals.user.UserService.DeleteUser:input_type -> swiftsignals.user.UserIDRequest
	0,  // 11: swiftsignals.user.UserService.GetUserIntersectionIDs:input_type -> swiftsignals.user.UserIDRequest
	7,  // 12: swiftsignals.user.UserService.AddIntersectionID:input_type -> swiftsignals.user.AddIntersectionIDRequest
	8,  // 13: swiftsignals.user.UserService.RemoveIntersectionIDs:input_type -> swiftsignals.user.RemoveIntersectionIDRequest
	9,  // 14: swiftsignals.user.UserService.ChangePassword:input_type -> swiftsignals.user.ChangePasswordRequest
	10, // 15: swiftsignals.user.UserService.ResetPassword:input_type -> swiftsignals.user.ResetPasswordRequest
	11, // 16: swiftsignals.user.UserService.MakeAdmin:input_type -> swiftsignals.user.AdminRequest
	11, // 17: swiftsignals.user.UserService.RemoveAdmin:input_type -> swiftsignals.user.AdminRequest
	12, // 18: swiftsignals.user.UserService.RegisterUser:output_type -> swiftsignals.user.UserResponse
	13, // 19: swiftsignals.user.UserService.LoginUser:output_type -> swiftsignals.user.LoginUserResponse
	15, // 20: swiftsignals.user.UserService.LogoutUser:output_type -> google.protobuf.Empty
	12, // 21: swiftsignals.user.UserService.GetUserByID:output_type -> swiftsignals.user.UserResponse
	12, // 22: swiftsignals.user.UserService.GetUserByEmail:output_type -> swiftsignals.user.UserResponse
	12, // 23: swiftsignals.user.UserService.GetAllUsers:output_type -> swiftsignals.user.UserResponse
	12, // 24: swiftsignals.user.UserService.UpdateUser:output_type -> swiftsignals.user.UserResponse
	15, // 25: swiftsignals.user.UserService.DeleteUser:output_type -> google.protobuf.Empty
	6,  // 26: swiftsignals.user.UserService.GetUserIntersectionIDs:output_type -> swiftsignals.user.IntersectionIDResponse
	15, // 27: swiftsignals.user.UserService.AddIntersectionID:output_type -> google.protobuf.Empty
	15, // 28: swiftsignals.user.UserService.RemoveIntersectionIDs:output_type -> google.protobuf.Empty
	15, // 29: swiftsignals.user.UserService.ChangePassword:output_type -> google.protobuf.Empty
	15, // 30: swiftsignals.user.UserService.ResetPassword:output_type -> google.protobuf.Empty
	15, // 31: swiftsignals.user.UserService.MakeAdmin:output_type -> google.protobuf.Empty
	15, // 32: swiftsignals.user.UserService.RemoveAdmin:output_type -> google.protobuf.Empty
	18, // [18:33] is the sub-list for method output_type
	3,  // [3:18] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_user_proto_rawDesc), len(file_user_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}
