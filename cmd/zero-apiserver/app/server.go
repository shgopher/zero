// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//nolint:nakedret
package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
	oteltrace "go.opentelemetry.io/otel/trace"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	kversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/endpoints/discovery/aggregated"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericfeatures "k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/filters"
	serveroptions "k8s.io/apiserver/pkg/server/options"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/apiserver/pkg/util/openapi"
	"k8s.io/client-go/rest"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/component-base/term"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	netutils "k8s.io/utils/net"

	"github.com/superproj/zero/cmd/zero-apiserver/app/options"
	"github.com/superproj/zero/internal/admission/initializer"
	"github.com/superproj/zero/internal/apiserver"
	"github.com/superproj/zero/internal/apiserver/storage"
	"github.com/superproj/zero/internal/pkg/config/minertype"
	"github.com/superproj/zero/pkg/generated/clientset/versioned"
	informers "github.com/superproj/zero/pkg/generated/informers/externalversions"
	zeroopenapi "github.com/superproj/zero/pkg/generated/openapi"
	"github.com/superproj/zero/pkg/version"
)

func init() {
	utilruntime.Must(logsapi.AddFeatureGates(utilfeature.DefaultMutableFeatureGate))
}

// NewAPIServerCommand creates a *cobra.Command object with default parameters.
func NewAPIServerCommand() *cobra.Command {
	opts := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use:   "zero-apiserver",
		Short: "Launch a zero API server",
		Long: `The Zero API server validates and configures data
for the api objects which include miners, minersets, configmaps, and
others. The API Server services REST operations and provides the frontend to the
zero's shared state through which all other components interact.`,

		// stop printing usage when the command errors
		SilenceUsage: true,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			// silence client-go warnings.
			// zero-apiserver loopback clients should not log self-issued warnings.
			rest.SetDefaultWarningHandler(rest.NoWarnings{})
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			version.PrintAndExitIfRequested("zero-apiserver")
			fs := cmd.Flags()

			// Activate logging as soon as possible, after that
			// show flags with the final logging configuration.
			if err := logsapi.ValidateAndApply(opts.Logs, utilfeature.DefaultFeatureGate); err != nil {
				return err
			}
			cliflag.PrintFlags(fs)

			// set default options
			completedOptions, err := Complete(opts)
			if err != nil {
				return err
			}

			// validate options
			if errs := completedOptions.Validate(); len(errs) != 0 {
				return utilerrors.NewAggregate(errs)
			}
			// add feature enablement metrics
			utilfeature.DefaultMutableFeatureGate.AddMetrics()
			return Run(completedOptions, genericapiserver.SetupSignalHandler())
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
	namedFlagSets := opts.Flags()
	version.AddFlags(namedFlagSets.FlagSet("global"))
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name(), logs.SkipLoggingConfigurationFlags())
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cliflag.SetUsageAndHelpFunc(cmd, namedFlagSets, cols)

	return cmd
}

// Run runs the specified APIServer. This should never exit.
func Run(completedOptions completedServerRunOptions, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Get().String())
	klog.InfoS("Golang settings", "GOGC", os.Getenv("GOGC"), "GOMAXPROCS", os.Getenv("GOMAXPROCS"), "GOTRACEBACK", os.Getenv("GOTRACEBACK"))

	config, err := CreateAPIServerConfig(completedOptions)
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	server.GenericAPIServer.AddPostStartHookOrDie(
		"start-zero-informers",
		func(context genericapiserver.PostStartHookContext) error {
			// remove dependence with kube-apiserver
			completedOptions.SharedInformerFactory.Start(context.StopCh)
			return nil
		},
	)

	server.GenericAPIServer.AddPostStartHookOrDie(
		"initialize-instance-config-client",
		func(ctx genericapiserver.PostStartHookContext) error {
			client, err := versioned.NewForConfig(ctx.LoopbackClientConfig)
			if err != nil {
				return err
			}

			if err := minertype.Init(context.Background(), client); err != nil {
				// When returning 'NotFound' error, we should not report an error, otherwise we can not
				// create 'MinerTypesConfigMapName' configmap via zero-apiserver
				if apierrors.IsNotFound(err) {
					return nil
				}

				klog.ErrorS(err, "Failed to init miner type cache")
				return err
			}

			return nil
		},
	)

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}

