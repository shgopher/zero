// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package options

import (
	"github.com/spf13/pflag"

	chainconfig "github.com/superproj/zero/internal/controller/apis/config"
)

// ChainControllerOptions holds the ChainController options.
type ChainControllerOptions struct {
	*chainconfig.ChainControllerConfiguration
}

// AddFlags adds flags related to ChainController for controller manager to the specified FlagSet.
func (o *ChainControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.Image, "node-image", o.Image, "The blockchain node image used by default."+
		"This parameter is ignored if a config file is specified by --config.")
}
