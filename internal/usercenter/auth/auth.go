// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package auth

//go:generate mockgen -self_package github.com/superproj/zero/internal/usercenter/auth -destination mock_auth.go -package auth github.com/superproj/zero/internal/usercenter/auth AuthProvider,AuthzInterface,AuthnInterface

import (
	"github.com/google/wire"
)

// ProviderSet is a Wire provider set that creates a new instance of auth.
var ProviderSet = wire.NewSet(NewAuth, wire.Bind(new(AuthProvider), new(*auth)), AuthnProviderSet, AuthzProviderSet)

// AuthProvider is an interface that combines both the AuthnInterface and AuthzInterface interfaces.
type AuthProvider interface {
	AuthnInterface
	AuthzInterface
}

// auth is a struct that implements AuthnInterface and AuthzInterface interfaces.
type auth struct {
	authn AuthnInterface
	authz AuthzInterface
}

// NewAuth is a constructor function that creates a new instance of auth struct.
func NewAuth(authn AuthnInterface, authz AuthzInterface) *auth {
	return &auth{authn: authn, authz: authz}
}

// Verify is a method that implements Verify method of AuthnInterface.
func (a *auth) Verify(accessToken string) (string, error) {
	return a.authn.Verify(accessToken)
}

// Authorize is a method that implements Authorize method of AuthzInterface.
func (a *auth) Authorize(rvals ...interface{}) (bool, error) {
	return a.authz.Authorize(rvals...)
}
