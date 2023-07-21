// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package options provides the flags used for the controller manager.
package options

import (
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/controller/garbagecollector"

	controllermanagerconfig "github.com/superproj/zero/cmd/zero-controller-manager/app/config"
	ctrlmgrconfig "github.com/superproj/zero/internal/controller/apis/config"
	"github.com/superproj/zero/internal/controller/apis/config/latest"
	"github.com/superproj/zero/internal/controller/apis/config/validation"
	kubeconfigutil "github.com/superproj/zero/internal/pkg/util/kubeconfig"
	clientset "github.com/superproj/zero/pkg/generated/clientset/versioned"
)

const (
	// ControllerManagerUserAgent is the userAgent name when starting zero-controller managers.
	ControllerManagerUserAgent = "zero-controller-manager"
)

// Options is the main context object for the zero-controller manager.
type Options struct {
	// ConfigFile is the location of the miner controller server's configuration file.
	ConfigFile string

	// WriteConfigTo is the path where the default configuration will be written.
	WriteConfigTo string

	// Generic                    *GenericControllerManagerConfigurationOptions
	// GarbageCollectorController *GarbageCollectorControllerOptions
	// ChainController            *ChainControllerOptions
	// NamespaceController *NamespaceControllerOptions

	// The address of the Kubernetes API server (overrides any value in kubeconfig).
	Master string
	// Path to kubeconfig file with authorization and master location information.
	Kubeconfig string
	Logs       *logs.Options

	// config is the zero controller manager server's configuration object.
	// The default values.
	config *ctrlmgrconfig.ZeroControllerManagerConfiguration
}

// NewOptions creates a new Options with a default config.
func NewOptions() (*Options, error) {
	defaultComponentConfig, err := latest.Default()
	if err != nil {
		return nil, err
	}

	o := Options{
		Logs:   logs.NewOptions(),
		config: defaultComponentConfig,
		/*
			Cloud:              cloud.NewCloudOptions(),
			Concurrency:        1,
			Logs:               logs.NewOptions(),
			Metrics:        metrics.NewOptions(),
			MetricsBindAddress: metrics.DefaultMetricsAddress,
			HealthAddr:         ":9441",
			LeaderElection: componentbaseconfig.LeaderElectionConfiguration{
				LeaseDuration:     metav1.Duration{Duration: 15 * time.Second},
				RenewDeadline:     metav1.Duration{Duration: 10 * time.Second},
				RetryPeriod:       metav1.Duration{Duration: 2 * time.Second},
				LeaderElect:       true,
				ResourceName:      "zero-controller-manager",
				ResourceNamespace: "kube-system",
			},
		*/
	}

	gcIgnoredResources := make([]ctrlmgrconfig.GroupResource, 0, len(garbagecollector.DefaultIgnoredResources()))
	for r := range garbagecollector.DefaultIgnoredResources() {
		gcIgnoredResources = append(gcIgnoredResources, ctrlmgrconfig.GroupResource{Group: r.Group, Resource: r.Resource})
	}
	o.config.GarbageCollectorController.GCIgnoredResources = gcIgnoredResources
	o.config.Generic.LeaderElection.ResourceName = "zero-controller-manager"
	o.config.Generic.LeaderElection.ResourceNamespace = "kube-system"

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
	(&GenericControllerManagerConfigurationOptions{&o.config.Generic}).AddFlags(&fss)
	(&GarbageCollectorControllerOptions{&o.config.GarbageCollectorController}).AddFlags(fss.FlagSet("garbage collector controller"))
	(&ChainControllerOptions{&o.config.ChainController}).AddFlags(fss.FlagSet("chain controller"))

	fs := fss.FlagSet("misc")
	fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "The path to the configuration file.")
	fs.StringVar(&o.WriteConfigTo, "write-config-to", o.WriteConfigTo, "If set, write the default configuration values to this file and exit.")
	fs.StringVar(&o.Master, "master", o.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig).")
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	logsapi.AddFlags(o.Logs, fss.FlagSet("logs"))
	utilfeature.DefaultMutableFeatureGate.AddFlag(fss.FlagSet("generic"))

	return fss
}

// Validate is used to validate the options and config before launching the controller.
func (o *Options) Validate() error {
	var errs []error

	if err := validation.Validate(o.config).ToAggregate(); err != nil {
		errs = append(errs, err.Errors()...)
	}

	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up zero controller manager config with options.
func (o *Options) ApplyTo(c *controllermanagerconfig.Config) error {
	c.ComponentConfig = o.config

	return nil
}

// Config return a controller manager config objective.
func (o Options) Config() (*controllermanagerconfig.Config, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags(o.Master, o.Kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(restclient.AddUserAgent(kubeconfig, ControllerManagerUserAgent))
	if err != nil {
		return nil, err
	}

	c := &controllermanagerconfig.Config{
		Kubeconfig: kubeconfigutil.SetClientOptionsForController(restclient.AddUserAgent(kubeconfig, ControllerManagerUserAgent)),
		Client:     client,
	}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}
