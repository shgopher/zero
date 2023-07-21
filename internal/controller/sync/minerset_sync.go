// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//nolint:dupl
package sync

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"gorm.io/gorm"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	gwmodel "github.com/superproj/zero/internal/gateway/model"
	"github.com/superproj/zero/internal/gateway/store"
	"github.com/superproj/zero/pkg/apis/apps/v1beta1"
)

const minerSetControllerName = "controller-manager.minerSetSync"

// MinerSetSyncReconciler sync a MinerSet object to database.
type MinerSetSyncReconciler struct {
	client client.Client

	Store store.IStore
}

// SetupWithManager sets up the controller with the Manager.
func (r *MinerSetSyncReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.MinerSet{}).
		WithOptions(options).
		Named(minerSetControllerName).
		Complete(r)
	if err != nil {
		return err
	}

	r.client = mgr.GetClient()
	return nil
}

func (r *MinerSetSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the MinerSet instance
	ms := &v1beta1.MinerSet{}
	if err := r.client.Get(ctx, req.NamespacedName, ms); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, r.Store.MinerSets().Delete(ctx, map[string]interface{}{"namespace": req.Namespace, "name": req.Name})
		}
		return ctrl.Result{}, err
	}

	msr, err := r.Store.MinerSets().Get(ctx, map[string]interface{}{"namespace": req.Namespace, "name": req.Name})
	if err != nil {
		// minerset record not exist, create it.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctrl.Result{}, addMinerSet(ctx, r.Store, ms)
		}

		return ctrl.Result{}, err
	}

	// minerset record exist, update it
	originMinerSet := new(gwmodel.MinerSetM)
	*originMinerSet = *msr

	msr = applyToMinerSet(msr, ms)
	if !reflect.DeepEqual(msr, originMinerSet) {
		//nolint: errchkjson
		data, _ := json.Marshal(msr)
		log.V(4).Info("minerset record changed", "newest", string(data))
		return ctrl.Result{}, r.Store.MinerSets().Update(ctx, msr)
	}

	return ctrl.Result{}, nil
}

// create minerset record.
func addMinerSet(ctx context.Context, dbcli store.IStore, ms *v1beta1.MinerSet) error {
	return dbcli.MinerSets().Create(ctx, applyToMinerSet(&gwmodel.MinerSetM{}, ms))
}

func applyToMinerSet(msr *gwmodel.MinerSetM, ms *v1beta1.MinerSet) *gwmodel.MinerSetM {
	msr.Namespace = ms.Namespace
	msr.Name = ms.Name
	msr.Replicas = *ms.Spec.Replicas
	msr.DisplayName = ms.Spec.DisplayName
	msr.DeletePolicy = ms.Spec.DeletePolicy
	msr.MinReadySeconds = ms.Spec.MinReadySeconds
	msr.FullyLabeledReplicas = ms.Status.FullyLabeledReplicas
	msr.ReadyReplicas = ms.Status.ReadyReplicas
	msr.AvailableReplicas = ms.Status.AvailableReplicas

	if ms.Status.FailureReason != nil {
		msr.FailureReason = string(*ms.Status.FailureReason)
	}
	if ms.Status.FailureMessage != nil {
		msr.FailureMessage = *ms.Status.FailureMessage
	}

	if len(ms.Status.Conditions) > 0 {
		//nolint:errchkjson
		data, _ := json.Marshal(ms.Status.Conditions)
		msr.Conditions = string(data)
	}

	return msr
}
