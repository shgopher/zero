// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package gateway

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/jinzhu/copier"
	"k8s.io/client-go/rest"

	"github.com/superproj/zero/internal/gateway/server"
	"github.com/superproj/zero/internal/pkg/bootstrap"
	"github.com/superproj/zero/internal/pkg/client/usercenter"
	"github.com/superproj/zero/pkg/db"
	"github.com/superproj/zero/pkg/generated/clientset/versioned"
	informers "github.com/superproj/zero/pkg/generated/informers/externalversions"
	"github.com/superproj/zero/pkg/log"
	genericoptions "github.com/superproj/zero/pkg/options"
	"github.com/superproj/zero/pkg/version"
)

var (
	// Name is the name of the compiled software.
	Name = "zero-gateway"

	ID, _ = os.Hostname()
)

// Config defines the config for the apiserver.
type Config struct {
	GRPCOptions       *genericoptions.GRPCOptions
	HTTPOptions       *genericoptions.HTTPOptions
	TLSOptions        *genericoptions.TLSOptions
	UserCenterOptions *usercenter.UserCenterOptions
	MySQLOptions      *genericoptions.MySQLOptions
	RedisOptions      *genericoptions.RedisOptions
	EtcdOptions       *genericoptions.EtcdOptions
	JaegerOptions     *genericoptions.JaegerOptions
	ConsulOptions     *genericoptions.ConsulOptions

	// the rest config for the zero-apiserver
	Kubeconfig *rest.Config
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

	appInfo := bootstrap.NewAppInfo(ID, Name, version.Get().String())

	conf := &server.Config{
		HTTP: *c.HTTPOptions,
		GRPC: *c.GRPCOptions,
		TLS:  *c.TLSOptions,
	}

	// You can use `sigs.k8s.io/controller-runtime/pkg/client`.New to created a client
	// which can support unstructured types also.
	// cl, err := client.New(c.Kubeconfig, client.Options{})
	client, err := versioned.NewForConfig(c.Kubeconfig)
	if err != nil {
		log.Errorw(err, "Cannot create connection with zero-apiserver")
		return nil, err
	}

	var mysqlOptions db.MySQLOptions
	var redisOptions db.RedisOptions
	_ = copier.Copy(&mysqlOptions, c.MySQLOptions)
	_ = copier.Copy(&redisOptions, c.RedisOptions)

	app, cleanup, err := wireApp(stopCh, appInfo, conf, client, &mysqlOptions, &redisOptions, c.UserCenterOptions, c.RedisOptions, c.EtcdOptions)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	return &Server{App: app}, nil
}

// // Server represents the server.
type Server struct {
	App                   *kratos.App
	SharedInformerFactory informers.SharedInformerFactory
}

func (s *Server) Run() error {
	return s.App.Run()
}
