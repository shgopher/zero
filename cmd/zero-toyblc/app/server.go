// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package app

import (
	genericapiserver "k8s.io/apiserver/pkg/server"

	"github.com/superproj/zero/cmd/zero-toyblc/app/options"
	"github.com/superproj/zero/internal/toyblc"
	"github.com/superproj/zero/pkg/app"
)

// Define the description of the command.
const commandDesc = `The toy blc is used to start a naive and simple blockchain node.`

// NewApp creates and returns a new App object with default parameters.
func NewApp(name string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp(name, "Launch a zero toy blockchain node",
		app.WithDescription(commandDesc),
		app.WithOptions(opts),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

// Returns the function to run the application.
func run(opts *options.Options) app.RunFunc {
	return func() error {
		cfg, err := opts.Config()
		if err != nil {
			return err
		}

		return Run(cfg, genericapiserver.SetupSignalHandler())
	}
}

// Run runs the specified APIServer. This should never exit.
func Run(c *toyblc.Config, stopCh <-chan struct{}) error {
	server, err := c.Complete().New()
	if err != nil {
		return err
	}

	return server.Run(stopCh)
}
