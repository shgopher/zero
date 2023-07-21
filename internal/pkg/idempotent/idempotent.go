// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package idempotent

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"

	"github.com/superproj/zero/pkg/idempotent"
	"github.com/superproj/zero/pkg/log"
)

var ProviderSet = wire.NewSet(NewIdempotent)

type Idempotent struct {
	idempotent *idempotent.Idempotent
}

func (idt Idempotent) Token(ctx context.Context) string {
	return idt.idempotent.Token(ctx)
}

func (idt Idempotent) Check(ctx context.Context, token string) bool {
	return idt.idempotent.Check(ctx, token)
}

// NewIdempotent is initialize idempotent from config.
func NewIdempotent(redis redis.UniversalClient) (idt *Idempotent, err error) {
	ins := idempotent.New(idempotent.WithRedis(redis))
	idt = &Idempotent{
		idempotent: ins,
	}

	log.Infow("initialize idempotent success")
	return idt, nil
}
