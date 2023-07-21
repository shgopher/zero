// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package miner

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/superproj/zero/internal/pkg/known"
	"github.com/superproj/zero/internal/pkg/util/conditions"
	minerutil "github.com/superproj/zero/internal/pkg/util/miner"
	"github.com/superproj/zero/pkg/apis/apps/v1beta1"
	"github.com/superproj/zero/pkg/util/record"
)

var externalReadyWait = 30 * time.Second

func (r *Reconciler) reconcilePhase(_ context.Context, m *v1beta1.Miner) {
	originalPhase := m.Status.Phase

	// Set the phase to "Pending" if nil.
	if m.Status.Phase == "" {
		m.Status.SetTypedPhase(v1beta1.MinerPhasePending)
	}

	// Set phase to "Provisioning" if podRef has been set and the pod is not ready.
	if m.Status.PodRef != nil && !conditions.IsTrue(m, v1beta1.MinerPodHealthyCondition) {
		m.Status.SetTypedPhase(v1beta1.MinerPhaseProvisioning)
	}

	// Set phase to "Running" if podRef has been set and the pod is ready
	if m.Status.PodRef != nil && conditions.IsTrue(m, v1beta1.MinerPodHealthyCondition) {
		m.Status.SetTypedPhase(v1beta1.MinerPhaseRunning)
	}

	// Set the phase to "Failed" if any of Status.FailureReason or Status.FailureMessage is not-nil.
	if m.Status.FailureReason != nil || m.Status.FailureMessage != nil {
		m.Status.SetTypedPhase(v1beta1.MinerPhaseFailed)
	}

	// Set the phase to "Deleting" if the deletion timestamp is set.
	if !m.DeletionTimestamp.IsZero() {
		m.Status.SetTypedPhase(v1beta1.MinerPhaseDeleting)
	}

	// If the phase has changed, update the LastUpdated timestamp
	if m.Status.Phase != originalPhase {
		now := metav1.Now()
		m.Status.LastUpdated = &now
	}
}

func (r *Reconciler) reconcileAnnotations(ctx context.Context, m *v1beta1.Miner) (ctrl.Result, error) {
	needReconcile := false
	for _, annotation := range []string{known.CPUAnnotation, known.MemoryAnnotation} {
		if _, ok := m.Annotations[annotation]; !ok {
			needReconcile = true
		}
	}

	if !needReconcile {
		return ctrl.Result{}, nil
	}

	if m.Annotations == nil {
		m.Annotations = make(map[string]string)
	}
	cpu := r.ComponentConfig.Types[m.Spec.MinerType].CPU
	memory := r.ComponentConfig.Types[m.Spec.MinerType].Memory
	m.Annotations[known.CPUAnnotation] = cpu.String()
	m.Annotations[known.MemoryAnnotation] = memory.String()

	return ctrl.Result{}, nil
}

/*
func (r *Reconciler) reconcileProviderConfigMap(ctx context.Context, m *v1beta1.Miner) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	if _, err := r.ProviderClient.CoreV1().ConfigMaps(m.Namespace).Get(ctx, m.Spec.AccountFrom.ConfigMapRef.Name, metav1.GetOptions{}); err == nil {
		return ctrl.Result{}, nil
	}

	cm := &corev1.ConfigMap{}
	key := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.AccountFrom.ConfigMapRef.Name}
	if err := r.client.Get(ctx, key, cm); err != nil {
		record.Warnf(m, "FailedCreate", "Failed to get ConfigMap %s: %v", key.Name, err)
		log.Error(err, "Failed to get configMap", "configMap", klog.KRef(key.Namespace, key.Name))
		return ctrl.Result{}, err
	}

	kindcm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cm.Name,
			Namespace: cm.Namespace,
		},
		Data: cm.Data,
	}

	_, err := r.ProviderClient.CoreV1().ConfigMaps(kindcm.Namespace).Create(ctx, kindcm, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return ctrl.Result{}, nil
		}

		record.Warnf(m, "FailedCreate", "Failed to create ConfigMap %s: %v", kindcm.Name, err)
		log.Error(err, "Failed to create configMap", "configMap", klog.KObj(kindcm))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
*/

