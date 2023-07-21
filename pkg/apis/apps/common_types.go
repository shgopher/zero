// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package apps

const (
	// MinerAnnotation is the annotation set on pods identifying the miner the pod belongs to.
	MinerAnnotation = "apps.zero.io/miner"

	// OwnerKindAnnotation is the annotation set on pods identifying the owner kind.
	OwnerKindAnnotation = "apps.zero.io/owner-kind"

	// OwnerNameAnnotation is the annotation set on pods identifying the owner name.
	OwnerNameAnnotation = "apps.zero.io/owner-name"

	// DisableMinerCreate is an annotation that can be used to signal a MinerSet to stop creating new miners.
	// It is utilized in the OnDelete MinerSetStrategy to allow the MinerSet controller to scale down
	// older MinerSets when Miners are deleted and add the new replicas to the latest MinerSet.
	DisableMinerCreate = "apps.zero.io/disable-miner-create"

	// DeleteMinerAnnotation marks control plane and worker nodes that will be given priority for deletion
	// when KCP or a minerset scales down. This annotation is given top priority on all delete policies.
	DeleteMinerAnnotation = "apps.zero.io/delete-miner"
)

// MinerAddressType describes a valid MinerAddress type.
type MinerAddressType string

// Define the MinerAddressType constants.
const (
	MinerHostName    MinerAddressType = "Hostname"
	MinerExternalIP  MinerAddressType = "ExternalIP"
	MinerInternalIP  MinerAddressType = "InternalIP"
	MinerExternalDNS MinerAddressType = "ExternalDNS"
	MinerInternalDNS MinerAddressType = "InternalDNS"
)

// MinerAddress contains information for the miner's address.
type MinerAddress struct {
	// Miner address type, one of Hostname, ExternalIP or InternalIP.
	Type MinerAddressType `json:"type"`

	// The machine address.
	Address string `json:"address"`
}

// MinerAddresses is a slice of MinerAddress items to be used by infrastructure providers.
type MinerAddresses []MinerAddress

// ObjectMeta is metadata that all persisted resources must have, which includes all objects
// users must create. This is a copy of customizable fields from metav1.ObjectMeta.
//
// ObjectMeta is embedded in `Miner.Spec` and `MinerSet.Template`,
// which are not top-level Kubernetes objects. Given that metav1.ObjectMeta has lots of special cases
// and read-only fields which end up in the generated CRD validation, having it as a subset simplifies
// the API and some issues that can impact user experience.
//
// During the [upgrade to controller-tools@v2](https://github.com/kubernetes-sigs/cluster-api/pull/1054)
// for v1alpha2, we noticed a failure would occur running Cluster API test suite against the new CRDs,
// specifically `spec.metadata.creationTimestamp in body must be of type string: "null"`.
// The investigation showed that `controller-tools@v2` behaves differently than its previous version
// when handling types from [metav1](k8s.io/apimachinery/pkg/apis/meta/v1) package.
//
// In more details, we found that embedded (non-top level) types that embedded `metav1.ObjectMeta`
// had validation properties, including for `creationTimestamp` (metav1.Time).
// The `metav1.Time` type specifies a custom json marshaller that, when IsZero() is true, returns `null`
// which breaks validation because the field isn't marked as nullable.
//
// In future versions, controller-tools@v2 might allow overriding the type and validation for embedded
// types. When that happens, this hack should be revisited.
type ObjectMeta struct {
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`
}
