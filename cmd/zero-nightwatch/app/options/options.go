// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package options contains flags and options for initializing an nightwatch.
package options

import (
	"github.com/spf13/viper"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/superproj/zero/internal/nightwatch"
	"github.com/superproj/zero/internal/pkg/feature"
	kubeconfigutil "github.com/superproj/zero/internal/pkg/util/kubeconfig"
	"github.com/superproj/zero/pkg/app"
	clientset "github.com/superproj/zero/pkg/generated/clientset/versioned"
	"github.com/superproj/zero/pkg/log"
	genericoptions "github.com/superproj/zero/pkg/options"
)

const (
	// UserAgent is the userAgent name when starting zero-nightwatch server.
	UserAgent = "zero-nightwatch"
)

var _ app.CliOptions = (*Options)(nil)

// Options contains everything necessary to create and run a nightwatch server.
type Options struct {
	HealthOptions *genericoptions.HealthOptions  `json:"health" mapstructure:"health"`
	MySQLOptions  *genericoptions.MySQLOptions   `json:"db" mapstructure:"db"`
	RedisOptions  *genericoptions.RedisOptions   `json:"redis" mapstructure:"redis"`
	Metrics       *genericoptions.MetricsOptions `json:"metrics" mapstructure:"metrics"`
	// Path to kubeconfig file with authorization and master location information.
	Kubeconfig   string          `json:"kubeconfig" mapstructure:"kubeconfig"`
	FeatureGates map[string]bool `json:"feature-gates"`
	Log          *log.Options    `json:"log" mapstructure:"log"`
}

// NewOptions returns initialized Options.
func NewOptions() *Options {
	o := &Options{
		HealthOptions: genericoptions.NewHealthOptions(),
		MySQLOptions:  genericoptions.NewMySQLOptions(),
		RedisOptions:  genericoptions.NewRedisOptions(),
		Metrics:       genericoptions.NewMetricsOptions(),
		Log:           log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.HealthOptions.AddFlags(fss.FlagSet("health"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.Metrics.AddFlags(fss.FlagSet("metrics"))
	o.Log.AddFlags(fss.FlagSet("log"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	feature.DefaultMutableFeatureGate.AddFlag(fs)

	return fss
}

// Complete completes all the required options.
func (o *Options) Complete() error {
	if err := viper.Unmarshal(&o); err != nil {
		return err
	}

	_ = feature.DefaultMutableFeatureGate.SetFromMap(o.FeatureGates)
	return nil
}

// Validate validates all the required options.
func (o *Options) Validate() error {
	errs := []error{}

	errs = append(errs, o.HealthOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.Metrics.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up zero-nightwatch config with options.
func (o *Options) ApplyTo(c *nightwatch.Config) error {
	c.HealthOptions = o.HealthOptions
	c.MySQLOptions = o.MySQLOptions
	c.RedisOptions = o.RedisOptions
	return nil
}

// Config return a zero-nightwatch config object.
func (o *Options) Config() (*nightwatch.Config, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", o.Kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeconfigutil.SetDefaultClientOptions(restclient.AddUserAgent(kubeconfig, UserAgent))

	client, err := clientset.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	c := &nightwatch.Config{
		Client: client,
	}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}
