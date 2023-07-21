// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package minertype

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/superproj/zero/internal/pkg/config"
)

type MinerProfile struct {
	CPU              resource.Quantity `json:"cpu,omitempty"`
	Memory           resource.Quantity `json:"memory,omitempty"`
	MiningDifficulty int               `json:"miningDifficulty,omitempty"`
}

type cacher struct {
	mu     sync.Mutex
	client interface{}

	data *cache.Cache
}

var (
	g    = new(cacher)
	once = sync.Once{}
)

func setDefault(tc *cacher) {
	g = tc
}

func (tc *cacher) loadData() error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	cm, err := config.MinerTypesName.GetConfig(tc.client)
	if err != nil {
		return err
	}

	c := cache.New(24*time.Hour, 48*time.Hour)
	for k, v := range cm.Data {
		var profile MinerProfile
		if err := json.Unmarshal([]byte(v), &profile); err != nil {
			return err
		}

		c.Set(k, &profile, cache.DefaultExpiration)
	}

	tc.data = c
	return nil
}

func Init(ctx context.Context, client interface{}) error {
	var err error
	once.Do(func() {
		tc := cacher{
			client: client,
		}

		if err = tc.loadData(); err != nil {
			return
		}

		setDefault(&tc)
	})

	return err
}

func GetProfile(key string) (*MinerProfile, bool) {
	obj, found := g.data.Get(key)
	if found {
		return obj.(*MinerProfile), true
	}

	if err := g.loadData(); err != nil {
		return nil, false
	}

	obj, found = g.data.Get(key)
	if found {
		return obj.(*MinerProfile), true
	}

	return nil, false
}
