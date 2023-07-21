// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//nolint:gocritic
package chain

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

// chainStrategy implements behavior for Chain objects.
type chainStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Chain
// objects via the REST API.
var Strategy = chainStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// Strategy should implement rest.RESTCreateStrategy.
var _ rest.RESTCreateStrategy = Strategy

// Strategy should implement rest.RESTUpdateStrategy.
var _ rest.RESTUpdateStrategy = Strategy

// NamespaceScoped is true for chains.
func (chainStrategy) NamespaceScoped() bool {
	return true
}

// AllowCreateOnUpdate is false for chains.
func (chainStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForCreate clears fields that are not allowed to be set by end users
// on creation.
func (chainStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	chain := obj.(*apps.Chain)
	chain.Generation = 1

	chain.Status = apps.ChainStatus{}
	chain.Status.Conditions = []apps.Condition{}
	dropDisabledFields(chain, nil)
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (chainStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newChain := obj.(*apps.Chain)
	oldChain := old.(*apps.Chain)
	newChain.Status = oldChain.Status
	dropDisabledFields(newChain, oldChain)
}

// Validate validates a new chain.
func (chainStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	chain := obj.(*apps.Chain)
	return validation.ValidateChain(chain)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (chainStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (chainStrategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end user.
func (chainStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateChainUpdate(obj.(*apps.Chain), old.(*apps.Chain))
}

// WarningsOnUpdate returns warnings for the given update.
func (chainStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// If AllowUnconditionalUpdate() is true and the object specified by
// the user does not have a resource version, then generic Update()
// populates it with the latest version. Else, it checks that the
// version specified by the user matches the version of latest etcd
// object.
func (chainStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// Storage strategy for the Status subresource.
type chainStatusStrategy struct {
	chainStrategy
}

// StatusStrategy is the default logic invoked when updating object status.
var StatusStrategy = chainStatusStrategy{Strategy}

// Make sure we correctly implement the interface.
var _ = rest.GarbageCollectionDeleteStrategy(Strategy)

// DefaultGarbageCollectionPolicy returns DeleteDependents for all currently served versions.
func (chainStrategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.DeleteDependents
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (chainStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"apps.cloudchain.io/v1beta1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
			fieldpath.MakePathOrDie("status", "conditions"),
		),
	}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update of status.
func (chainStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newChain := obj.(*apps.Chain)
	oldChain := old.(*apps.Chain)

	// Updating /status should not modify spec
	newChain.Spec = oldChain.Spec
	newChain.DeletionTimestamp = nil

	// don't allow the chains/status endpoint to touch owner references since old kubelets corrupt them in a way
	// that breaks garbage collection
	newChain.OwnerReferences = oldChain.OwnerReferences
}

// ValidateUpdate is the default update validation for an end user updating status.
func (chainStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateChainStatusUpdate(obj.(*apps.Chain), old.(*apps.Chain))
}

// WarningsOnUpdate returns warnings for the given update.
func (chainStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (chainStatusStrategy) Canonicalize(obj runtime.Object) {
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	c, ok := obj.(*apps.Chain)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a chain")
	}
	return labels.Set(c.Labels), ToSelectableFields(c), nil
}

// NameTriggerFunc returns value metadata.namespace of given object.
func NameTriggerFunc(obj runtime.Object) string {
	return obj.(*apps.Chain).ObjectMeta.Name
}

// ToSelectableFields returns a field set that can be used for filter selection.
func ToSelectableFields(obj *apps.Chain) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

// Matcher is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func Matcher(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:       label,
		Field:       field,
		GetAttrs:    GetAttrs,
		IndexFields: []string{"metadata.name"},
	}
}

func dropDisabledFields(chain *apps.Chain, oldChain *apps.Chain) {
}
