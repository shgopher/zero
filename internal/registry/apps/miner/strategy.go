// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//nolint:gocritic
package miner

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

// minerStrategy implements behavior for Miner objects.
type minerStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Miner
// objects via the REST API.
var Strategy = minerStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped is true for miners.
func (minerStrategy) NamespaceScoped() bool {
	return true
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (minerStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"apps.zero.io/v1beta1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
			fieldpath.MakePathOrDie("status"),
		),
	}

	return fields
}

// AllowCreateOnUpdate is false for miners.
func (minerStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForCreate clears fields that are not allowed to be set by end users
// on creation.
func (minerStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	miner := obj.(*apps.Miner)
	miner.Generation = 1

	// Be explicit that users cannot create pre-provisioned miners.
	miner.Status = apps.MinerStatus{}
	miner.Status.Conditions = []apps.Condition{}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (minerStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newMiner := obj.(*apps.Miner)
	oldMiner := old.(*apps.Miner)
	// If you want to change Spec through subresources, you can uncomment following assignments
	// newMiner.Spec = oldMiner.Spec
	newMiner.Status = oldMiner.Status
}

// Validate validates a new miner.
func (minerStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	miner := obj.(*apps.Miner)
	return validation.ValidateMiner(miner)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (minerStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (minerStrategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end user.
func (minerStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateMinerUpdate(obj.(*apps.Miner), old.(*apps.Miner))
}

// WarningsOnUpdate returns warnings for the given update.
func (minerStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// If AllowUnconditionalUpdate() is true and the object specified by
// the user does not have a resource version, then generic Update()
// populates it with the latest version. Else, it checks that the
// version specified by the user matches the version of latest etcd
// object.
func (minerStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// Storage strategy for the Status subresource.
type minerStatusStrategy struct {
	minerStrategy
}

// StatusStrategy is the default logic invoked when updating object status.
var StatusStrategy = minerStatusStrategy{Strategy}

// Make sure we correctly implement the interface.
var _ = rest.GarbageCollectionDeleteStrategy(Strategy)

// DefaultGarbageCollectionPolicy returns DeleteDependents for all currently served versions.
func (minerStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.DeleteDependents
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (minerStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"apps.zero.io/v1beta1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
			fieldpath.MakePathOrDie("status", "conditions"),
		),
	}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update of status.
func (minerStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newMiner := obj.(*apps.Miner)
	oldMiner := old.(*apps.Miner)

	// Updating /status should not modify spec
	newMiner.Spec = oldMiner.Spec
	newMiner.DeletionTimestamp = nil

	// don't allow the miners/status endpoint to touch owner references since old kubelets corrupt them in a way
	// that breaks garbage collection
	newMiner.OwnerReferences = oldMiner.OwnerReferences
}

// ValidateUpdate is the default update validation for an end user updating status.
func (minerStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateMinerStatusUpdate(obj.(*apps.Miner), old.(*apps.Miner))
}

// WarningsOnUpdate returns warnings for the given update.
func (minerStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (minerStatusStrategy) Canonicalize(obj runtime.Object) {
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	m, ok := obj.(*apps.Miner)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a miner")
	}
	return labels.Set(m.Labels), ToSelectableFields(m), nil
}

// ToSelectableFields returns a field set that can be used for filter selection.
func ToSelectableFields(obj *apps.Miner) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
	minerSpecificFieldsSet := fields.Set{
		"spec.minerType": obj.Spec.MinerType,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, minerSpecificFieldsSet)
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
