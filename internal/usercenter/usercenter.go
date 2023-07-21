// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package usercenter

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/jinzhu/copier"

	"github.com/superproj/zero/internal/pkg/bootstrap"
	"github.com/superproj/zero/internal/usercenter/server"
	"github.com/superproj/zero/pkg/db"
	genericoptions "github.com/superproj/zero/pkg/options"
	"github.com/superproj/zero/pkg/version"
)

var (
	// Name is the name of the compiled software.
	Name = "zero-usercenter"

	// ID contains the host name and any error encountered during the retrieval.
	ID, _ = os.Hostname()
)

// Config represents the configuration of the service.
type Config struct {
	GRPCOptions   *genericoptions.GRPCOptions
	HTTPOptions   *genericoptions.HTTPOptions
	TLSOptions    *genericoptions.TLSOptions
	JWTOptions    *genericoptions.JWTOptions
	MySQLOptions  *genericoptions.MySQLOptions
	RedisOptions  *genericoptions.RedisOptions
	EtcdOptions   *genericoptions.EtcdOptions
	KafkaOptions  *genericoptions.KafkaOptions
	JaegerOptions *genericoptions.JaegerOptions
	ConsulOptions *genericoptions.ConsulOptions
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

// completedConfig holds the configuration after it has been completed.
type completedConfig struct {
	*Config
}

// New returns a new instance of Server from the given config.
func (c completedConfig) New(stopCh <-chan struct{}) (*Server, error) {
	if err := c.JaegerOptions.SetTracerProvider(); err != nil {
		return nil, err
	}

	appInfo := bootstrap.NewAppInfo(ID, Name, version.Get().String())

	conf := &server.Config{
		HTTP: *c.HTTPOptions,
		GRPC: *c.GRPCOptions,
		TLS:  *c.TLSOptions,
	}

	var dbOptions db.MySQLOptions
	_ = copier.Copy(&dbOptions, c.MySQLOptions)

	// Initialize Kratos application with the provided configurations.
	app, cleanup, err := wireApp(appInfo, conf, &dbOptions, c.JWTOptions, c.RedisOptions, c.EtcdOptions, c.KafkaOptions)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	return &Server{App: app}, nil
}

// Server represents the server.
type Server struct {
	App *kratos.App
}

// Run is a method of the Server struct that starts the server.
func (s *Server) Run() error {
	return s.App.Run()
}
