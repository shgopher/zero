// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package watcher provides functions used by all watchers.
package watcher

import (
	"github.com/superproj/zero/internal/pkg/client/store"
	clientset "github.com/superproj/zero/pkg/generated/clientset/versioned"
)

// Config contains all the options need the watchers.
type Config struct {
	// The purpose of nightwatch is to handle asynchronous tasks on the zero platform
	// in a unified manner, so a store aggregation type is needed here.
	Store  store.Interface
	Client clientset.Interface
}
