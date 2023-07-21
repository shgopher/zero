// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package service

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/errors"

	v1 "github.com/superproj/zero/pkg/api/usercenter/v1"
	"github.com/superproj/zero/pkg/api/zerrors"
)

// Login authenticates the user credentials and returns a token on success.
func (s *UserCenterService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	resp, err := s.biz.Auths().Login(ctx, req)
	if err != nil {
		return &v1.LoginReply{}, v1.ErrorUserLoginFailed(err.Error())
	}

	return resp, nil
}

// Logout invalidates the user token.
func (s *UserCenterService) Logout(ctx context.Context, req *v1.LogoutRequest) (*emptypb.Empty, error) {
	if err := s.biz.Auths().Logout(ctx, req); err != nil {
		return &emptypb.Empty{}, zerrors.ErrorUnknown(err.Error())
	}

	return &emptypb.Empty{}, nil
}

// RefreshToken generates a new token using the refresh token.
func (s *UserCenterService) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.LoginReply, error) {
	resp, err := s.biz.Auths().RefreshToken(ctx, req)
	if err != nil {
		return &v1.LoginReply{}, errors.New(401, "UserLoginFailed", err.Error())
	}

	return resp, nil
}

// Auth authenticates and authorizes the user token for an object/action.
func (s *UserCenterService) Auth(ctx context.Context, req *v1.AuthRequest) (*v1.AuthResponse, error) {
	authn, err := s.Authenticate(ctx, &v1.AuthenticateRequest{Token: req.Token})
	if err != nil {
		return &v1.AuthResponse{}, err
	}

	authz, err := s.Authorize(ctx, &v1.AuthorizeRequest{Sub: authn.UserID, Obj: req.Obj, Act: req.Act})
	if err != nil {
		return &v1.AuthResponse{}, err
	}

	return &v1.AuthResponse{UserID: authn.UserID, Allowed: authz.Allowed}, nil
}

// Authenticate validates the user token and returns the user ID.
func (s *UserCenterService) Authenticate(ctx context.Context, req *v1.AuthenticateRequest) (*v1.AuthenticateResponse, error) {
	resp, err := s.biz.Auths().Authenticate(ctx, req.Token)
	if err != nil {
		return &v1.AuthenticateResponse{}, err
	}

	return resp, nil
}

// Authorize checks whether the user is authorized for the object/action.
func (s *UserCenterService) Authorize(ctx context.Context, req *v1.AuthorizeRequest) (*v1.AuthorizeResponse, error) {
	allowed, err := s.biz.Auths().Authorize(ctx, req.Sub, req.Obj, req.Act)
	if err != nil {
		return &v1.AuthorizeResponse{}, err
	}

	return allowed, nil
}
