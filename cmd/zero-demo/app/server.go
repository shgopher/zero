// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package app

import (
	genericapiserver "k8s.io/apiserver/pkg/server"

	"github.com/superproj/zero/cmd/zero-demo/app/options"
	"github.com/superproj/zero/internal/demo"
	"github.com/superproj/zero/pkg/app"
)

const commandDesc = `The demo server is a standard, specification-compliant demo 
example of the zero service.

Find more zero-demo information at:
    https://github.com/superproj/zero/blob/master/docs/guide/en-US/cmd/zero-demo.md`

// NewApp creates an App object with default parameters.
func NewApp(name string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp(name, "Launch a zero demo server",
		app.WithDescription(commandDesc),
		app.WithOptions(opts),
		// app.WithLogFeatureGate(feature.DefaultFeatureGate),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

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
func Run(c *demo.Config, stopCh <-chan struct{}) error {
	server, err := c.Complete().New(stopCh)
	if err != nil {
		return err
	}

	return server.Run()
}
