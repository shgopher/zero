// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package config

import (
	restclient "k8s.io/client-go/rest"

	"github.com/superproj/zero/internal/controller/minerset/apis/config"
)

// Config is the main context object for the controller.
type Config struct {
	ComponentConfig *config.MinerSetControllerConfiguration

	// the rest config for the master
	Kubeconfig *restclient.Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	*Config
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() *CompletedConfig {
	return &CompletedConfig{c}
}
