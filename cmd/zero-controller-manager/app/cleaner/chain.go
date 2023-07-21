// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package cleaner

import (
	"context"
	"sync"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/superproj/zero/internal/gateway/store"
	"github.com/superproj/zero/pkg/apis/apps/v1beta1"
)

type ChainCleaner struct {
	mu     sync.Mutex
	client client.Client
	dbcli  store.IStore
}

func (c *ChainCleaner) Initialize(client client.Client, dbcli store.IStore) {
	c.client = client
	c.dbcli = dbcli
}

func (c *ChainCleaner) Sync(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	klog.V(4).InfoS("Cleanup chains from chain table")
	_, chains, err := c.dbcli.Chains().List(ctx, "")
	if err != nil {
		klog.ErrorS(err, "Failed to list chains")
		return err
	}

	klog.V(4).InfoS("Successfully got chains", "count", len(chains))
	for _, chain := range chains {
		ch := v1beta1.Chain{}
		key := client.ObjectKey{Namespace: chain.Namespace, Name: chain.Name}
		if err := c.client.Get(ctx, key, &ch); err != nil {
			if apierrors.IsNotFound(err) {
				filter := map[string]interface{}{"namespace": chain.Namespace, "name": chain.Name}
				if derr := c.dbcli.Chains().Delete(ctx, filter); derr != nil {
					klog.V(1).InfoS("Failed to delete chain", "chain", klog.KRef(chain.Namespace, chain.Name), "err", derr)
					continue
				}
				klog.V(4).InfoS("Successfully delete chain", "chain", klog.KRef(chain.Namespace, chain.Name))
			}

			klog.ErrorS(err, "Failed to get chain", "chain", klog.KRef(key.Namespace, key.Name))
			return err
		}
	}

	return nil
}
