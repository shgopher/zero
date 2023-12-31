// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	addr = "127.0.0.1:6379"
)

func TestStore(t *testing.T) {
	store := NewStore(&Config{
		Addr:      addr,
		Database:  1,
		KeyPrefix: "prefix",
	})

	defer store.Close()

	key := "test"
	ctx := context.Background()
	err := store.Set(ctx, key, 0)
	assert.Nil(t, err)

	b, err := store.Check(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	b, err = store.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)
}
