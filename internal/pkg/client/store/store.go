// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package store

import (
	"sync"

	"github.com/google/wire"

	gwstore "github.com/superproj/zero/internal/gateway/store"
	ucstore "github.com/superproj/zero/internal/usercenter/store"
)

// ProviderSet is store providers.
var ProviderSet = wire.NewSet(NewStore, wire.Bind(new(Interface), new(*datastore)))

var (
	once sync.Once
	S    *datastore
)

// Interface defines the storage interface.
type Interface interface {
	Gateway() gwstore.IStore
	UserCenter() ucstore.IStore
}

type datastore struct {
	gw gwstore.IStore
	uc ucstore.IStore
}

var _ Interface = (*datastore)(nil)

func (ds *datastore) Gateway() gwstore.IStore {
	return ds.gw
}

func (ds *datastore) UserCenter() ucstore.IStore {
	return ds.uc
}

func NewStore(gw gwstore.IStore, uc ucstore.IStore) *datastore {
	once.Do(func() {
		S = &datastore{gw: gw, uc: uc}
	})

	return S
}
