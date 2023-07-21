// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	componentbaseconfigv1alpha1 "k8s.io/component-base/config/v1alpha1"

	genericconfigv1beta1 "github.com/superproj/zero/pkg/config/v1beta1"
)

const (
	// ZeroControllerManagerDefaultLockObjectNamespace defines default zero controller manager lock object namespace ("kube-system").
	ZeroControllerManagerDefaultLockObjectNamespace string = metav1.NamespaceSystem

	// ZeroControllerManagerDefaultLockObjectName defines default zero controller manager lock object name ("zero-controller-manager").
	ZeroControllerManagerDefaultLockObjectName = "zero-controller-manager"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZeroControllerManagerConfiguration contains elements describing zero-controller manager.
type ZeroControllerManagerConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// FeatureGates is a map of feature names to bools that enable or disable alpha/experimental features.
	FeatureGates map[string]bool `json:"featureGates,omitempty"`

	// Generic holds configuration for a generic controller-manager
	Generic GenericControllerManagerConfiguration `json:"generic,omitempty"`

	// GarbageCollectorControllerConfiguration holds configuration for
	// GarbageCollectorController related features.
	GarbageCollectorController GarbageCollectorControllerConfiguration `json:"garbageCollectorController,omitempty"`

	// ChainControllerConfiguration holds configuration for ChainController related features.
	ChainController ChainControllerConfiguration `json:"chainController,omitempty"`
}

// GenericControllerManagerConfiguration holds configuration for a generic controller-manager.
type GenericControllerManagerConfiguration struct {
	// MySQL defines the configuration of mysql client.
	MySQL genericconfigv1beta1.MySQLConfiguration `json:"mysql,omitempty"`

	// leaderElection defines the configuration of leader election client.
	LeaderElection componentbaseconfigv1alpha1.LeaderElectionConfiguration `json:"leaderElection,omitempty"`

	// Namespace that the controller watches to reconcile zero-apiserver objects.
	Namespace string `json:"namespace,omitempty"`

	// bindAddress is the IP address for the proxy server to serve on (set to 0.0.0.0
	// for all interfaces)
	BindAddress string `json:"bindAddress,omitempty"`

	// MetricsBindAddress is the IP address and port for the metrics server to serve on,
	// defaulting to 127.0.0.1:20249 (set to 0.0.0.0 for all interfaces)
	MetricsBindAddress string `json:"metricsBindAddress,omitempty"`

	// HealthzBindAddress is the IP address and port for the health check server to serve on,
	// defaulting to 0.0.0.0:20250
	HealthzBindAddress string `json:"healthzBindAddress,omitempty"`

	// Parallelism defines the amount of parallelism to process miners. Must be greater than 0. Defaults to 16
	Parallelism int32 `json:"parallelism,omitempty"`

	// SyncPeriod determines the minimum frequency at which watched resources are
	// reconciled. A lower period will correct entropy more quickly, but reduce
	// responsiveness to change if there are many watched resources. Change this
	// value only if you know what you are doing. Defaults to 10 hours if unset.
	SyncPeriod metav1.Duration `json:"syncPeriod,omitempty"`

	// Label value that the controller watches to reconcile cloud miner objects
	WatchFilterValue string `json:"watchFilterValue,omitempty"`
}

type ChainControllerConfiguration struct {
	// Image specify the blockchain node image.
	Image string `json:"image,omitempty"`
}

// GroupResource describes an group resource.
type GroupResource struct {
	// group is the group portion of the GroupResource.
	Group string `json:"group,omitempty"`
	// resource is the resource portion of the GroupResource.
	Resource string `json:"resource,omitempty"`
}

// GarbageCollectorControllerConfiguration contains elements describing GarbageCollectorController.
type GarbageCollectorControllerConfiguration struct {
	// enables the generic garbage collector. MUST be synced with the
	// corresponding flag of the kube-apiserver. WARNING: the generic garbage
	// collector is an alpha feature.
	EnableGarbageCollector *bool `json:"enableGarbageCollector,omitempty"`
	// concurrentGCSyncs is the number of garbage collector workers that are
	// allowed to sync concurrently.
	ConcurrentGCSyncs int32 `json:"concurrentGCSyncs,omitempty"`
	// gcIgnoredResources is the list of GroupResources that garbage collection should ignore.
	GCIgnoredResources []GroupResource `json:"gcIgnoredResources,omitempty"`
}
