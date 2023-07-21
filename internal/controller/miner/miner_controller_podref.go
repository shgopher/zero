// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package miner

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/superproj/zero/internal/pkg/util/annotations"
	"github.com/superproj/zero/internal/pkg/util/conditions"
	minerutil "github.com/superproj/zero/internal/pkg/util/miner"
	podutil "github.com/superproj/zero/internal/pkg/util/pod"
	"github.com/superproj/zero/pkg/apis/apps/v1beta1"
	cmerrors "github.com/superproj/zero/pkg/errors"
	"github.com/superproj/zero/pkg/util/record"
)

// ErrPodNotFound signals that a corev1.Pod could not be found for the given provider id.
var ErrPodNotFound = errors.New("cannot find pod with matching ProviderID")

func (r *Reconciler) reconcileProviderPod(ctx context.Context, m *v1beta1.Miner) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// Even if Status.PodRef exists, continue to do the following checks to make sure Pod is healthy
	pod, err := r.ProviderClient.CoreV1().Pods(m.Namespace).Get(ctx, m.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// While a PodRef is set in the status, failing to get that pod means the pod is deleted.
			// If Status.PodRef is not set before, pod still can be in the provisioning state.
			if m.Status.PodRef != nil {
				conditions.MarkFalse(m, v1beta1.MinerPodHealthyCondition, v1beta1.PodNotFoundReason, v1beta1.ConditionSeverityError, "")
				m.Status.FailureReason = cmerrors.MinerStatusErrorPtr(cmerrors.InvalidConfigurationMinerError)
				m.Status.FailureMessage = pointer.String(fmt.Sprintf("Miner pod %q has been deleted after being ready", m.Status.PodRef.Name))
				record.Warnf(m, v1beta1.DeletedReason, "Miner pod %q has been deleted after being ready", m.Status.PodRef.Name)
				log.Error(err, "no matching Pod for Miner in the same namespace")
				return ctrl.Result{}, err
			}
			conditions.MarkFalse(m, v1beta1.MinerPodHealthyCondition, v1beta1.PodProvisioningReason, v1beta1.ConditionSeverityWarning, "")
			// No need to requeue here. Pods emit an event that triggers reconciliation.
			// return ctrl.Result{}, nil
			return r.createMinerPod(ctx, m)
		}

		log.Error(err, "Failed to retrieve pod by miner name")
		record.Warn(m, "FailedCreate", err.Error())
		return ctrl.Result{}, err
	}

	// Set the Miner PodRef.
	if m.Status.PodRef == nil {
		m.Status.PodRef = &corev1.ObjectReference{
			Kind:       pod.Kind,
			APIVersion: pod.APIVersion,
			Name:       pod.Name,
			Namespace:  pod.Namespace,
			UID:        pod.UID,
		}
		log.Info("Kubernetes pod is now available", "pod", klog.KObj(pod))
		record.Eventf(m, "SuccessfulSetPodRef", "Set podRef with %q", m.Status.PodRef.Name)
	}

	// Reconcile pod annotations.
	objPatch := client.MergeFrom(pod.DeepCopy())
	desired := map[string]string{
		v1beta1.MinerNamespaceAnnotation: m.GetNamespace(),
		v1beta1.MinerAnnotation:          m.Name,
	}
	if owner := metav1.GetControllerOfNoCopy(m); owner != nil {
		desired[v1beta1.OwnerKindAnnotation] = owner.Kind
		desired[v1beta1.OwnerNameAnnotation] = owner.Name
	}
	if annotations.AddAnnotations(pod, desired) {
		patchBytes, _ := objPatch.Data(pod)
		if _, err := r.ProviderClient.CoreV1().Pods(pod.Namespace).Patch(ctx, pod.Name, types.MergePatchType, patchBytes, metav1.PatchOptions{}); err != nil {
			log.V(2).Info("Failed patch pod to set annotations", "err", err, "podName", pod.Name)
			return ctrl.Result{}, err
		}
	}

	// Do the remaining pod health checks, then set the pod health to true if all checks pass.
	markPodHealthyCondition(m, pod)
	if podutil.IsPodReady(pod) {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{Requeue: true}, nil
}

