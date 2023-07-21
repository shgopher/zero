// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package db

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisOptions defines optsions for mysql database.
type RedisOptions struct {
	// host:port address.
	Addr       string
	Username   string
	Password   string
	Database   int
	MaxRetries int
	Timeout    time.Duration
}

// NewRedis create a new gorm db instance with the given options.
func NewRedis(opts *RedisOptions) (*redis.Client, error) {
	options := &redis.Options{
		Addr:        opts.Addr,
		Username:    opts.Username,
		Password:    opts.Password,
		DB:          opts.Database,
		MaxRetries:  opts.MaxRetries,
		DialTimeout: opts.Timeout,
	}

	rdb := redis.NewClient(options)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}
