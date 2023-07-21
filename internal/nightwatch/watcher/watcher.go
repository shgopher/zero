// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package watcher

import (
	"context"
	"errors"
	"sync"

	"github.com/robfig/cron/v3"

	reflectutil "github.com/superproj/zero/pkg/util/reflect"
)

const (
	Every3Seconds = "@every 3s"
)

// Watcher is the interface for watchers. It use cron job as a scheduling engine.
type Watcher interface {
	Init(ctx context.Context, config *Config) error
	cron.Job
}

// Spec interface provides methods to set spec for a cron job.
type ISpec interface {
	// Spec return the spec for a cron job.
	// There are two cron spec formats in common usage:
	// - standard cron format: https://en.wikipedia.org/wiki/Cron
	// - quartz scheduler format: http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/tutorial-lesson-06.html
	// This method is optional for a watcher.
	Spec() string
}

var (
	registryLock = new(sync.Mutex)
	registry     = make(map[string]Watcher)
)

var (
	// ErrRegistered will be returned when watcher is already been registered.
	ErrRegistered = errors.New("watcher has already been registered")
	// ErrConfigUnavailable will be returned when the configuration input is not the expected type.
	ErrConfigUnavailable = errors.New("configuration is not available")
)

// Register registers a watcher and save in global variable `registry`.
func Register(watcher Watcher) {
	registryLock.Lock()
	defer registryLock.Unlock()

	name := reflectutil.StructName(watcher)
	if _, ok := registry[name]; ok {
		panic("duplicate watcher entry: " + name)
	}

	registry[name] = watcher
}

// ListWatchers returns registered watchers in map format.
func ListWatchers() map[string]Watcher {
	registryLock.Lock()
	defer registryLock.Unlock()

	return registry
}