func (r *Reconciler) createMinerPod(ctx context.Context, m *v1beta1.Miner) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	ch := &v1beta1.Chain{}
	key := client.ObjectKey{Namespace: metav1.NamespaceSystem, Name: m.Spec.ChainName}
	if err := r.client.Get(ctx, key, ch); err != nil {
		record.Warnf(m, "FailedCreate", "Failed to get chain %s: %v", key, err)
		return ctrl.Result{}, err
	}

	// service name is the same as miner name
	bootstrapIP := minerutil.GenesisDNSServiceNameFromMiner(ch.Status.MinerRef.Name)
	minerType, ok := r.ComponentConfig.Types[m.Spec.MinerType]
	if !ok {
		errMessage := fmt.Sprintf("Miner's miner type %s is unsupported", m.Spec.MinerType)
		log.Error(fmt.Errorf("Miner's miner type is unsupported"), errMessage)
		m.Status.FailureReason = cmerrors.MinerStatusErrorPtr(cmerrors.InvalidConfigurationMinerError)
		m.Status.FailureMessage = pointer.String(errMessage)
		record.Warn(m, "MinerTypeUnsupported", errMessage)
		return ctrl.Result{}, fmt.Errorf(errMessage)
	}

	args := []string{
		"--p2p-addr=0.0.0.0:6001",
		//nolint:nosprintfhostport
		"--peers=" + fmt.Sprintf("ws://%s:6001", bootstrapIP),
		"--http.addr=0.0.0.0:58000",
	}
	if !minerutil.IsGenesisMiner(m) {
		args = append(args,
			"--miner",
			"--address="+m.Namespace,
			"--min-mine-interval="+(time.Duration(ch.Spec.MinMineIntervalSeconds)*time.Second).String(),
			"--mining-difficulty="+strconv.Itoa(minerType.MiningDifficulty),
		)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   m.Namespace,
			Name:        m.Name,
			Annotations: map[string]string{v1beta1.MinerAnnotation: m.Name},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: m.Spec.RestartPolicy,
			Containers: []corev1.Container{
				{
					ImagePullPolicy: corev1.PullIfNotPresent,
					Name:            "toyblc",
					Image:           pointer.StringDeref(GetMinerEnv().Image, ch.Spec.Image),
					Args:            args,
					Resources: corev1.ResourceRequirements{
						Limits: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceCPU:    minerType.CPU,
							corev1.ResourceMemory: minerType.Memory,
						},
					},
				},
			},
		},
	}

	if _, err := r.ProviderClient.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

/*
func (r *Reconciler) getPod(ctx context.Context, podsGetter v1clientset.PodsGetter, namespace string, name string) (*corev1.Pod, error) {
	log := ctrl.LoggerFrom(ctx)

	podList := corev1.PodList{}
	if err := c.List(ctx, &podList, client.MatchingFields{index.PodProviderIDField: providerID.IndexKey()}); err != nil {
		return nil, err
	}
	if len(nodeList.Items) == 0 {
		// If for whatever reason the index isn't registered or available, we fallback to loop over the whole list.
		nl := corev1.PodList{}
		for {
			if err := c.List(ctx, &nl, client.Continue(nl.Continue)); err != nil {
				return nil, err
			}

			for key, node := range nl.Items {
				nodeProviderID, err := podutil.NewProviderID(node.Spec.ProviderID)
				if err != nil {
					log.Error(err, "Failed to parse ProviderID", "Pod", klog.KRef("", nl.Items[key].GetName()))
					continue
				}

				if providerID.Equals(nodeProviderID) {
					return &node, nil
				}
			}

			if nl.Continue == "" {
				break
			}
		}

		return nil, ErrPodNotFound
	}

	if len(nodeList.Items) != 1 {
		return nil, fmt.Errorf("unexpectedly found more than one Pod matching the providerID %s", providerID.String())
	}

	return &nodeList.Items[0], nil
}
*/

func markPodHealthyCondition(miner *v1beta1.Miner, pod *corev1.Pod) {
	cond := podutil.GetReadyCondition(&pod.Status)
	if cond == nil {
		return
	}

	switch cond.Status {
	case corev1.ConditionTrue:
		conditions.MarkTrue(miner, v1beta1.MinerPodHealthyCondition)
	case corev1.ConditionFalse:
		conditions.MarkFalse(miner, v1beta1.MinerPodHealthyCondition, v1beta1.PodConditionsFailedReason, v1beta1.ConditionSeverityWarning, cond.Message)
	default:
		conditions.MarkUnknown(miner, v1beta1.MinerPodHealthyCondition, v1beta1.PodConditionsFailedReason, cond.Message)
	}
}
