// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package auth

//go:generate mockgen -self_package github.com/superproj/zero/internal/usercenter/biz/auth -destination mock_auth.go -package auth github.com/superproj/zero/internal/usercenter/biz/auth AuthBiz

import (
	"context"

	"github.com/superproj/zero/internal/pkg/zerox"
	"github.com/superproj/zero/internal/usercenter/auth"
	"github.com/superproj/zero/internal/usercenter/locales"
	"github.com/superproj/zero/internal/usercenter/store"
	v1 "github.com/superproj/zero/pkg/api/usercenter/v1"
	"github.com/superproj/zero/pkg/authn"
	"github.com/superproj/zero/pkg/i18n"
	"github.com/superproj/zero/pkg/log"
)

// AuthBiz defines functions used for authentication and authorization.
type AuthBiz interface {
	// Login authenticates a user and returns a token.
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error)

	// Logout invalidates a token.
	Logout(ctx context.Context, req *v1.LogoutRequest) error

	// RefreshToken refreshes an existing token and returns a new one.
	RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.LoginReply, error)

	// Authenticate validates an access token and returns the associated user ID.
	Authenticate(ctx context.Context, accessToken string) (*v1.AuthenticateResponse, error)

	// Authorize checks if a user has the necessary permissions to perform an action on an object.
	Authorize(ctx context.Context, sub, obj, act string) (*v1.AuthorizeResponse, error)
}

// The authBiz struct contains dependencies required for authentication and authorization.
type authBiz struct {
	ds    store.IStore
	authn authn.Authenticator
	auth  auth.AuthProvider
}

var _ AuthBiz = (*authBiz)(nil)

// New creates a new authBiz instance.
func New(ds store.IStore, authn authn.Authenticator, auth auth.AuthProvider) *authBiz {
	return &authBiz{authn: authn, auth: auth, ds: ds}
}

// Login authenticates a user and returns a token.
func (b *authBiz) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	userM, err := b.ds.Users().GetByUsername(ctx, req.Username)
	if err != nil {
		log.C(ctx).Errorf(err.Error())
		return nil, i18n.FromContext(ctx).E(locales.RecordNotFound)
	}

	if err := authn.Compare(userM.Password, req.Password); err != nil {
		log.C(ctx).Errorf(err.Error())
		return nil, i18n.FromContext(ctx).E(locales.IncorrectPassword)
	}

	token, err := b.authn.Sign(ctx, userM.UserID)
	if err != nil {
		log.C(ctx).Errorf(err.Error())
		return nil, i18n.FromContext(ctx).E(locales.JwtTokenSignFail)
	}

	return &v1.LoginReply{
		Token:     token.GetAccessToken(),
		Type:      token.GetTokenType(),
		ExpiresAt: token.GetExpiresAt(),
	}, nil
}

// Logout invalidates a token.
func (b *authBiz) Logout(ctx context.Context, req *v1.LogoutRequest) error {
	if err := b.authn.Destroy(ctx, zerox.FromAccessToken(ctx)); err != nil {
		log.C(ctx).Errorf(err.Error())
		return err
	}

	return nil
}

// RefreshToken refreshes an existing token and returns a new one.
func (b *authBiz) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.LoginReply, error) {
	token, err := b.authn.Sign(ctx, zerox.FromUserID(ctx))
	if err != nil {
		log.C(ctx).Errorf(err.Error())
		return nil, err
	}

	return &v1.LoginReply{
		Token:     token.GetAccessToken(),
		Type:      token.GetTokenType(),
		ExpiresAt: token.GetExpiresAt(),
	}, nil
}

// Authenticate validates an access token and returns the associated user ID.
func (b *authBiz) Authenticate(ctx context.Context, accessToken string) (*v1.AuthenticateResponse, error) {
	userID, err := b.auth.Verify(accessToken)
	if err != nil {
		log.C(ctx).Errorf(err.Error())
		return nil, err
	}

	return &v1.AuthenticateResponse{UserID: userID}, nil
}

// Authorize checks if a user has the necessary permissions to perform an action on an object.
func (b *authBiz) Authorize(ctx context.Context, sub, obj, act string) (*v1.AuthorizeResponse, error) {
	allowed, err := b.auth.Authorize(sub, obj, act)
	if err != nil {
		log.C(ctx).Errorf(err.Error())
		return nil, err
	}
	return &v1.AuthorizeResponse{Allowed: allowed}, nil
}
