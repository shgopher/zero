// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package gateway

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"

	"github.com/superproj/zero/internal/gateway/biz"
	"github.com/superproj/zero/internal/gateway/server"
	"github.com/superproj/zero/internal/gateway/service"
	"github.com/superproj/zero/internal/gateway/store"
	customvalidation "github.com/superproj/zero/internal/gateway/validation"
	"github.com/superproj/zero/internal/pkg/bootstrap"
	"github.com/superproj/zero/internal/pkg/client/usercenter"
	"github.com/superproj/zero/internal/pkg/idempotent"
	"github.com/superproj/zero/internal/pkg/validation"
	"github.com/superproj/zero/pkg/db"
	clientset "github.com/superproj/zero/pkg/generated/clientset/versioned"
	genericoptions "github.com/superproj/zero/pkg/options"
)

// wireApp init kratos application.
func wireApp(
	<-chan struct{},
	bootstrap.AppInfo,
	*server.Config,
	clientset.Interface,
	*db.MySQLOptions,
	*db.RedisOptions,
	*usercenter.UserCenterOptions,
	*genericoptions.RedisOptions,
	*genericoptions.EtcdOptions,
) (*kratos.App, func(), error) {
	wire.Build(
		bootstrap.ProviderSet,
		bootstrap.NewEtcdRegistrar,
		server.ProviderSet,
		store.ProviderSet,
		usercenter.ProviderSet,
		db.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		validation.ProviderSet,
		idempotent.ProviderSet,
		customvalidation.ProviderSet,
		createInformers,
	)

	return nil, nil, nil
}
