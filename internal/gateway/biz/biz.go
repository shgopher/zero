// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package biz

//go:generate mockgen -self_package github.com/superproj/zero/internal/gateway/biz -destination mock_biz.go -package biz github.com/superproj/zero/internal/gateway/biz BizFactory

import (
	"github.com/google/wire"

	"github.com/superproj/zero/internal/gateway/biz/miner"
	"github.com/superproj/zero/internal/gateway/biz/minerset"
	"github.com/superproj/zero/internal/gateway/store"
	clientset "github.com/superproj/zero/pkg/generated/clientset/versioned"
	informers "github.com/superproj/zero/pkg/generated/informers/externalversions"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewBiz, wire.Bind(new(BizFactory), new(*biz)))

// BizFactory defines functions used to return resource interface.
type BizFactory interface {
	Miners() miner.MinerBiz
	MinerSets() minerset.MinerSetBiz
}

type biz struct {
	ds store.IStore
	cl clientset.Interface
	f  informers.SharedInformerFactory
}

// NewBiz returns BizFactory interface.
func NewBiz(ds store.IStore, cl clientset.Interface, f informers.SharedInformerFactory) *biz {
	return &biz{ds, cl, f}
}

func (b *biz) MinerSets() minerset.MinerSetBiz {
	return minerset.New(b.ds, b.cl, b.f)
}

func (b *biz) Miners() miner.MinerBiz {
	return miner.New(b.ds, b.cl, b.f)
}
