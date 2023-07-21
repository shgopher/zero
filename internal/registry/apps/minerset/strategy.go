// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//nolint:gocritic
package minerset

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"

	"github.com/superproj/zero/pkg/apis/apps"
	"github.com/superproj/zero/pkg/apis/apps/validation"
)

// minerSetStrategy implements behavior for MinerSet objects.
type minerSetStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating MinerSet
// objects via the REST API.
var Strategy = minerSetStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped is true for minersets.
func (minerSetStrategy) NamespaceScoped() bool {
	return true
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (minerSetStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"apps.zero.io/v1beta1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
			fieldpath.MakePathOrDie("status"),
		),
	}

	return fields
}

// AllowCreateOnUpdate is false for minersets.
func (minerSetStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForCreate clears fields that are not allowed to be set by end users
// on creation.
func (minerSetStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	minerset := obj.(*apps.MinerSet)
	minerset.Generation = 1

	// Be explicit that users cannot create pre-provisioned minersets.
	minerset.Status = apps.MinerSetStatus{}
	minerset.Status.Conditions = []apps.Condition{}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (minerSetStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newMinerSet := obj.(*apps.MinerSet)
	oldMinerSet := old.(*apps.MinerSet)
	// If you want to change Spec through subresources, you can uncomment following assignments
	// newMinerSet.Spec = oldMiner.Spec
	newMinerSet.Status = oldMinerSet.Status
}

// Validate validates a new minerset.
func (minerSetStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	minerset := obj.(*apps.MinerSet)
	return validation.ValidateMinerSet(minerset)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (minerSetStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (minerSetStrategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end user.
func (minerSetStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateMinerSetUpdate(obj.(*apps.MinerSet), old.(*apps.MinerSet))
}

// WarningsOnUpdate returns warnings for the given update.
func (minerSetStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// If AllowUnconditionalUpdate() is true and the object specified by
// the user does not have a resource version, then generic Update()
// populates it with the latest version. Else, it checks that the
// version specified by the user matches the version of latest etcd
// object.
func (minerSetStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// Storage strategy for the Status subresource.
type minerSetStatusStrategy struct {
	minerSetStrategy
}

// StatusStrategy is the default logic invoked when updating object status.
var StatusStrategy = minerSetStatusStrategy{Strategy}

// Make sure we correctly implement the interface.
var _ = rest.GarbageCollectionDeleteStrategy(Strategy)

// DefaultGarbageCollectionPolicy returns DeleteDependents for all currently served versions.
func (minerSetStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.DeleteDependents
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (minerSetStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"apps.zero.io/v1beta1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
			fieldpath.MakePathOrDie("status", "conditions"),
		),
	}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update of status.
func (minerSetStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newMinerSet := obj.(*apps.MinerSet)
	oldMinerSet := old.(*apps.MinerSet)

	// Updating /status should not modify spec
	newMinerSet.Spec = oldMinerSet.Spec
	newMinerSet.DeletionTimestamp = nil

	// don't allow the minersets/status endpoint to touch owner references since old kubelets corrupt them in a way
	// that breaks garbage collection
	newMinerSet.OwnerReferences = oldMinerSet.OwnerReferences
}

// ValidateUpdate is the default update validation for an end user updating status.
func (minerSetStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateMinerSetStatusUpdate(obj.(*apps.MinerSet), old.(*apps.MinerSet))
}

// WarningsOnUpdate returns warnings for the given update.
func (minerSetStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (minerSetStatusStrategy) Canonicalize(obj runtime.Object) {
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	m, ok := obj.(*apps.MinerSet)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a minerset")
	}
	return labels.Set(m.Labels), ToSelectableFields(m), nil
}

// ToSelectableFields returns a field set that can be used for filter selection.
func ToSelectableFields(obj *apps.MinerSet) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
	minerSetSpecificFieldsSet := fields.Set{
		// "spec.type":    obj.Spec.Type, TODO ?
		// "spec.address": obj.Spec.Address,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, minerSetSpecificFieldsSet)
}

// Matcher is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func Matcher(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}
