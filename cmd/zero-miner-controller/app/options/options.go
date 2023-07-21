// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package options provides the flags used for the miner controller.
package options

import (
	"fmt"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	componentbaseoptions "k8s.io/component-base/config/options"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog/v2"

	controllerconfig "github.com/superproj/zero/cmd/zero-miner-controller/app/config"
	minercontrollerconfig "github.com/superproj/zero/internal/controller/miner/apis/config"
	"github.com/superproj/zero/internal/controller/miner/apis/config/latest"
	"github.com/superproj/zero/internal/controller/miner/apis/config/validation"
	kubeconfigutil "github.com/superproj/zero/internal/pkg/util/kubeconfig"
)

const (
	// ControllerUserAgent is the userAgent name when starting zero-miner controller.
	ControllerUserAgent = "zero-miner-controller"
)

// Options is the main context object for the zero-miner controller.
type Options struct {
	// ConfigFile is the location of the miner controller server's configuration file.
	ConfigFile string

	// WriteConfigTo is the path where the default configuration will be written.
	WriteConfigTo string

	// The address of the Kubernetes API server (overrides any value in kubeconfig).
	Master string

	// Path to kubeconfig file with authorization and master location information.
	Kubeconfig string

	// Path to miner provider kubeconfig file with authorization and master location information.
	ProviderKubeconfig string

	InCluster bool

	Logs *logs.Options

	// config is the miner controller server's configuration object.
	// The default values.
	config *minercontrollerconfig.MinerControllerConfiguration
}

// NewOptions creates a new Options with a default config.
func NewOptions() (*Options, error) {
	o := Options{
		Logs: logs.NewOptions(),
	}

	defaultComponentConfig, err := latest.Default()
	if err != nil {
		return nil, err
	}
	o.config = defaultComponentConfig

	return &o, nil
}

// Complete completes all the required options.
func (o *Options) Complete() error {
	if len(o.ConfigFile) == 0 && len(o.WriteConfigTo) == 0 {
		klog.InfoS("Warning, all flags other than --config, --write-config-to are deprecated, please begin using a config file ASAP")
	}

	if len(o.ConfigFile) > 0 {
		cfg, err := loadConfigFromFile(o.ConfigFile)
		if err != nil {
			return err
		}

		o.config = cfg
	}

	return utilfeature.DefaultMutableFeatureGate.SetFromMap(o.config.FeatureGates)
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	// o.Logs.AddFlags(fss.FlagSet("logs"))
	componentbaseoptions.BindLeaderElectionFlags(&o.config.LeaderElection, fss.FlagSet("leader elect"))
	///o.config.Cloud.AddFlags(fss.FlagSet("cloud"))

	fs := fss.FlagSet("misc")
	fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "The path to the configuration file.")
	fs.StringVar(&o.WriteConfigTo, "write-config-to", o.WriteConfigTo, "If set, write the default configuration values to this file and exit.")
	fs.StringVar(&o.Master, "master", o.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig).")
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	fs.BoolVar(&o.InCluster, "in-cluster", o.InCluster, "Create miner pod in the cluster where controller is located.")
	fs.StringVar(&o.ProviderKubeconfig, "provider-kubeconfig", o.ProviderKubeconfig, ""+
		"Path to miner provider kubeconfig file with authorization and master location information.")

	logsapi.AddFlags(o.Logs, fss.FlagSet("logs"))
	utilfeature.DefaultMutableFeatureGate.AddFlag(fss.FlagSet("generic"))

	return fss
}

// Validate is used to validate the options and config before launching the controller.
func (o *Options) Validate() error {
	var errs []error

	if !o.InCluster && o.ProviderKubeconfig == "" {
		errs = append(errs, fmt.Errorf("`provider-kubeconfig` must be provided when `in-cluster` is false"))
	}

	if err := validation.Validate(o.config).ToAggregate(); err != nil {
		errs = append(errs, err.Errors()...)
	}

	// TODO: validate master and kubeconfig
	if o.config.Parallelism <= 0 {
		errs = append(errs, fmt.Errorf("--parallelism must be greater than or equal to 0"))
	}

	// errs = append(errs, o.Cloud.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up miner controller config with options.
func (o *Options) ApplyTo(c *controllerconfig.Config) error {
	c.ComponentConfig = o.config
	return nil
}

// Config return a miner controller config objective.
func (o *Options) Config() (*controllerconfig.Config, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags(o.Master, o.Kubeconfig)
	if err != nil {
		return nil, err
	}

	var providerKubeconfig *restclient.Config
	if !o.InCluster {
		providerKubeconfig, err = clientcmd.BuildConfigFromFlags("", o.ProviderKubeconfig)
	} else {
		providerKubeconfig, err = restclient.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	c := &controllerconfig.Config{
		Kubeconfig:         kubeconfigutil.SetClientOptionsForController(restclient.AddUserAgent(kubeconfig, ControllerUserAgent)),
		ProviderKubeconfig: kubeconfigutil.SetClientOptionsForController(restclient.AddUserAgent(providerKubeconfig, ControllerUserAgent)),
	}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}
