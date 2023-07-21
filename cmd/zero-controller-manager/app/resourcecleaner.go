// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package app

import (
	"context"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/superproj/zero/internal/gateway/store"
)

type Cleaner interface {
	Sync(ctx context.Context) error
	Initialize(client client.Client, storeClient store.IStore)
}

type resourceCleaner struct {
	cleaners []Cleaner
}

// newResourceCleaner return a cleaner set used to clean deleted resources from zero db.
func newResourceCleaner(client client.Client, dbcli store.IStore, cleaners ...Cleaner) *resourceCleaner {
	for _, cleaner := range cleaners {
		cleaner.Initialize(client, dbcli)
	}

	return &resourceCleaner{cleaners}
}

func (c *resourceCleaner) Start(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Hour)
	for {
		c.SyncAll(ctx)
		<-ticker.C
	}
}

func (c *resourceCleaner) SyncAll(ctx context.Context) {
	for _, cleaner := range c.cleaners {
		//nolint: errcheck
		go cleaner.Sync(ctx)
	}
}
