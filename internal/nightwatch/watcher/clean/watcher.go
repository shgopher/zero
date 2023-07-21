// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package clean is a watcher implement.
package clean

import (
	"context"

	"github.com/superproj/zero/internal/nightwatch/watcher"
	"github.com/superproj/zero/internal/pkg/client/store"
	"github.com/superproj/zero/pkg/log"
)

var _ watcher.Watcher = (*cleanWatcher)(nil)

// watcher implement.
type cleanWatcher struct {
	store store.Interface
}

// Run runs the watcher.
func (w *cleanWatcher) Run() {
	_, miners, err := w.store.Gateway().Miners().List(context.Background(), "")
	if err != nil {
		log.Errorw(err, "Failed to list miners")
	}

	for _, m := range miners {
		log.Infow("Get a miner", "miner", m.Name)
	}
}

// Init initializes the watcher for later execution.
func (w *cleanWatcher) Init(ctx context.Context, config *watcher.Config) error {
	w.store = config.Store
	return nil
}

func init() {
	watcher.Register(&cleanWatcher{})
}
