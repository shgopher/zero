// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package app implements a server that runs a set of active components.
package app

import (
	"context"
	"fmt"
	"os"

	"github.com/jinzhu/copier"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/metadata"
	restclient "k8s.io/client-go/rest"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/component-base/term"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	ctrlmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/superproj/zero/pkg/version"

	"github.com/superproj/zero/cmd/zero-controller-manager/app/cleaner"
	"github.com/superproj/zero/cmd/zero-controller-manager/app/config"
	"github.com/superproj/zero/cmd/zero-controller-manager/app/options"
	zerocontroller "github.com/superproj/zero/internal/controller"
	configv1beta1 "github.com/superproj/zero/internal/controller/apis/config/v1beta1"
	"github.com/superproj/zero/internal/gateway/store"
	"github.com/superproj/zero/internal/pkg/metrics"
	contextutil "github.com/superproj/zero/internal/pkg/util/context"
	"github.com/superproj/zero/internal/pkg/util/ratelimiter"
	appsv1beta1 "github.com/superproj/zero/pkg/apis/apps/v1beta1"
	apiv1 "github.com/superproj/zero/pkg/apis/core/v1"
	"github.com/superproj/zero/pkg/db"
	"github.com/superproj/zero/pkg/util/record"
)

const appName = "zero-controller-manager"

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(logsapi.AddFeatureGates(utilfeature.DefaultMutableFeatureGate))
	// utilruntime.Must(features.AddFeatureGates(utilfeature.DefaultMutableFeatureGate)) TODO: add this line to keep up with the latest kube-scheduler

	// applies all the stored functions to the scheme created by controller-runtime
	_ = apiv1.AddToScheme(scheme)
	_ = appsv1beta1.AddToScheme(scheme)
	_ = configv1beta1.AddToScheme(scheme)
	// _ = corev1.AddToScheme(scheme)
}

