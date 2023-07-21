// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package collections

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/superproj/zero/internal/pkg/util/conditions"
	coreutil "github.com/superproj/zero/internal/pkg/util/core"
	appsv1beta1 "github.com/superproj/zero/pkg/apis/apps/v1beta1"
)

// Func is the functon definition for a filter.
type Func func(miner *appsv1beta1.Miner) bool

// And returns a filter that returns true if all of the given filters returns true.
func And(filters ...Func) Func {
	return func(miner *appsv1beta1.Miner) bool {
		for _, f := range filters {
			if !f(miner) {
				return false
			}
		}
		return true
	}
}

// Or returns a filter that returns true if any of the given filters returns true.
func Or(filters ...Func) Func {
	return func(miner *appsv1beta1.Miner) bool {
		for _, f := range filters {
			if f(miner) {
				return true
			}
		}
		return false
	}
}

// Not returns a filter that returns the opposite of the given filter.
func Not(mf Func) Func {
	return func(miner *appsv1beta1.Miner) bool {
		return !mf(miner)
	}
}

// HasControllerRef is a filter that returns true if the miner has a controller ref.
func HasControllerRef(miner *appsv1beta1.Miner) bool {
	if miner == nil {
		return false
	}
	return metav1.GetControllerOf(miner) != nil
}

/*
// InFailureDomains returns a filter to find all miners
// in any of the given failure domains.
func InFailureDomains(failureDomains ...*string) Func {
	return func(miner *appsv1beta1.Miner) bool {
		if miner == nil {
			return false
		}
		for i := range failureDomains {
			fd := failureDomains[i]
			if fd == nil {
				if fd == miner.Spec.FailureDomain {
					return true
				}
				continue
			}
			if miner.Spec.FailureDomain == nil {
				continue
			}
			if *fd == *miner.Spec.FailureDomain {
				return true
			}
		}
		return false
	}
}
*/

// OwnedMiners returns a filter to find all miners owned by specified owner.
// Usage: GetFilteredMinersForCluster(ctx, client, cluster, OwnedMiners(controlPlane)).
func OwnedMiners(owner client.Object) func(miner *appsv1beta1.Miner) bool {
	return func(miner *appsv1beta1.Miner) bool {
		if miner == nil {
			return false
		}
		return coreutil.IsOwnedByObject(miner, owner)
	}
}

/*
// ControlPlaneMiners returns a filter to find all control plane miners for a cluster, regardless of ownership.
// Usage: GetFilteredMinersForCluster(ctx, client, cluster, ControlPlaneMiners(cluster.Name)).
func ControlPlaneMiners(clusterName string) func(miner *appsv1beta1.Miner) bool {
	selector := ControlPlaneSelectorForCluster(clusterName)
	return func(miner *appsv1beta1.Miner) bool {
		if miner == nil {
			return false
		}
		return selector.Matches(labels.Set(miner.Labels))
	}
}

// AdoptableControlPlaneMiners returns a filter to find all un-controlled control plane miners.
// Usage: GetFilteredMinersForCluster(ctx, client, cluster, AdoptableControlPlaneMiners(cluster.Name, controlPlane)).
func AdoptableControlPlaneMiners(clusterName string) func(miner *appsv1beta1.Miner) bool {
	return And(
		ControlPlaneMiners(clusterName),
		Not(HasControllerRef),
	)
}
*/

// ActiveMiners returns a filter to find all active miners.
// Usage: GetFilteredMinersForCluster(ctx, client, cluster, ActiveMiners).
func ActiveMiners(miner *appsv1beta1.Miner) bool {
	if miner == nil {
		return false
	}
	return miner.DeletionTimestamp.IsZero()
}

// HasDeletionTimestamp returns a filter to find all miners that have a deletion timestamp.
func HasDeletionTimestamp(miner *appsv1beta1.Miner) bool {
	if miner == nil {
		return false
	}
	return !miner.DeletionTimestamp.IsZero()
}

// HasUnhealthyCondition returns a filter to find all miners that have a MinerHealthCheckSucceeded condition set to False,
// indicating a problem was detected on the miner, and the MinerOwnerRemediated condition set, indicating that KCP is
// responsible of performing remediation as owner of the miner.
func HasUnhealthyCondition(miner *appsv1beta1.Miner) bool {
	if miner == nil {
		return false
	}
	return conditions.IsFalse(miner, appsv1beta1.MinerHealthCheckSucceededCondition) && conditions.IsFalse(miner, appsv1beta1.MinerOwnerRemediatedCondition)
}

// IsReady returns a filter to find all miners with the ReadyCondition equals to True.
func IsReady() Func {
	return func(miner *appsv1beta1.Miner) bool {
		if miner == nil {
			return false
		}
		return conditions.IsTrue(miner, appsv1beta1.ReadyCondition)
	}
}

// ShouldRolloutAfter returns a filter to find all miners where
// CreationTimestamp < rolloutAfter < reconciliationTIme.
func ShouldRolloutAfter(reconciliationTime, rolloutAfter *metav1.Time) Func {
	return func(miner *appsv1beta1.Miner) bool {
		if miner == nil {
			return false
		}
		return miner.CreationTimestamp.Before(rolloutAfter) && rolloutAfter.Before(reconciliationTime)
	}
}

// HasAnnotationKey returns a filter to find all miners that have the
// specified Annotation key present.
func HasAnnotationKey(key string) Func {
	return func(miner *appsv1beta1.Miner) bool {
		if miner == nil || miner.Annotations == nil {
			return false
		}
		if _, ok := miner.Annotations[key]; ok {
			return true
		}
		return false
	}
}

/*
// ControlPlaneSelectorForCluster returns the label selector necessary to get control plane miners for a given cluster.
func ControlPlaneSelectorForCluster(clusterName string) labels.Selector {
	must := func(r *labels.Requirement, err error) labels.Requirement {
		if err != nil {
			panic(err)
		}
		return *r
	}
	return labels.NewSelector().Add(
		must(labels.NewRequirement(appsv1beta1.ClusterLabelName, selection.Equals, []string{clusterName})),
		must(labels.NewRequirement(appsv1beta1.MinerControlPlaneLabelName, selection.Exists, []string{})),
	)
}

// MatchesKubernetesVersion returns a filter to find all miners that match a given Kubernetes version.
func MatchesKubernetesVersion(kubernetesVersion string) Func {
	return func(miner *appsv1beta1.Miner) bool {
		if miner == nil {
			return false
		}
		if miner.Spec.Version == nil {
			return false
		}
		return *miner.Spec.Version == kubernetesVersion
	}
}

// WithVersion returns a filter to find miner that have a non empty and valid version.
func WithVersion() Func {
	return func(miner *appsv1beta1.Miner) bool {
		if miner == nil {
			return false
		}
		if miner.Spec.Version == nil {
			return false
		}
		if _, err := semver.ParseTolerant(*miner.Spec.Version); err != nil {
			return false
		}
		return true
	}
}

// HealthyAPIServer returns a filter to find all miners that have a MinerAPIServerPodHealthyCondition
// set to true.
func HealthyAPIServer() Func {
	return func(miner *appsv1beta1.Miner) bool {
		if miner == nil {
			return false
		}
		return conditions.IsTrue(miner, controlplanev1.MinerAPIServerPodHealthyCondition)
	}
}
*/
