// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package apiserver

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	apiserverfeatures "k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/routes"

	"github.com/superproj/zero/internal/apiserver/storage"
	appsrest "github.com/superproj/zero/internal/registry/apps/rest"
	coordinationrest "github.com/superproj/zero/internal/registry/coordination/rest"
	corerest "github.com/superproj/zero/internal/registry/core/rest"
	appsv1beta1 "github.com/superproj/zero/pkg/apis/apps/v1beta1"
	coordinationv1 "github.com/superproj/zero/pkg/apis/coordination/v1"
	apiv1 "github.com/superproj/zero/pkg/apis/core/v1"
	informers "github.com/superproj/zero/pkg/generated/informers/externalversions"
)

// ExtraConfig holds custom apiserver config.
type ExtraConfig struct {
	// Place you custom config here.
	APIResourceConfigSource serverstorage.APIResourceConfigSource
	StorageFactory          serverstorage.StorageFactory
	EventTTL                time.Duration
	EnableLogsSupport       bool
	VersionedInformers      informers.SharedInformerFactory
}

// Config defines the config for the apiserver.
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

// ZeroServer contains state for a Kubernetes cluster master/api server.
type ZeroServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	cfg := completedConfig{c.GenericConfig.Complete(), &c.ExtraConfig}

	return CompletedConfig{&cfg}
}

// New returns a new instance of ZeroServer from the given config.
// Certain config fields will be set to a default value if unset.
func (c completedConfig) New() (*ZeroServer, error) {
	genericServer, err := c.GenericConfig.New("zero-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	if c.ExtraConfig.EnableLogsSupport {
		routes.Logs{}.Install(genericServer.Handler.GoRestfulContainer)
	}

	s := &ZeroServer{
		GenericAPIServer: genericServer,
	}

	// install zero legacy rest storage
	if err := s.InstallLegacyAPI(&c, c.GenericConfig.RESTOptionsGetter); err != nil {
		return nil, err
	}

	// The order here is preserved in discovery.
	// If resources with identical names exist in more than one of these groups (e.g. "miners.apps"" and "miners.core"),
	// the order of this list determines which group an unqualified resource name (e.g. "deployments") should prefer.
	// This priority order is used for local discovery, but it ends up aggregated in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go
	// with specific priorities.
	// TODO: describe the priority all the way down in the RESTStorageProviders and plumb it back through the various discovery
	// handlers that we have.
	restStorageProviders := []storage.RESTStorageProvider{
		// &admissionrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		appsrest.RESTStorageProvider{},
		coordinationrest.RESTStorageProvider{},
	}
	if err := s.InstallAPIs(c.ExtraConfig.APIResourceConfigSource, c.GenericConfig.RESTOptionsGetter, restStorageProviders...); err != nil {
		return nil, err
	}

	// TODO: copy from kube-apiserver
	s.GenericAPIServer.AddPostStartHookOrDie(
		"start-zero-server-informers",
		func(context genericapiserver.PostStartHookContext) error {
			// remove dependence with kube-apiserver
			c.ExtraConfig.VersionedInformers.Start(context.StopCh)
			return nil
		},
	)

	if utilfeature.DefaultFeatureGate.Enabled(apiserverfeatures.APIServerIdentity) {
		// put some post start hook here
	}

	return s, nil
}

// InstallLegacyAPI will install the legacy APIs for the restStorageProviders if they are enabled.
func (s *ZeroServer) InstallLegacyAPI(c *completedConfig, restOptionsGetter generic.RESTOptionsGetter) error {
	legacyRESTStorageProvider := corerest.LegacyRESTStorageProvider{
		EventTTL: c.ExtraConfig.EventTTL,
	}
	apiGroupInfo, err := legacyRESTStorageProvider.NewLegacyRESTStorage(restOptionsGetter)
	if err != nil {
		return fmt.Errorf("error building core storage: %w", err)
	}
	if len(apiGroupInfo.VersionedResourcesStorageMap) == 0 { // if all core storage is disabled, return.
		return nil
	}

	if err := s.GenericAPIServer.InstallLegacyAPIGroup(genericapiserver.DefaultLegacyAPIPrefix, &apiGroupInfo); err != nil {
		return fmt.Errorf("error in registering group versions: %w", err)
	}
	return nil
}