// NewControllerManagerCommand creates a *cobra.Command object with default parameters.
func NewControllerManagerCommand() *cobra.Command {
	o, err := options.NewOptions()
	if err != nil {
		klog.Fatalf("Unable to initialize command options: %v", err)
	}

	cmd := &cobra.Command{
		Use: appName,
		Long: `The zero controller manager is a daemon that embeds
the core control loops. In applications of robotics and
automation, a control loop is a non-terminating loop that regulates the state of
the system. In Zero , a controller is a control loop that watches the shared
state of the miner through the zero-apiserver and makes changes attempting to move the
current state towards the desired state.`,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			// silence client-go warnings.
			// zero-controller-manager generically watches APIs (including deprecated ones),
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

			// add feature enablement metrics
			// utilfeature.DefaultMutableFeatureGate.AddMetrics() TODO: add this to keep up with the latest kube-controller-manager
			return Run(c.Complete(), wait.NeverStop)
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

// Run runs the controller manager options. This should never exit.
func Run(c *config.CompletedConfig, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Get())

	klog.InfoS("Golang settings", "GOGC", os.Getenv("GOGC"), "GOMAXPROCS", os.Getenv("GOMAXPROCS"), "GOTRACEBACK", os.Getenv("GOTRACEBACK"))

	// Do some initialization here
	var mysqlOptions db.MySQLOptions
	_ = copier.Copy(&mysqlOptions, c.ComponentConfig.Generic.MySQL)
	storeClient, err := wireStoreClient(&mysqlOptions)
	if err != nil {
		return err
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := ctrl.NewManager(c.Kubeconfig, ctrl.Options{
		Scheme:                     scheme,
		MetricsBindAddress:         c.ComponentConfig.Generic.MetricsBindAddress,
		Namespace:                  c.ComponentConfig.Generic.Namespace,
		HealthProbeBindAddress:     c.ComponentConfig.Generic.HealthzBindAddress,
		SyncPeriod:                 &c.ComponentConfig.Generic.SyncPeriod.Duration,
		LeaderElection:             c.ComponentConfig.Generic.LeaderElection.LeaderElect,
		LeaderElectionResourceLock: c.ComponentConfig.Generic.LeaderElection.ResourceLock,
		LeaderElectionNamespace:    c.ComponentConfig.Generic.LeaderElection.ResourceNamespace,
		LeaderElectionID:           c.ComponentConfig.Generic.LeaderElection.ResourceName,
		LeaseDuration:              &c.ComponentConfig.Generic.LeaderElection.LeaseDuration.Duration,
		RetryPeriod:                &c.ComponentConfig.Generic.LeaderElection.RetryPeriod.Duration,
		RenewDeadline:              &c.ComponentConfig.Generic.LeaderElection.RenewDeadline.Duration,
		ClientDisableCacheFor: []client.Object{
			&corev1.ConfigMap{},
			&corev1.Secret{},
		},
		// TLSOpts:                tlsOptionOverrides,
	})
	if err != nil {
		klog.ErrorS(err, "Unable to new controller manager")
		return err
	}

	machineMetricsCollector := metrics.NewMinerCollector(mgr.GetClient(), c.ComponentConfig.Generic.Namespace)
	ctrlmetrics.Registry.MustRegister(machineMetricsCollector)

	// klog.Background will automatically use the right logger.
	ctrl.SetLogger(klog.Background())

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("zero-controller-manager"))

	// TODO: Switch to `wait.ContextForChannel` when wait package support `ContextForChannel` function
	ctx, _ := contextutil.ContextForChannel(stopCh)

	// setup resource cleaner controller
	clean := newResourceCleaner(mgr.GetClient(), storeClient, &cleaner.MinerCleaner{}, &cleaner.MinerSetCleaner{}, &cleaner.ChainCleaner{})
	if err := mgr.Add(clean); err != nil {
		klog.ErrorS(err, "Unable to create resource cleaner", "controller", "ResourceCleaner")
		return err
	}

	if err := setupReconcilers(ctx, c, storeClient, mgr); err != nil {
		return err
	}

	if err := setupChecks(mgr); err != nil {
		return err
	}

	// setupReconcilers(ctx, mgr)

	// Start the Cmd
	klog.InfoS("Starting controller manager")

	return mgr.Start(ctx)
}

func setupChecks(mgr ctrl.Manager) error {
	// add handlers
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		klog.ErrorS(err, "Unable to set up ready check")
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		klog.ErrorS(err, "Unable to set up health check")
		return err
	}

	return nil
}

func setupReconcilers(ctx context.Context, c *config.CompletedConfig, storeClient store.IStore, mgr ctrl.Manager) error {
	// setup garbage collector controller
	gc := &garbageCollector{completedConfig: c}
	if err := mgr.Add(gc); err != nil {
		klog.ErrorS(err, "Unable to create controller", "controller", "GarbageCollector")
		return err
	}

	defaultOptions := controller.Options{
		MaxConcurrentReconciles: int(c.ComponentConfig.Generic.Parallelism),
		RecoverPanic:            pointer.Bool(true),
		RateLimiter:             ratelimiter.DefaultControllerRateLimiter(),
	}

	// setup chain controller
	if err := (&zerocontroller.ChainReconciler{
		ComponentConfig:  &c.ComponentConfig.ChainController,
		WatchFilterValue: c.ComponentConfig.Generic.WatchFilterValue,
	}).SetupWithManager(ctx, mgr, defaultOptions); err != nil {
		klog.ErrorS(err, "Unable to create controller", "controller", "Chain")
		return err
	}

	// setup sync controller
	if err := (&zerocontroller.SyncReconciler{
		Store: storeClient,
	}).SetupWithManager(ctx, mgr, defaultOptions); err != nil {
		klog.ErrorS(err, "Unable to create controller", "controller", "Sync")
		return err
	}

	metadataClient, err := metadata.NewForConfig(c.Kubeconfig)
	if err != nil {
		klog.ErrorS(err, "Failed to create metadata client")
		return err
	}

	if err := (&zerocontroller.NamespacedResourcesDeleterReconciler{
		Client:         c.Client,
		MetadataClient: metadataClient,
	}).SetupWithManager(ctx, mgr, defaultOptions); err != nil {
		klog.ErrorS(err, "Unable to create controller", "controller", "Namespace")
		return err
	}

	return nil
}
