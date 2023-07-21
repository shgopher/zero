// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/superproj/zero/internal/gateway/store"
	"github.com/superproj/zero/pkg/db"
)

// Injectors from wire.go:

func wireStoreClient(mySQLOptions *db.MySQLOptions) (store.IStore, error) {
	gormDB, err := db.NewMySQL(mySQLOptions)
	if err != nil {
		return nil, err
	}
	datastore := store.NewStore(gormDB)
	return datastore, nil
}