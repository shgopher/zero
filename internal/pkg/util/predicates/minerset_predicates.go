// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package predicates implements predicate utilities.
package predicates

import (
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	appsv1beta1 "github.com/superproj/zero/pkg/apis/apps/v1beta1"
)

/*
// MinerSetCreateInfraReady returns a predicate that returns true for a create event when a minerset has Status.InfrastructureReady set as true
// it also returns true if the resource provided is not a MinerSet to allow for use with controller-runtime NewControllerManagedBy.
func MinerSetCreateInfraReady(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			log := logger.WithValues("predicate", "MinerSetCreateInfraReady", "eventType", "create")

			c, ok := e.Object.(*appsv1beta1.MinerSet)
			if !ok {
				log.V(4).Info("Expected MinerSet", "type", fmt.Sprintf("%T", e.Object))
				return false
			}
			log = log.WithValues("MinerSet", klog.KObj(c))

			// Only need to trigger a reconcile if the MinerSet.Status.InfrastructureReady is true
			if c.Status.InfrastructureReady {
				log.V(6).Info("MinerSet infrastructure is ready, allowing further processing")
				return true
			}

			log.V(4).Info("MinerSet infrastructure is not ready, blocking further processing")
			return false
		},
		UpdateFunc:  func(e event.UpdateEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}
*/

// MinerSetCreateNotPaused returns a predicate that returns true for a create event when a minerset has Spec.Paused set as false
// it also returns true if the resource provided is not a MinerSet to allow for use with controller-runtime NewControllerManagedBy.
func MinerSetCreateNotPaused(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			log := logger.WithValues("predicate", "MinerSetCreateNotPaused", "eventType", "create")

			c, ok := e.Object.(*appsv1beta1.MinerSet)
			if !ok {
				log.V(4).Info("Expected MinerSet", "type", fmt.Sprintf("%T", e.Object))
				return false
			}
			log = log.WithValues("MinerSet", klog.KObj(c))

			// Only need to trigger a reconcile if the MinerSet.Spec.Paused is false
			if !isPaused(c) {
				log.V(6).Info("MinerSet is not paused, allowing further processing")
				return true
			}

			log.V(4).Info("MinerSet is paused, blocking further processing")
			return false
		},
		UpdateFunc:  func(e event.UpdateEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}

/*
// MinerSetUpdateInfraReady returns a predicate that returns true for an update event when a minerset has Status.InfrastructureReady changed from false to true
// it also returns true if the resource provided is not a MinerSet to allow for use with controller-runtime NewControllerManagedBy.
func MinerSetUpdateInfraReady(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := logger.WithValues("predicate", "MinerSetUpdateInfraReady", "eventType", "update")

			oldMinerSet, ok := e.ObjectOld.(*appsv1beta1.MinerSet)
			if !ok {
				log.V(4).Info("Expected MinerSet", "type", fmt.Sprintf("%T", e.ObjectOld))
				return false
			}
			log = log.WithValues("MinerSet", klog.KObj(oldMinerSet))

			newMinerSet := e.ObjectNew.(*appsv1beta1.MinerSet)

			if !oldMinerSet.Status.InfrastructureReady && newMinerSet.Status.InfrastructureReady {
				log.V(6).Info("MinerSet infrastructure became ready, allowing further processing")
				return true
			}

			log.V(4).Info("MinerSet infrastructure did not become ready, blocking further processing")
			return false
		},
		CreateFunc:  func(e event.CreateEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}
*/

// MinerSetUpdateUnpaused returns a predicate that returns true for an update event when a minerset has Spec.Paused changed from true to false
// it also returns true if the resource provided is not a MinerSet to allow for use with controller-runtime NewControllerManagedBy.
func MinerSetUpdateUnpaused(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := logger.WithValues("predicate", "MinerSetUpdateUnpaused", "eventType", "update")

			oldMinerSet, ok := e.ObjectOld.(*appsv1beta1.MinerSet)
			if !ok {
				log.V(4).Info("Expected MinerSet", "type", fmt.Sprintf("%T", e.ObjectOld))
				return false
			}
			log = log.WithValues("MinerSet", klog.KObj(oldMinerSet))

			newMinerSet := e.ObjectNew.(*appsv1beta1.MinerSet)

			if isPaused(oldMinerSet) && !isPaused(newMinerSet) {
				log.V(4).Info("MinerSet was unpaused, allowing further processing")
				return true
			}

			// This predicate always work in "or" with Paused predicates
			// so the logs are adjusted to not provide false negatives/verbosity al V<=5.
			log.V(6).Info("MinerSet was not unpaused, blocking further processing")
			return false
		},
		CreateFunc:  func(e event.CreateEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}

