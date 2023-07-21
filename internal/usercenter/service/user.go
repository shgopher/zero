// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package service

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/superproj/zero/pkg/api/usercenter/v1"
)

// CreateUser receives a CreateUserRequest and creates a new user record in the datastore.
// It returns an empty response (emptypb.Empty) and an error if there's any.
func (s *UserCenterService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.biz.Users().Create(ctx, req)
}

// ListUser receives a ListUserRequest and returns a ListUserResponse containing a list
// of users with pagination and an error if there is any.
func (s *UserCenterService) ListUser(ctx context.Context, req *v1.ListUserRequest) (*v1.ListUserResponse, error) {
	return s.biz.Users().List(ctx, req)
}

// GetUser receives a GetUserRequest and returns a UserReply with the corresponding user information
// and an error if there's any.
func (s *UserCenterService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.UserReply, error) {
	return s.biz.Users().Get(ctx, req)
}

// UpdateUser receives an UpdateUserRequest and updates the user record in the datastore.
// It returns an empty response (emptypb.Empty) and an error if there's any.
func (s *UserCenterService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.biz.Users().Update(ctx, req)
}

// UpdatePassword receives an UpdatePasswordRequest and updates the user's password in the datastore.
// It returns an empty response (emptypb.Empty) and an error if there's any.
func (s *UserCenterService) UpdatePassword(ctx context.Context, req *v1.UpdatePasswordRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.biz.Users().UpdatePassword(ctx, req)
}

// DeleteUser receives a DeleteUserRequest and removes the user record from the datastore.
// It returns an empty response (emptypb.Empty) and an error if there's any.
func (s *UserCenterService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.biz.Users().Delete(ctx, req)
}
