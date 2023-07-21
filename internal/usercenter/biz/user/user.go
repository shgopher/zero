// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package user

//go:generate mockgen -self_package github.com/superproj/zero/internal/usercenter/biz/user -destination mock_user.go -package user github.com/superproj/zero/internal/usercenter/biz/user UserBiz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	validationutil "github.com/superproj/zero/internal/pkg/util/validation"
	"github.com/superproj/zero/internal/pkg/zerox"
	"github.com/superproj/zero/internal/usercenter/model"
	"github.com/superproj/zero/internal/usercenter/store"
	v1 "github.com/superproj/zero/pkg/api/usercenter/v1"
	"github.com/superproj/zero/pkg/authn"
)

// UserBiz defines methods used to handle user request.
type UserBiz interface {
	Create(ctx context.Context, req *v1.CreateUserRequest) error
	List(ctx context.Context, req *v1.ListUserRequest) (*v1.ListUserResponse, error)
	Get(ctx context.Context, req *v1.GetUserRequest) (*v1.UserReply, error)
	Update(ctx context.Context, req *v1.UpdateUserRequest) error
	Delete(ctx context.Context, req *v1.DeleteUserRequest) error

	// extensions apis
	UpdatePassword(ctx context.Context, req *v1.UpdatePasswordRequest) error
}

// userBiz struct implements the UserBiz interface and contains a store.IStore instance.
type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

// New returns a new instance of userBiz.
func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

// Create creates a new user and stores it in the database.
func (b *userBiz) Create(ctx context.Context, req *v1.CreateUserRequest) error {
	err := b.ds.TX(ctx, func(ctx context.Context) error {
		var userM model.UserM
		_ = copier.Copy(&userM, req)

		if err := b.ds.Users().Create(ctx, &userM); err != nil {
			return v1.ErrorUserCreateFailed("create user failed: %s", err.Error())
		}

		secretM := &model.SecretM{
			UserID:      userM.UserID,
			Name:        "generated",
			Expires:     0,
			Description: "automatically generated when user is created",
		}
		if err := b.ds.Secrets().Create(ctx, secretM); err != nil {
			return v1.ErrorSecretCreateFailed("create secret failed: %s", err.Error())
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// List retrieves a list of all users from the database.
func (b *userBiz) List(ctx context.Context, req *v1.ListUserRequest) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*v1.UserReply, 0)
	for _, item := range list {
		var u v1.UserReply
		_ = copier.Copy(&u, &item)
		u.CreatedAt = timestamppb.New(item.CreatedAt)
		u.UpdatedAt = timestamppb.New(item.UpdatedAt)
		u.Password = "******"
		users = append(users, &u)
	}

	return &v1.ListUserResponse{TotalCount: count, Users: users}, nil
}

// Get retrieves a single user from the database.
func (b *userBiz) Get(ctx context.Context, req *v1.GetUserRequest) (*v1.UserReply, error) {
	filters := map[string]interface{}{"username": req.Username}
	if !validationutil.IsAdminUser(zerox.FromUserID(ctx)) {
		filters["userID"] = zerox.FromUserID(ctx)
	}

	userM, err := b.ds.Users().Fetch(ctx, filters)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrorUserNotFound(err.Error())
		}

		return nil, err
	}

	var user v1.UserReply
	_ = copier.Copy(&user, userM)
	user.Password = "******"
	user.CreatedAt = timestamppb.New(userM.CreatedAt)
	user.UpdatedAt = timestamppb.New(userM.UpdatedAt)
	return &user, nil
}

// Update updates a user's information in the database.
func (b *userBiz) Update(ctx context.Context, req *v1.UpdateUserRequest) error {
	filters := map[string]interface{}{"username": req.Username}
	if !validationutil.IsAdminUser(zerox.FromUserID(ctx)) {
		filters["userID"] = zerox.FromUserID(ctx)
	}

	userM, err := b.ds.Users().Fetch(ctx, filters)
	if err != nil {
		return err
	}

	if req.Nickname != nil {
		userM.Nickname = *req.Nickname
	}
	if req.Email != nil {
		userM.Email = *req.Email
	}
	if req.Phone != nil {
		userM.Phone = *req.Phone
	}

	return b.ds.Users().Update(ctx, userM)
}

// UpdatePassword updates a user's password in the database.
// Note that after updating the password, if the JWT Token has not expired, it can
// still be accessed through the token, the token is not deleted synchronously here.
func (b *userBiz) UpdatePassword(ctx context.Context, req *v1.UpdatePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, zerox.FromUserID(ctx), req.Username)
	if err != nil {
		return err
	}

	if err := authn.Compare(userM.Password, req.OldPassword); err != nil {
		return v1.ErrorUserLoginFailed("password incorrect")
	}
	userM.Password, _ = authn.Encrypt(req.NewPassword)

	return b.ds.Users().Update(ctx, userM)
}

// Delete deletes a user from the database.
func (b *userBiz) Delete(ctx context.Context, req *v1.DeleteUserRequest) error {
	filters := map[string]interface{}{"username": req.Username}
	if !validationutil.IsAdminUser(zerox.FromUserID(ctx)) {
		filters["userID"] = zerox.FromUserID(ctx)
	}

	return b.ds.Users().Delete(ctx, filters)
}
