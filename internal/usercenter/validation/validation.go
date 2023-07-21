// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package validation

import (
	"context"
	"fmt"

	"github.com/google/wire"

	"github.com/superproj/zero/internal/pkg/known"
	ucknown "github.com/superproj/zero/internal/pkg/known/usercenter"
	"github.com/superproj/zero/internal/pkg/validation"
	"github.com/superproj/zero/internal/pkg/zerox"
	"github.com/superproj/zero/internal/usercenter/locales"
	"github.com/superproj/zero/internal/usercenter/store"
	v1 "github.com/superproj/zero/pkg/api/usercenter/v1"
	"github.com/superproj/zero/pkg/i18n"
)

// ProviderSet is validator providers.
var ProviderSet = wire.NewSet(New)

// validator struct implements the validation.CustomValidator interface.
type validator struct {
	ds store.IStore
}

// New creates and initializes a custom validator.
// It receives an instance of store.IStore interface as parameter ds
// and returns a new validation.CustomValidator and an error.
func New(ds store.IStore) (validation.CustomValidator, error) {
	vd := &validator{ds: ds}

	return vd, nil
}

// ValidateCreateUserRequest validates the request to create a user.
// If the validation fails, it returns an error; otherwise, it returns nil.
func (vd *validator) ValidateCreateUserRequest(ctx context.Context, req *v1.CreateUserRequest) error {
	if _, err := vd.ds.Users().GetByUsername(ctx, req.Username); err == nil {
		return i18n.FromContext(ctx).E(locales.UserAlreadyExists)
	}

	return nil
}

// ValidateListUserRequest validates the request to list users.
// Ensures that only a user with the AdminUserID can view the list of users, otherwise returning an error.
func (vd *validator) ValidateListUserRequest(ctx context.Context, req *v1.ListUserRequest) error {
	if userID := zerox.FromUserID(ctx); userID != known.AdminUserID {
		return i18n.FromContext(ctx).E(locales.UserListUnauthorized)
	}

	return nil
}

// ValidateCreateSecretRequest validates the request to create a secret.
// Returns an error if the maximum number of secrets is reached.
func (vd *validator) ValidateCreateSecretRequest(ctx context.Context, req *v1.CreateSecretRequest) error {
	_, secrets, err := vd.ds.Secrets().List(ctx, zerox.FromUserID(ctx))
	if err != nil {
		return err
	}

	if len(secrets) >= ucknown.MaxSecretCount {
		return fmt.Errorf("secret reach the max count %d", ucknown.MaxSecretCount)
	}

	return nil
}

// ValidateAuthRequest validates the authentication request.
// In this sample, no actual validation is needed, so it returns nil directly.
func (vd *validator) ValidateAuthRequest(ctx context.Context, req *v1.AuthRequest) error {
	return nil
}

// ValidateAuthorizeRequest validates the authorization request.
// In this sample, no actual validation is needed, so it returns nil directly.
func (vd *validator) ValidateAuthorizeRequest(ctx context.Context, req *v1.AuthorizeRequest) error {
	return nil
}
