// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package app implements a server that runs a set of active components.
package app

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/component-base/term"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/superproj/zero/cmd/zero-miner-controller/app/config"
	"github.com/superproj/zero/cmd/zero-miner-controller/app/options"
	zerocontroller "github.com/superproj/zero/internal/controller"
	contextutil "github.com/superproj/zero/internal/pkg/util/context"
	"github.com/superproj/zero/internal/pkg/util/ratelimiter"
	"github.com/superproj/zero/pkg/apis/apps/v1beta1"
	"github.com/superproj/zero/pkg/apis/apps/v1beta1/index"
	"github.com/superproj/zero/pkg/db"
	"github.com/superproj/zero/pkg/util/record"
	"github.com/superproj/zero/pkg/version"
)

const appName = "zero-miner-controller"

func init() {
	utilruntime.Must(logsapi.AddFeatureGates(utilfeature.DefaultMutableFeatureGate))
}

// NewControllerCommand creates a *cobra.Command object with default parameters.
func NewControllerCommand() *cobra.Command {
	o, err := options.NewOptions()
	if err != nil {
		klog.Fatalf("Unable to initialize command options: %v", err)
	}

	cmd := &cobra.Command{
		Use: appName,
		Long: `The cloud miner controller is a daemon that embeds
the core control loops. In applications of robotics and
automation, a control loop is a non-terminating loop that regulates the state of
the system. In Zero, a controller is a control loop that watches the shared
state of the miner through the zero-apiserver and makes changes attempting to move the
current state towards the desired state.`,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			// silence client-go warnings.
			// zero-miner-controller generically watches APIs (including deprecated ones),
			// and CI ensures it works properly against matching zero-apiserver versions.
			restclient.SetDefaultWarningHandler(restclient.NoWarnings{})
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			version.PrintAndExitIfRequested(appName)

			// Activate logging as soon as possible, after that
			// show flags with the final logging configuration.
			if err := logsapi.ValidateAndApply(o.Logs, utilfeature.DefaultFeatureGate); err != nil {
				return err
			}
			cliflag.PrintFlags(cmd.Flags())

			if err := o.Complete(); err != nil {
				return err
			}

			if err := o.Validate(); err != nil {
				return err
			}

			c, err := o.Config()
			if err != nil {
				return err
			}

			cc := c.Complete()

			if err := options.LogOrWriteConfig(o.WriteConfigTo, cc.ComponentConfig); err != nil {
				return err
			}

			if err := Run(cc, wait.NeverStop); err != nil {
				return err
			}

			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	fs := cmd.Flags()
	namedFlagSets := o.Flags()
	version.AddFlags(namedFlagSets.FlagSet("global"))
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name(), logs.SkipLoggingConfigurationFlags())
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cliflag.SetUsageAndHelpFunc(cmd, namedFlagSets, cols)

	return cmd
}

// Run runs the controller options. This should never exit.
func Run(c *config.CompletedConfig, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Get())

	// Do some initialization here
	var redisOptions db.RedisOptions
	_ = copier.Copy(&redisOptions, &c.ComponentConfig.Redis)
	rdb, err := db.NewRedis(&redisOptions)
	if err != nil {
		return err
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := ctrl.NewManager(c.Kubeconfig, ctrl.Options{
		MetricsBindAddress:         c.ComponentConfig.MetricsBindAddress,
		Namespace:                  c.ComponentConfig.Namespace,
		HealthProbeBindAddress:     c.ComponentConfig.HealthzBindAddress,
		SyncPeriod:                 &c.ComponentConfig.SyncPeriod.Duration,
		LeaderElection:             c.ComponentConfig.LeaderElection.LeaderElect,
		LeaderElectionResourceLock: c.ComponentConfig.LeaderElection.ResourceLock,
		LeaderElectionNamespace:    c.ComponentConfig.LeaderElection.ResourceNamespace,
		LeaderElectionID:           c.ComponentConfig.LeaderElection.ResourceName,
		LeaseDuration:              &c.ComponentConfig.LeaderElection.LeaseDuration.Duration,
		RetryPeriod:                &c.ComponentConfig.LeaderElection.RetryPeriod.Duration,
		RenewDeadline:              &c.ComponentConfig.LeaderElection.RenewDeadline.Duration,
	})
	if err != nil {
		klog.ErrorS(err, "Unable to new miner controller")
		return err
	}

	// applies all the stored functions to the scheme created by controller-runtime
	_ = v1beta1.AddToScheme(mgr.GetScheme())
	_ = corev1.AddToScheme(mgr.GetScheme())

	// klog.Background will automatically use the right logger.
	ctrl.SetLogger(klog.Background())

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("zero-miner-controller"))

	// TODO: Switch to `wait.ContextForChannel` when wait package support `ContextForChannel` function
	ctx, _ := contextutil.ContextForChannel(stopCh)

	// controller-runtime for multi-cluster support, reference:
	// https://github.com/kubernetes-sigs/controller-runtime/blob/main/designs/move-cluster-specific-code-out-of-manager.md
	providerClient, providerCluster, err := prepareProvider(c.ProviderKubeconfig)
	if err != nil {
		return err
	}
	if err := mgr.Add(providerCluster); err != nil {
		return err
	}

	if err := index.AddDefaultIndexes(ctx, mgr); err != nil {
		klog.ErrorS(err, "unable to setup indexes")
		return err
	}

	if err = (&zerocontroller.MinerReconciler{
		ProviderClient:   providerClient,
		ProviderCluster:  providerCluster,
		RedisClient:      rdb,
		ComponentConfig:  c.ComponentConfig,
		WatchFilterValue: c.ComponentConfig.WatchFilterValue,
	}).SetupWithManager(ctx, mgr, controller.Options{
		MaxConcurrentReconciles: int(c.ComponentConfig.Parallelism),
		RecoverPanic:            pointer.Bool(true),
		RateLimiter:             ratelimiter.DefaultControllerRateLimiter(),
	}); err != nil {
		klog.ErrorS(err, "Unable to create controller", "controller", "miner")
		return err
	}

	// add handlers
	if err := mgr.AddReadyzCheck("ping", healthz.Ping); err != nil {
		klog.ErrorS(err, "Unable to set up health check")
		return err
	}

	if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
		klog.ErrorS(err, "Unable to set up ready check")
		return err
	}

	// Start the Cmd
	klog.InfoS("Starting miner controller")

	return mgr.Start(ctx)
}

func prepareProvider(kubeconfig *restclient.Config) (kubernetes.Interface, cluster.Cluster, error) {
	providerClient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	providerCluster, err := cluster.New(kubeconfig)
	if err != nil {
		return providerClient, nil, err
	}

	return providerClient, providerCluster, nil
}