// ZeroServer will install the APIs for the restStorageProviders if they are enabled.
func (s *ZeroServer) InstallAPIs(
	apiResourceConfigSource serverstorage.APIResourceConfigSource,
	restOptionsGetter generic.RESTOptionsGetter,
	restStorageProviders ...storage.RESTStorageProvider,
) error {
	apiGroupsInfo := []*genericapiserver.APIGroupInfo{}

	// used later in the loop to filter the served resource by those that have expired.
	resourceExpirationEvaluator, err := genericapiserver.NewResourceExpirationEvaluator(*s.GenericAPIServer.Version)
	if err != nil {
		return err
	}

	for _, restStorageBuilder := range restStorageProviders {
		groupName := restStorageBuilder.GroupName()
		if !apiResourceConfigSource.AnyResourceForGroupEnabled(groupName) {
			klog.V(1).InfoS("Skipping disabled API group", "groupName", groupName)
			continue
		}
		apiGroupInfo, err := restStorageBuilder.NewRESTStorage(apiResourceConfigSource, restOptionsGetter)
		if err != nil {
			return fmt.Errorf("problem initializing API group %q: %w", groupName, err)
		}

		if len(apiGroupInfo.VersionedResourcesStorageMap) == 0 {
			// If we have no storage for any resource configured, this API group is effectively disabled.
			// This can happen when an entire API group, version, or development-stage (alpha, beta, GA) is disabled.
			klog.V(1).InfoS("API group is not enabled, skipping.", "groupName", groupName)
			continue
		}

		// Remove resources that serving kinds that are removed.
		// We do this here so that we don't accidentally serve versions without resources or openapi information that for kinds we don't serve.
		// This is a spot above the construction of individual storage handlers so that no sig accidentally forgets to check.
		resourceExpirationEvaluator.RemoveDeletedKinds(groupName, apiGroupInfo.Scheme, apiGroupInfo.VersionedResourcesStorageMap)
		if len(apiGroupInfo.VersionedResourcesStorageMap) == 0 {
			klog.V(1).InfoS("Removing API group because it is time to stop serving it because it has no versions per APILifecycle.", "groupName", groupName)
			continue
		}

		klog.V(1).InfoS("Enabling API group", "groupName", groupName)

		if postHookProvider, ok := restStorageBuilder.(genericapiserver.PostStartHookProvider); ok {
			name, hook, err := postHookProvider.PostStartHook()
			if err != nil {
				klog.Fatalf("Error building PostStartHook: %v", err)
			}
			s.GenericAPIServer.AddPostStartHookOrDie(name, hook)
		}

		apiGroupsInfo = append(apiGroupsInfo, &apiGroupInfo)
	}

	if err := s.GenericAPIServer.InstallAPIGroups(apiGroupsInfo...); err != nil {
		return fmt.Errorf("error in registering group versions: %w", err)
	}
	return nil
}

var (
	// stableAPIGroupVersionsEnabledByDefault is a list of our stable versions.
	stableAPIGroupVersionsEnabledByDefault = []schema.GroupVersion{
		apiv1.SchemeGroupVersion,
		appsv1beta1.SchemeGroupVersion,
		coordinationv1.SchemeGroupVersion,
	}

	// legacyBetaEnabledByDefaultResources is the list of beta resources we enable.  You may only add to this list
	// if your resource is already enabled by default in a beta level we still serve AND there is no stable API for it.
	// see https://github.com/kubernetes/enhancements/tree/master/keps/sig-architecture/3136-beta-apis-off-by-default
	// for more details.
	legacyBetaEnabledByDefaultResources = []schema.GroupVersionResource{}

	// betaAPIGroupVersionsDisabledByDefault is for all future beta groupVersions.
	betaAPIGroupVersionsDisabledByDefault = []schema.GroupVersion{}
)

// DefaultAPIResourceConfigSource returns which groupVersion enabled and its
// resources enabled/disabled.
func DefaultAPIResourceConfigSource() *serverstorage.ResourceConfig {
	ret := serverstorage.NewResourceConfig()
	// NOTE: GroupVersions listed here will be enabled by default. Don't put alpha versions in the list.
	ret.EnableVersions(stableAPIGroupVersionsEnabledByDefault...)

	// disable alpha and beta versions explicitly so we have a full list of what's possible to serve
	ret.DisableVersions(betaAPIGroupVersionsDisabledByDefault...)

	// enable the legacy beta resources that were present before stopped serving new beta APIs by default.
	ret.EnableResources(legacyBetaEnabledByDefaultResources...)

	return ret
}