// MinerSetUnpaused returns a Predicate that returns true on MinerSet creation events where MinerSet.Spec.Paused is false
// and Update events when MinerSet.Spec.Paused transitions to false.
// This implements a common requirement for many minerset-api and provider controllers (such as MinerSet Infrastructure
// controllers) to resume reconciliation when the MinerSet is unpaused.
// Example use:
//
//	err := controller.Watch(
//	    &source.Kind{Type: &appsv1beta1.MinerSet{}},
//	    &handler.EnqueueRequestsFromMapFunc{
//	        ToRequests: minersetToMiners,
//	    },
//	    predicates.MinerSetUnpaused(r.Log),
//	)
func MinerSetUnpaused(logger logr.Logger) predicate.Funcs {
	log := logger.WithValues("predicate", "MinerSetUnpaused")

	// Use any to ensure we process either create or update events we care about
	return Any(log, MinerSetCreateNotPaused(log), MinerSetUpdateUnpaused(log))
}

/*
// MinerSetControlPlaneInitialized returns a Predicate that returns true on Update events
// when ControlPlaneInitializedCondition on a MinerSet changes to true.
// Example use:
//
//	err := controller.Watch(
//	    &source.Kind{Type: &appsv1beta1.MinerSet{}},
//	    &handler.EnqueueRequestsFromMapFunc{
//	        ToRequests: minersetToMiners,
//	    },
//	    predicates.MinerSetControlPlaneInitialized(r.Log),
//	)
func MinerSetControlPlaneInitialized(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := logger.WithValues("predicate", "MinerSetControlPlaneInitialized", "eventType", "update")

			oldMinerSet, ok := e.ObjectOld.(*appsv1beta1.MinerSet)
			if !ok {
				log.V(4).Info("Expected MinerSet", "type", fmt.Sprintf("%T", e.ObjectOld))
				return false
			}
			log = log.WithValues("MinerSet", klog.KObj(oldMinerSet))

			newMinerSet := e.ObjectNew.(*appsv1beta1.MinerSet)

			if !conditions.IsTrue(oldMinerSet, appsv1beta1.ControlPlaneInitializedCondition) &&
				conditions.IsTrue(newMinerSet, appsv1beta1.ControlPlaneInitializedCondition) {
				log.V(6).Info("MinerSet ControlPlaneInitialized was set, allow further processing")
				return true
			}

			log.V(6).Info("MinerSet ControlPlaneInitialized hasn't changed, blocking further processing")
			return false
		},
		CreateFunc:  func(e event.CreateEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}

// MinerSetUnpausedAndInfrastructureReady returns a Predicate that returns true on MinerSet creation events where
// both MinerSet.Spec.Paused is false and MinerSet.Status.InfrastructureReady is true and Update events when
// either MinerSet.Spec.Paused transitions to false or MinerSet.Status.InfrastructureReady transitions to true.
// This implements a common requirement for some minerset-api and provider controllers (such as Miner Infrastructure
// controllers) to resume reconciliation when the MinerSet is unpaused and when the infrastructure becomes ready.
// Example use:
//
//	err := controller.Watch(
//	    &source.Kind{Type: &appsv1beta1.MinerSet{}},
//	    &handler.EnqueueRequestsFromMapFunc{
//	        ToRequests: minersetToMiners,
//	    },
//	    predicates.MinerSetUnpausedAndInfrastructureReady(r.Log),
//	)
func MinerSetUnpausedAndInfrastructureReady(logger logr.Logger) predicate.Funcs {
	log := logger.WithValues("predicate", "MinerSetUnpausedAndInfrastructureReady")

	// Only continue processing create events if both not paused and infrastructure is ready
	createPredicates := All(log, MinerSetCreateNotPaused(log), MinerSetCreateInfraReady(log))

	// Process update events if either MinerSet is unpaused or infrastructure becomes ready
	updatePredicates := Any(log, MinerSetUpdateUnpaused(log), MinerSetUpdateInfraReady(log))

	// Use any to ensure we process either create or update events we care about
	return Any(log, createPredicates, updatePredicates)
}

// MinerSetHasTopology returns a Predicate that returns true when minerset.Spec.Topology
// is NOT nil and false otherwise.
func MinerSetHasTopology(logger logr.Logger) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return processIfTopologyManaged(logger.WithValues("predicate", "MinerSetHasTopology", "eventType", "update"), e.ObjectNew)
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return processIfTopologyManaged(logger.WithValues("predicate", "MinerSetHasTopology", "eventType", "create"), e.Object)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return processIfTopologyManaged(logger.WithValues("predicate", "MinerSetHasTopology", "eventType", "delete"), e.Object)
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return processIfTopologyManaged(logger.WithValues("predicate", "MinerSetHasTopology", "eventType", "generic"), e.Object)
		},
	}
}

func processIfTopologyManaged(logger logr.Logger, object client.Object) bool {
	minerset, ok := object.(*appsv1beta1.MinerSet)
	if !ok {
		logger.V(4).Info("Expected MinerSet", "type", fmt.Sprintf("%T", object))
		return false
	}

	log := logger.WithValues("MinerSet", klog.KObj(minerset))

	if minerset.Spec.Topology != nil {
		log.V(6).Info("MinerSet has topology, allowing further processing")
		return true
	}

	log.V(6).Info("MinerSet does not have topology, blocking further processing")
	return false
}
*/

func isPaused(ms *appsv1beta1.MinerSet) bool {
	_, ok := ms.GetAnnotations()[appsv1beta1.PausedAnnotation]
	return ok
}
