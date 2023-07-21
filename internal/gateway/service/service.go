// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package service

import (
	"github.com/google/wire"

	"github.com/superproj/zero/internal/gateway/biz"
	"github.com/superproj/zero/internal/pkg/idempotent"
	v1 "github.com/superproj/zero/pkg/api/gateway/v1"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGatewayService)

type GatewayService struct {
	v1.UnimplementedGatewayServer

	biz biz.BizFactory
	idt *idempotent.Idempotent
}

// func NewGatewayService(biz biz.BizFactory, lister appsv1beta1.Interface, client clientset.Interface) *GatewayService {.
func NewGatewayService(biz biz.BizFactory, idt *idempotent.Idempotent) *GatewayService {
	return &GatewayService{biz: biz, idt: idt}
}
