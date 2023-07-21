// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package demo

import (
	"os"

	"github.com/go-kratos/kratos/v2"

	"github.com/superproj/zero/pkg/log"
	genericoptions "github.com/superproj/zero/pkg/options"
)

var (
	// Name is the name of the compiled software.
	Name = "zero-demo"

	ID, _ = os.Hostname()
)

// Config represents the configuration of the service.
type Config struct {
	GRPCOptions   *genericoptions.GRPCOptions
	HTTPOptions   *genericoptions.HTTPOptions
	TLSOptions    *genericoptions.TLSOptions
	MySQLOptions  *genericoptions.MySQLOptions
	RedisOptions  *genericoptions.RedisOptions
	JWTOptions    *genericoptions.JWTOptions
	JaegerOptions *genericoptions.JaegerOptions
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

type completedConfig struct {
	*Config
}

// New returns a new instance of Server from the given config.
func (c completedConfig) New(stopCh <-chan struct{}) (*Server, error) {
	if err := c.JaegerOptions.SetTracerProvider(); err != nil {
		return nil, err
	}
	log.Infow("New demo server", "httpAddr", c.HTTPOptions.Addr)

	return &Server{}, nil
}

// Server represents the server.
type Server struct {
	App *kratos.App
}

func (s *Server) Run() error {
	select {}
	// return s.App.Run()
}
