// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package storage

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/fieldmanager"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	autoscalingv1 "k8s.io/kubernetes/pkg/apis/autoscaling/v1"
	autoscalingvalidation "k8s.io/kubernetes/pkg/apis/autoscaling/validation"
	"k8s.io/kubernetes/pkg/printers"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"

	printersinternal "github.com/superproj/zero/internal/pkg/printers/internalversion"
	"github.com/superproj/zero/internal/registry/apps/minerset"
	"github.com/superproj/zero/pkg/apis/apps"
)

// ReplicasPathMappings returns the mappings between each group version and a replicas path.
func ReplicasPathMappings() fieldmanager.ResourcePathMappings {
	return replicasPathInMinerSet
}

// maps a group version to the replicas path in a minerset object.
var replicasPathInMinerSet = fieldmanager.ResourcePathMappings{
	schema.GroupVersion{Group: "apps.zero.io", Version: "v1beta1"}.String(): fieldpath.MakePathOrDie("spec", "replicas"),
}

// REST implements a RESTStorage for minersets.
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against minersets.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, *StatusREST, *ScaleREST, error) {
	store := &genericregistry.Store{
		NewFunc:       func() runtime.Object { return &apps.MinerSet{} },
		NewListFunc:   func() runtime.Object { return &apps.MinerSetList{} },
		PredicateFunc: minerset.Matcher,
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*apps.MinerSet).Name, nil
		},
		DefaultQualifiedResource: apps.Resource("minersets"),

		CreateStrategy:      minerset.Strategy,
		UpdateStrategy:      minerset.Strategy,
		DeleteStrategy:      minerset.Strategy,
		ResetFieldsStrategy: minerset.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: minerset.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, nil, nil, err
	}

	// Subresources use the same store and creation strategy, which only
	// allows empty subs. Updates to an existing subresource are handled by
	// dedicated strategies.
	statusStore := *store
	statusStore.UpdateStrategy = minerset.StatusStrategy
	statusStore.ResetFieldsStrategy = minerset.StatusStrategy
	return &REST{store}, &StatusREST{store: &statusStore}, &ScaleREST{store: store}, nil
}

// Implement ShortNamesProvider.
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"ms"}
}

var _ rest.CategoriesProvider = &REST{}

// Categories implements the CategoriesProvider interface. Returns a list of categories a resource is part of.
func (r *REST) Categories() []string {
	return []string{"all"}
}

// StatusREST implements the REST endpoint for changing the status of a minerset.
type StatusREST struct {
	store *genericregistry.Store
}

// New creates a new MinerSet object.
func (r *StatusREST) New() runtime.Object {
	return &apps.MinerSet{}
}

// Destroy cleans up resources on shutdown.
func (r *StatusREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *metav1.UpdateOptions,
) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// GetResetFields implements rest.ResetFieldsStrategy.
func (r *StatusREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return r.store.GetResetFields()
}

func (r *StatusREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.store.ConvertToTable(ctx, object, tableOptions)
}

// ScaleREST implements a Scale for MinerSet.
type ScaleREST struct {
	store *genericregistry.Store
}

// ScaleREST implements Patcher.
var (
	_ = rest.Patcher(&ScaleREST{})
	_ = rest.GroupVersionKindProvider(&ScaleREST{})
)

// GroupVersionKind returns GroupVersionKind for MinerSet Scale object.
func (r *ScaleREST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return autoscalingv1.SchemeGroupVersion.WithKind("Scale")
}

// New creates a new Scale object.
func (r *ScaleREST) New() runtime.Object {
	return &autoscaling.Scale{}
}

// Destroy cleans up resources on shutdown.
func (r *ScaleREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

// Get retrieves object from Scale storage.
func (r *ScaleREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := r.store.Get(ctx, name, options)
	if err != nil {
		return nil, errors.NewNotFound(apps.Resource("minerset/scale"), name)
	}
	minerset := obj.(*apps.MinerSet)
	scale, err := scaleFromMinerSet(minerset)
	if err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("%v", err))
	}
	return scale, nil
}

// Update alters scale subset of MinerSet object.
func (r *ScaleREST) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *metav1.UpdateOptions,
) (runtime.Object, bool, error) {
	obj, _, err := r.store.Update(
		ctx,
		name,
		&scaleUpdatedObjectInfo{name, objInfo},
		toScaleCreateValidation(createValidation),
		toScaleUpdateValidation(updateValidation),
		false,
		options,
	)
	if err != nil {
		return nil, false, err
	}

	minerset := obj.(*apps.MinerSet)
	newScale, err := scaleFromMinerSet(minerset)
	if err != nil {
		return nil, false, errors.NewBadRequest(fmt.Sprintf("%v", err))
	}
	return newScale, false, nil
}

func (r *ScaleREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.store.ConvertToTable(ctx, object, tableOptions)
}

func toScaleCreateValidation(f rest.ValidateObjectFunc) rest.ValidateObjectFunc {
	return func(ctx context.Context, obj runtime.Object) error {
		scale, err := scaleFromMinerSet(obj.(*apps.MinerSet))
		if err != nil {
			return err
		}
		return f(ctx, scale)
	}
}

