// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//go:build wireinject
// +build wireinject

package app

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/google/wire"

	"github.com/superproj/zero/internal/gateway/store"
	"github.com/superproj/zero/pkg/db"
)

func wireStoreClient(*db.MySQLOptions) (store.IStore, error) {
	wire.Build(
		db.ProviderSet,
		store.ProviderSet,
	)

	return nil, nil
}