// CreateAPIServerConfig creates all the resources for running the API server, but runs none of them.
func CreateAPIServerConfig(completedOptions completedServerRunOptions) (*apiserver.Config, error) {
	genericConfig, versionedInformers, storageFactory, err := buildGenericConfig(completedOptions)
	if err != nil {
		return nil, err
	}

	completedOptions.Metrics.Apply()

	config := &apiserver.Config{
		GenericConfig: genericConfig,
		ExtraConfig: apiserver.ExtraConfig{
			APIResourceConfigSource: storageFactory.APIResourceConfigSource,
			StorageFactory:          storageFactory,
			EventTTL:                completedOptions.EventTTL,
			EnableLogsSupport:       completedOptions.EnableLogsHandler,
			VersionedInformers:      versionedInformers,
		},
	}

	return config, nil
}

// BuildGenericConfig takes the master server options and produces the genericapiserver.Config associated with it.
func buildGenericConfig(s completedServerRunOptions) (
	genericConfig *genericapiserver.RecommendedConfig,
	versionedInformers informers.SharedInformerFactory,
	storageFactory *serverstorage.DefaultStorageFactory,
	lastErr error,
) {
	genericConfig = genericapiserver.NewRecommendedConfig(legacyscheme.Codecs)
	genericConfig.MergedResourceConfig = apiserver.DefaultAPIResourceConfigSource()

	if lastErr = s.GenericServerRunOptions.ApplyTo(&genericConfig.Config); lastErr != nil {
		return
	}

	if utilfeature.DefaultFeatureGate.Enabled(genericfeatures.APIServerTracing) {
		if lastErr = s.Traces.ApplyTo(genericConfig.EgressSelector, &genericConfig.Config); lastErr != nil {
			return
		}
	}

	// wrap the definitions to revert any changes from disabled features
	getOpenAPIDefinitions := openapi.GetOpenAPIDefinitionsWithoutDisabledFeatures(zeroopenapi.GetOpenAPIDefinitions)
	genericConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(
		getOpenAPIDefinitions,
		openapinamer.NewDefinitionNamer(legacyscheme.Scheme),
	)
	genericConfig.OpenAPIConfig.Info.Title = "Zero"
	genericConfig.OpenAPIConfig.Info.Version = "1.0"
	if utilfeature.DefaultFeatureGate.Enabled(genericfeatures.OpenAPIV3) {
		genericConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(
			getOpenAPIDefinitions,
			openapinamer.NewDefinitionNamer(legacyscheme.Scheme),
		)
		genericConfig.OpenAPIV3Config.Info.Title = "Zero"
	}
	// Placeholder
	genericConfig.LongRunningFunc = filters.BasicLongRunningRequestCheck(
		sets.NewString("watch", "proxy"), sets.NewString("attach", "exec", "proxy", "log", "portforward"),
	)
	genericConfig.Version = convertVersion()

	if genericConfig.EgressSelector != nil {
		s.RecommendedOptions.Etcd.StorageConfig.Transport.EgressLookup = genericConfig.EgressSelector.Lookup
	}
	if utilfeature.DefaultFeatureGate.Enabled(genericfeatures.APIServerTracing) {
		s.RecommendedOptions.Etcd.StorageConfig.Transport.TracerProvider = genericConfig.TracerProvider
	} else {
		s.RecommendedOptions.Etcd.StorageConfig.Transport.TracerProvider = oteltrace.NewNoopTracerProvider()
	}

	if lastErr = s.RecommendedOptions.Etcd.Complete(genericConfig.StorageObjectCountTracker, genericConfig.DrainedNotify(), genericConfig.AddPostStartHook); lastErr != nil {
		return
	}

	storageFactoryConfig := storage.NewStorageFactoryConfig()
	storageFactoryConfig.APIResourceConfig = genericConfig.MergedResourceConfig
	storageFactory, lastErr = storageFactoryConfig.Complete(s.RecommendedOptions.Etcd).New()
	if lastErr != nil {
		return
	}
	if lastErr = s.RecommendedOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = s.RecommendedOptions.Etcd.ApplyWithStorageFactoryTo(storageFactory, &genericConfig.Config); lastErr != nil {
		return
	}

	// Use protobufs for self-communication.
	// Since not every generic apiserver has to support protobufs, we
	// cannot default to it in generic apiserver and need to explicitly
	// set it in kube-apiserver.
	genericConfig.LoopbackClientConfig.ContentConfig.ContentType = "application/vnd.kubernetes.protobuf"
	// Disable compression for self-communication, since we are going to be
	// on a fast local network
	genericConfig.LoopbackClientConfig.DisableCompression = true

	zeroClientConfig := genericConfig.LoopbackClientConfig
	clientgoExternalClient, err := versioned.NewForConfig(zeroClientConfig)
	if err != nil {
		lastErr = fmt.Errorf("failed to create real external clientset: %w", err)
		return
	}
	versionedInformers = informers.NewSharedInformerFactory(clientgoExternalClient, 10*time.Minute)

	if utilfeature.DefaultFeatureGate.Enabled(genericfeatures.AggregatedDiscoveryEndpoint) {
		genericConfig.AggregatedDiscoveryGroupManager = aggregated.NewResourceManager()
	}

	return
}

// completedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type completedServerRunOptions struct {
	*options.ServerRunOptions
}

// Validate checks completedServerRunOptions and return a slice of found errs.
func (o completedServerRunOptions) Validate() []error {
	errs := []error{}
	errs = append(errs, o.RecommendedOptions.Validate()...)
	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.Metrics.Validate()...)
	// errs = append(errs, o.CloudOptions.Validate()...)

	return errs
}

// Complete fills in fields required to have valid data.
func Complete(o *options.ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions
	// set defaults
	if err := o.GenericServerRunOptions.DefaultAdvertiseAddress(o.RecommendedOptions.SecureServing.SecureServingOptions); err != nil {
		return options, err
	}

	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts(
		o.GenericServerRunOptions.AdvertiseAddress.String(),
		[]string{"zero.io"},
		[]net.IP{netutils.ParseIPSloppy("127.0.0.1")}); err != nil {
		return options, fmt.Errorf("error creating self-signed certificates: %w", err)
	}

	//nolint: nestif
	if len(o.GenericServerRunOptions.ExternalHost) == 0 {
		if len(o.GenericServerRunOptions.AdvertiseAddress) > 0 {
			o.GenericServerRunOptions.ExternalHost = o.GenericServerRunOptions.AdvertiseAddress.String()
		} else {
			if hostname, err := os.Hostname(); err == nil {
				o.GenericServerRunOptions.ExternalHost = hostname
			} else {
				return options, fmt.Errorf("error finding host name: %w", err)
			}
		}
		klog.Infof("external host was not specified, using %v", o.GenericServerRunOptions.ExternalHost)
	}

	if o.RecommendedOptions.Etcd.EnableWatchCache {
		sizes := storage.DefaultWatchCacheSizes()
		// Ensure that overrides parse correctly.
		userSpecified, err := serveroptions.ParseWatchCacheSizes(o.RecommendedOptions.Etcd.WatchCacheSizes)
		if err != nil {
			return options, err
		}
		for resource, size := range userSpecified {
			sizes[resource] = size
		}
		o.RecommendedOptions.Etcd.WatchCacheSizes, err = serveroptions.WriteWatchCacheSizes(sizes)
		if err != nil {
			return options, err
		}
	}

	o.RecommendedOptions.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(genericfeatures.APIListChunking)

	o.RecommendedOptions.ExtraAdmissionInitializers = func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
		client, err := versioned.NewForConfig(c.LoopbackClientConfig)
		if err != nil {
			return nil, err
		}
		informerFactory := informers.NewSharedInformerFactory(client, c.LoopbackClientConfig.Timeout)
		o.SharedInformerFactory = informerFactory
		return []admission.PluginInitializer{initializer.New(informerFactory, client)}, nil
	}

	options.ServerRunOptions = o
	return options, nil
}

func convertVersion() *kversion.Info {
	zeroVersion := version.Get()
	v, _ := semver.Make(zeroVersion.GitVersion)

	return &kversion.Info{
		Major:        strconv.FormatUint(v.Major, 10),
		Minor:        strconv.FormatUint(v.Minor, 10),
		GitVersion:   zeroVersion.GitVersion,
		GitCommit:    zeroVersion.GitCommit,
		GitTreeState: zeroVersion.GitTreeState,
		BuildDate:    zeroVersion.BuildDate,
		GoVersion:    zeroVersion.GoVersion,
		Compiler:     zeroVersion.Compiler,
		Platform:     zeroVersion.Platform,
	}
}
