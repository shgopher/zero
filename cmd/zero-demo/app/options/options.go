// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package options contains flags and options for initializing an apiserver
package options

import (
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/superproj/zero/internal/demo"
	"github.com/superproj/zero/internal/pkg/feature"
	"github.com/superproj/zero/pkg/app"
	"github.com/superproj/zero/pkg/log"
	genericoptions "github.com/superproj/zero/pkg/options"
)

const (
	// UserAgent is the userAgent name when starting zero-demo server.
	UserAgent = "zero-demo"
)

var _ app.CliOptions = (*Options)(nil)

// Options contains state for master/api server.
type Options struct {
	GRPCOptions   *genericoptions.GRPCOptions    `json:"grpc" mapstructure:"grpc"`
	HTTPOptions   *genericoptions.HTTPOptions    `json:"http" mapstructure:"http"`
	TLSOptions    *genericoptions.TLSOptions     `json:"tls" mapstructure:"tls"`
	MySQLOptions  *genericoptions.MySQLOptions   `json:"db" mapstructure:"db"`
	RedisOptions  *genericoptions.RedisOptions   `json:"redis" mapstructure:"redis"`
	JaegerOptions *genericoptions.JaegerOptions  `json:"jaeger" mapstructure:"jaeger"`
	JWTOptions    *genericoptions.JWTOptions     `json:"jwt" mapstructure:"jwt"`
	Metrics       *genericoptions.MetricsOptions `json:"metrics" mapstructure:"metrics"`
	Log           *log.Options                   `json:"log" mapstructure:"log"`

	EnableTLS bool `json:"enable-tls" mapstructure:"enable-tls"`
	// Path to kubeconfig file with authorization and master location information.
	Kubeconfig   string          `json:"kubeconfig" mapstructure:"kubeconfig"`
	FeatureGates map[string]bool `json:"feature-gates" mapstructure:"-"`
}

// NewOptions returns initialized Options.
func NewOptions() *Options {
	o := &Options{
		GRPCOptions:   genericoptions.NewGRPCOptions(),
		HTTPOptions:   genericoptions.NewHTTPOptions(),
		TLSOptions:    genericoptions.NewTLSOptions(),
		MySQLOptions:  genericoptions.NewMySQLOptions(),
		RedisOptions:  genericoptions.NewRedisOptions(),
		JaegerOptions: genericoptions.NewJaegerOptions(),
		JWTOptions:    genericoptions.NewJWTOptions(),
		Metrics:       genericoptions.NewMetricsOptions(),
		Log:           log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.HTTPOptions.AddFlags(fss.FlagSet("http"))
	o.TLSOptions.AddFlags(fss.FlagSet("tls"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.JaegerOptions.AddFlags(fss.FlagSet("jaeger"))
	o.JWTOptions.AddFlags(fss.FlagSet("jwt"))
	o.Metrics.AddFlags(fss.FlagSet("metrics"))
	o.Log.AddFlags(fss.FlagSet("log"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.BoolVar(&o.EnableTLS, "enable-tls", o.EnableTLS, "Enable TLS for grpc server.")
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	feature.DefaultMutableFeatureGate.AddFlag(fs)

	return fss
}

// Complete completes all the required options.
func (o *Options) Complete() error {
	if o.JaegerOptions.ServiceName == "" {
		o.JaegerOptions.ServiceName = UserAgent
	}
	_ = feature.DefaultMutableFeatureGate.SetFromMap(o.FeatureGates)
	return nil
}

// Validate validates all the required options.
func (o *Options) Validate() error {
	errs := []error{}

	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.JaegerOptions.Validate()...)
	errs = append(errs, o.JWTOptions.Validate()...)
	errs = append(errs, o.Metrics.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up zero-demo config with options.
func (o *Options) ApplyTo(c *demo.Config) error {
	c.GRPCOptions = o.GRPCOptions
	c.HTTPOptions = o.HTTPOptions
	c.TLSOptions = o.TLSOptions
	c.MySQLOptions = o.MySQLOptions
	c.RedisOptions = o.RedisOptions
	c.JaegerOptions = o.JaegerOptions
	c.JWTOptions = o.JWTOptions

	return nil
}

// Config return a zero-demo config object.
func (o *Options) Config() (*demo.Config, error) {
	c := &demo.Config{}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}
