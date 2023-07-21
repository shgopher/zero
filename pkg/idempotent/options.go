// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package idempotent

import (
	"github.com/go-redis/redis/v8"
)

type Options struct {
	redis  redis.UniversalClient
	prefix string
	expire int
}

func WithRedis(rd redis.UniversalClient) func(*Options) {
	return func(options *Options) {
		if rd != nil {
			getOptionsOrSetDefault(options).redis = rd
		}
	}
}

func WithPrefix(prefix string) func(*Options) {
	return func(options *Options) {
		if prefix != "" {
			getOptionsOrSetDefault(options).prefix = prefix
		}
	}
}

func WithExpire(min int) func(*Options) {
	return func(options *Options) {
		if min > 0 {
			getOptionsOrSetDefault(options).expire = min
		}
	}
}

func getOptionsOrSetDefault(options *Options) *Options {
	if options == nil {
		return &Options{
			prefix: "idempotent",
			expire: 60,
		}
	}
	return options
}