func (r *Reconciler) reconcileProviderService(ctx context.Context, m *v1beta1.Miner) (ctrl.Result, error) {
	// only create service for genesis miners
	if !minerutil.IsGenesisMiner(m) {
		return ctrl.Result{}, nil
	}

	log := ctrl.LoggerFrom(ctx)
	if _, err := r.ProviderClient.CoreV1().Services(m.Namespace).Get(ctx, minerutil.GetProviderServiceName(m), metav1.GetOptions{}); err == nil {
		return ctrl.Result{}, nil
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: m.Namespace,
			Name:      minerutil.GetProviderServiceName(m),
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: map[string]string{v1beta1.MinerSetNameLabel: m.Name},
			// ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Name:       "websocket",
					Protocol:   corev1.ProtocolTCP,
					Port:       6001,
					TargetPort: intstr.IntOrString{IntVal: 6001},
				},
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.IntOrString{IntVal: 58000},
				},
			},
		},
	}

	_, err := r.ProviderClient.CoreV1().Services(svc.Namespace).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return ctrl.Result{}, nil
		}

		record.Warnf(m, "FailedCreate", "Failed to get Service %s: %v", svc.Name, err)
		log.Error(err, "Failed to create service", "service", klog.KObj(svc))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

/*
// reconcileExternal handles generic unstructured objects referenced by a Miner.
func (r *Reconciler) reconcileExternal(ctx context.Context, cluster *v1beta1.Cluster, m *clusterv1.Miner, ref *corev1.ObjectReference) (external.ReconcileOutput, error) {
	log := ctrl.LoggerFrom(ctx)

	if err := utilconversion.UpdateReferenceAPIContract(ctx, r.Client, r.APIReader, ref); err != nil {
		return external.ReconcileOutput{}, err
	}

	obj, err := external.Get(ctx, r.Client, ref, m.Namespace)
	if err != nil {
		if apierrors.IsNotFound(errors.Cause(err)) {
			log.Info("could not find external ref, requeueing", "refGVK", ref.GroupVersionKind(), "refName", ref.Name, "Miner", klog.KObj(m))
			return external.ReconcileOutput{RequeueAfter: externalReadyWait}, nil
		}
		return external.ReconcileOutput{}, err
	}

	// if external ref is paused, return error.
	if annotations.IsPaused(cluster, obj) {
		log.V(3).Info("External object referenced is paused")
		return external.ReconcileOutput{Paused: true}, nil
	}

	// Initialize the patch helper.
	helper, err := patch.NewHelper(obj, r.Client)
	if err != nil {
		return external.ReconcileOutput{}, err
	}

	// With the migration from v1alpha2 to v1alpha3, Miner controllers should be the owner for the
	// infra Miners, hence remove any existing minerset controller owner reference
	if controller := metav1.GetControllerOf(obj); controller != nil && controller.Kind == "MinerSet" {
		gv, err := schema.ParseGroupVersion(controller.APIVersion)
		if err != nil {
			return external.ReconcileOutput{}, err
		}
		if gv.Group == v1beta1.GroupVersion.Group {
			ownerRefs := util.RemoveOwnerRef(obj.GetOwnerReferences(), *controller)
			obj.SetOwnerReferences(ownerRefs)
		}
	}

	// Set external object ControllerReference to the Miner.
	if err := controllerutil.SetControllerReference(m, obj, r.Client.Scheme()); err != nil {
		return external.ReconcileOutput{}, err
	}

	// Set the Cluster label.
	labels := obj.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[v1beta1.ClusterLabelName] = m.Spec.ClusterName
	obj.SetLabels(labels)

	// Always attempt to Patch the external object.
	if err := helper.Patch(ctx, obj); err != nil {
		return external.ReconcileOutput{}, err
	}

	// Ensure we add a watcher to the external object.
	if err := r.externalTracker.Watch(log, obj, &handler.EnqueueRequestForOwner{OwnerType: &v1beta1.Miner{}}); err != nil {
		return external.ReconcileOutput{}, err
	}

	// Set failure reason and message, if any.
	failureReason, failureMessage, err := external.FailuresFrom(obj)
	if err != nil {
		return external.ReconcileOutput{}, err
	}
	if failureReason != "" {
		minerStatusError := capierrors.MinerStatusError(failureReason)
		m.Status.FailureReason = &minerStatusError
	}
	if failureMessage != "" {
		m.Status.FailureMessage = pointer.String(
			fmt.Sprintf("Failure detected from referenced resource %v with name %q: %s",
				obj.GroupVersionKind(), obj.GetName(), failureMessage),
		)
	}

	return external.ReconcileOutput{Result: obj}, nil
}

// reconcileBootstrap reconciles the Spec.Bootstrap.ConfigRef object on a Miner.
func (r *Reconciler) reconcileBootstrap(ctx context.Context, cluster *v1beta1.Cluster, m *clusterv1.Miner) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// If the bootstrap data is populated, set ready and return.
	if m.Spec.Bootstrap.DataSecretName != nil {
		m.Status.BootstrapReady = true
		conditions.MarkTrue(m, v1beta1.BootstrapReadyCondition)
		return ctrl.Result{}, nil
	}

	// If the Bootstrap ref is nil (and so the miner should use user generated data secret), return.
	if m.Spec.Bootstrap.ConfigRef == nil {
		return ctrl.Result{}, nil
	}

	// Call generic external reconciler if we have an external reference.
	externalResult, err := r.reconcileExternal(ctx, cluster, m, m.Spec.Bootstrap.ConfigRef)
	if err != nil {
		return ctrl.Result{}, err
	}
	if externalResult.RequeueAfter > 0 {
		return ctrl.Result{RequeueAfter: externalResult.RequeueAfter}, nil
	}
	if externalResult.Paused {
		return ctrl.Result{}, nil
	}
	bootstrapConfig := externalResult.Result

	// If the bootstrap config is being deleted, return early.
	if !bootstrapConfig.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}

	// Determine if the bootstrap provider is ready.
	ready, err := external.IsReady(bootstrapConfig)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Report a summary of current status of the bootstrap object defined for this miner.
	conditions.SetMirror(m, v1beta1.BootstrapReadyCondition,
		conditions.UnstructuredGetter(bootstrapConfig),
		conditions.WithFallbackValue(ready, v1beta1.WaitingForDataSecretFallbackReason, clusterv1.ConditionSeverityInfo, ""),
	)

	// If the bootstrap provider is not ready, requeue.
	if !ready {
		log.Info("Waiting for bootstrap provider to generate data secret and report status.ready", bootstrapConfig.GetKind(), klog.KObj(bootstrapConfig))
		return ctrl.Result{RequeueAfter: externalReadyWait}, nil
	}

	// Get and set the name of the secret containing the bootstrap data.
	secretName, _, err := unstructured.NestedString(bootstrapConfig.Object, "status", "dataSecretName")
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to retrieve dataSecretName from bootstrap provider for Miner %q in namespace %q", m.Name, m.Namespace)
	} else if secretName == "" {
		return ctrl.Result{}, errors.Errorf("retrieved empty dataSecretName from bootstrap provider for Miner %q in namespace %q", m.Name, m.Namespace)
	}
	m.Spec.Bootstrap.DataSecretName = pointer.String(secretName)
	if !m.Status.BootstrapReady {
		log.Info("Bootstrap provider generated data secret and reports status.ready", bootstrapConfig.GetKind(), klog.KObj(bootstrapConfig), "Secret", klog.KRef(m.Namespace, secretName))
	}
	m.Status.BootstrapReady = true
	return ctrl.Result{}, nil
}

// reconcileInfrastructure reconciles the Spec.InfrastructureRef object on a Miner.
func (r *Reconciler) reconcileInfrastructure(ctx context.Context, cluster *v1beta1.Cluster, m *clusterv1.Miner) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// Call generic external reconciler.
	infraReconcileResult, err := r.reconcileExternal(ctx, cluster, m, &m.Spec.InfrastructureRef)
	if err != nil {
		return ctrl.Result{}, err
	}
	if infraReconcileResult.RequeueAfter > 0 {
		// Infra object went missing after the miner was up and running
		if m.Status.InfrastructureReady {
			log.Error(err, "Miner infrastructure reference has been deleted after being ready, setting failure state")
			m.Status.FailureReason = capierrors.MinerStatusErrorPtr(capierrors.InvalidConfigurationMinerError)
			m.Status.FailureMessage = pointer.String(fmt.Sprintf("Miner infrastructure resource %v with name %q has been deleted after being ready",
				m.Spec.InfrastructureRef.GroupVersionKind(), m.Spec.InfrastructureRef.Name))
			return ctrl.Result{}, errors.Errorf("could not find %v %q for Miner %q in namespace %q, requeueing", m.Spec.InfrastructureRef.GroupVersionKind().String(), m.Spec.InfrastructureRef.Name, m.Name, m.Namespace)
		}
		return ctrl.Result{RequeueAfter: infraReconcileResult.RequeueAfter}, nil
	}
	// if the external object is paused, return without any further processing
	if infraReconcileResult.Paused {
		return ctrl.Result{}, nil
	}
	infraConfig := infraReconcileResult.Result

	if !infraConfig.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}

	// Determine if the infrastructure provider is ready.
	ready, err := external.IsReady(infraConfig)
	if err != nil {
		return ctrl.Result{}, err
	}
	if ready && !m.Status.InfrastructureReady {
		log.Info("Infrastructure provider has completed miner infrastructure provisioning and reports status.ready", infraConfig.GetKind(), klog.KObj(infraConfig))
	}
	m.Status.InfrastructureReady = ready

	// Report a summary of current status of the infrastructure object defined for this miner.
	conditions.SetMirror(m, v1beta1.InfrastructureReadyCondition,
		conditions.UnstructuredGetter(infraConfig),
		conditions.WithFallbackValue(ready, v1beta1.WaitingForInfrastructureFallbackReason, clusterv1.ConditionSeverityInfo, ""),
	)

	// If the infrastructure provider is not ready, return early.
	if !ready {
		log.Info("Waiting for infrastructure provider to create miner infrastructure and report status.ready", infraConfig.GetKind(), klog.KObj(infraConfig))
		return ctrl.Result{RequeueAfter: externalReadyWait}, nil
	}

	// Get Spec.ProviderID from the infrastructure provider.
	var providerID string
	if err := util.UnstructuredUnmarshalField(infraConfig, &providerID, "spec", "providerID"); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to retrieve Spec.ProviderID from infrastructure provider for Miner %q in namespace %q", m.Name, m.Namespace)
	} else if providerID == "" {
		return ctrl.Result{}, errors.Errorf("retrieved empty Spec.ProviderID from infrastructure provider for Miner %q in namespace %q", m.Name, m.Namespace)
	}

	// Get and set Status.Addresses from the infrastructure provider.
	err = util.UnstructuredUnmarshalField(infraConfig, &m.Status.Addresses, "status", "addresses")
	if err != nil && err != util.ErrUnstructuredFieldNotFound {
		return ctrl.Result{}, errors.Wrapf(err, "failed to retrieve addresses from infrastructure provider for Miner %q in namespace %q", m.Name, m.Namespace)
	}

	// Get and set the failure domain from the infrastructure provider.
	var failureDomain string
	err = util.UnstructuredUnmarshalField(infraConfig, &failureDomain, "spec", "failureDomain")
	switch {
	case err == util.ErrUnstructuredFieldNotFound: // no-op
	case err != nil:
		return ctrl.Result{}, errors.Wrapf(err, "failed to failure domain from infrastructure provider for Miner %q in namespace %q", m.Name, m.Namespace)
	default:
		m.Spec.FailureDomain = pointer.String(failureDomain)
	}

	m.Spec.ProviderID = pointer.String(providerID)
	return ctrl.Result{}, nil
}
*/
