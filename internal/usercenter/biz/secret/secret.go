// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package secret

//go:generate mockgen -self_package github.com/superproj/zero/internal/usercenter/biz/secret -destination mock_secret.go -package secret github.com/superproj/zero/internal/usercenter/biz/secret SecretBiz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/superproj/zero/internal/pkg/zerox"
	"github.com/superproj/zero/internal/usercenter/model"
	"github.com/superproj/zero/internal/usercenter/store"
	v1 "github.com/superproj/zero/pkg/api/usercenter/v1"
)

// SecretBiz defines functions used to handle secret request.
type SecretBiz interface {
	Create(ctx context.Context, req *v1.CreateSecretRequest) error
	List(ctx context.Context, req *v1.ListSecretRequest) (*v1.ListSecretResponse, error)
	Get(ctx context.Context, req *v1.GetSecretRequest) (*v1.SecretReply, error)
	Update(ctx context.Context, req *v1.UpdateSecretRequest) error
	Delete(ctx context.Context, req *v1.DeleteSecretRequest) error
}

type secretBiz struct {
	ds store.IStore
}

var _ SecretBiz = (*secretBiz)(nil)

// New creates a new instance of the secretBiz struct.
func New(ds store.IStore) *secretBiz {
	return &secretBiz{ds: ds}
}

// Create creates a new secret.
func (b *secretBiz) Create(ctx context.Context, req *v1.CreateSecretRequest) error {
	var secretM model.SecretM
	_ = copier.Copy(&secretM, req)
	secretM.UserID = zerox.FromUserID(ctx)

	if err := b.ds.Secrets().Create(ctx, &secretM); err != nil {
		return v1.ErrorSecretCreateFailed("create secret failed: %s", err.Error())
	}

	return nil
}

// List returns a list of secrets.
func (b *secretBiz) List(ctx context.Context, req *v1.ListSecretRequest) (*v1.ListSecretResponse, error) {
	count, list, err := b.ds.Secrets().List(ctx, zerox.FromUserID(ctx))
	if err != nil {
		return nil, err
	}

	secrets := make([]*v1.SecretReply, 0)
	for _, item := range list {
		var s v1.SecretReply
		_ = copier.Copy(&s, &item)
		s.CreatedAt = timestamppb.New(item.CreatedAt)
		s.UpdatedAt = timestamppb.New(item.UpdatedAt)
		secrets = append(secrets, &s)
	}

	return &v1.ListSecretResponse{TotalCount: count, Secrets: secrets}, nil
}

// Get returns a single secret.
func (b *secretBiz) Get(ctx context.Context, req *v1.GetSecretRequest) (*v1.SecretReply, error) {
	secretM, err := b.ds.Secrets().Get(ctx, zerox.FromUserID(ctx), req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrorSecretNotFound(err.Error())
		}

		return nil, err
	}

	var secret v1.SecretReply
	_ = copier.Copy(&secret, secretM)
	secret.CreatedAt = timestamppb.New(secretM.CreatedAt)
	secret.UpdatedAt = timestamppb.New(secretM.UpdatedAt)

	return &secret, nil
}

// Update updates a secret.
func (b *secretBiz) Update(ctx context.Context, req *v1.UpdateSecretRequest) error {
	secret, err := b.ds.Secrets().Get(ctx, zerox.FromUserID(ctx), req.Name)
	if err != nil {
		return err
	}

	if req.Expires != nil {
		secret.Expires = *req.Expires
	}
	if req.Status != nil {
		secret.Status = *req.Status
	}
	if req.Description != nil {
		secret.Description = *req.Description
	}

	return b.ds.Secrets().Update(ctx, secret)
}

// Delete deletes a secret.
func (b *secretBiz) Delete(ctx context.Context, req *v1.DeleteSecretRequest) error {
	return b.ds.Secrets().Delete(ctx, zerox.FromUserID(ctx), req.Name)
}
