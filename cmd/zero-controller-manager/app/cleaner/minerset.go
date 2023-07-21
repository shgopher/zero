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

type MinerSetCleaner struct {
	mu     sync.Mutex
	client client.Client
	dbcli  store.IStore
}

func (c *MinerSetCleaner) Initialize(client client.Client, dbcli store.IStore) {
	c.client = client
	c.dbcli = dbcli
}

func (c *MinerSetCleaner) Sync(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	klog.V(4).InfoS("Cleanup minersets from minerset table")
	_, minersets, err := c.dbcli.MinerSets().List(ctx, "")
	if err != nil {
		klog.ErrorS(err, "Failed to list minersets")
		return err
	}

	klog.V(4).InfoS("Successfully got minersets", "count", len(minersets))
	for _, minerset := range minersets {
		ms := v1beta1.MinerSet{}
		key := client.ObjectKey{Namespace: minerset.Namespace, Name: minerset.Name}
		if err := c.client.Get(ctx, key, &ms); err != nil {
			if apierrors.IsNotFound(err) {
				filter := map[string]interface{}{"namespace": minerset.Namespace, "name": minerset.Name}
				if derr := c.dbcli.MinerSets().Delete(ctx, filter); derr != nil {
					klog.V(1).InfoS("Failed to delete minerset", "minerset", klog.KRef(minerset.Namespace, minerset.Name), "err", derr)
					continue
				}
				klog.V(4).InfoS("Successfully delete minerset", "minerset", klog.KRef(minerset.Namespace, minerset.Name))
			}

			klog.ErrorS(err, "Failed to get minerset", "minerset", klog.KRef(key.Namespace, key.Name))
			return err
		}
	}

	return nil
}
