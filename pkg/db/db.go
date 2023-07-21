// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

// ProviderSet is db providers.
var ProviderSet = wire.NewSet(NewMySQL, NewRedis, wire.Bind(new(redis.UniversalClient), new(*redis.Client)))