func toScaleUpdateValidation(f rest.ValidateObjectUpdateFunc) rest.ValidateObjectUpdateFunc {
	return func(ctx context.Context, obj, old runtime.Object) error {
		newScale, err := scaleFromMinerSet(obj.(*apps.MinerSet))
		if err != nil {
			return err
		}
		oldScale, err := scaleFromMinerSet(old.(*apps.MinerSet))
		if err != nil {
			return err
		}
		return f(ctx, newScale, oldScale)
	}
}

// scaleFromMinerSet returns a scale subresource for a minerset.
func scaleFromMinerSet(minerset *apps.MinerSet) (*autoscaling.Scale, error) {
	selector, err := metav1.LabelSelectorAsSelector(&minerset.Spec.Selector)
	if err != nil {
		return nil, err
	}

	return &autoscaling.Scale{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Scale",
			APIVersion: "autoscaling/v1",
		},
		// TODO: Create a variant of ObjectMeta type that only contains the fields below.
		ObjectMeta: metav1.ObjectMeta{
			Name:              minerset.Name,
			Namespace:         minerset.Namespace,
			UID:               minerset.UID,
			ResourceVersion:   minerset.ResourceVersion,
			CreationTimestamp: minerset.CreationTimestamp,
		},
		Spec: autoscaling.ScaleSpec{
			Replicas: *minerset.Spec.Replicas,
		},
		Status: autoscaling.ScaleStatus{
			Replicas: minerset.Status.Replicas,
			Selector: selector.String(),
		},
	}, nil
}

// scaleUpdatedObjectInfo transforms existing minerset -> existing scale -> new scale -> new minerset.
type scaleUpdatedObjectInfo struct {
	name       string
	reqObjInfo rest.UpdatedObjectInfo
}

func (i *scaleUpdatedObjectInfo) Preconditions() *metav1.Preconditions {
	return i.reqObjInfo.Preconditions()
}

func (i *scaleUpdatedObjectInfo) UpdatedObject(ctx context.Context, oldObj runtime.Object) (runtime.Object, error) {
	minerset, ok := oldObj.DeepCopyObject().(*apps.MinerSet)
	if !ok {
		return nil, errors.NewBadRequest(fmt.Sprintf("expected existing object type to be MinerSet, got %T", minerset))
	}
	// if zero-value, the existing object does not exist
	if len(minerset.ResourceVersion) == 0 {
		return nil, errors.NewNotFound(apps.Resource("minerset/scale"), i.name)
	}

	groupVersion := schema.GroupVersion{Group: "apps.zero.io", Version: "v1beta1"}
	if requestInfo, found := genericapirequest.RequestInfoFrom(ctx); found {
		requestGroupVersion := schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}
		if _, ok := replicasPathInMinerSet[requestGroupVersion.String()]; ok {
			groupVersion = requestGroupVersion
		} else {
			klog.Fatalf("Unrecognized group/version in request info %q", requestGroupVersion.String())
		}
	}

	managedFieldsHandler := fieldmanager.NewScaleHandler(
		minerset.ManagedFields,
		groupVersion,
		replicasPathInMinerSet,
	)

	// minerset -> old scale
	oldScale, err := scaleFromMinerSet(minerset)
	if err != nil {
		return nil, err
	}

	scaleManagedFields, err := managedFieldsHandler.ToSubresource()
	if err != nil {
		return nil, err
	}
	oldScale.ManagedFields = scaleManagedFields

	// old scale -> new scale
	newScaleObj, err := i.reqObjInfo.UpdatedObject(ctx, oldScale)
	if err != nil {
		return nil, err
	}
	if newScaleObj == nil {
		return nil, errors.NewBadRequest("nil update passed to Scale")
	}
	scale, ok := newScaleObj.(*autoscaling.Scale)
	if !ok {
		return nil, errors.NewBadRequest(fmt.Sprintf("expected input object type to be Scale, but %T", newScaleObj))
	}

	// validate
	if errs := autoscalingvalidation.ValidateScale(scale); len(errs) > 0 {
		return nil, errors.NewInvalid(autoscaling.Kind("Scale"), minerset.Name, errs)
	}

	// validate precondition if specified (resourceVersion matching is handled by storage)
	if len(scale.UID) > 0 && scale.UID != minerset.UID {
		return nil, errors.NewConflict(
			apps.Resource("minerset/scale"),
			minerset.Name,
			fmt.Errorf("precondition failed: UID in precondition: %v, UID in object meta: %v", scale.UID, minerset.UID),
		)
	}

	// move replicas/resourceVersion fields to object and return
	minerset.Spec.Replicas = &scale.Spec.Replicas
	minerset.ResourceVersion = scale.ResourceVersion

	updatedEntries, err := managedFieldsHandler.ToParent(scale.ManagedFields)
	if err != nil {
		return nil, err
	}
	minerset.ManagedFields = updatedEntries

	return minerset, nil
}
