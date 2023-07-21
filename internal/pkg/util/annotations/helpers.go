// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

// Package annotations implements annotation helper functions.
package annotations

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/superproj/zero/pkg/apis/apps/v1beta1"
)

/*
// IsPaused returns true if the MinerSet is paused.
func IsPaused(ms *v1beta1.MinerSet) bool {
	if ms.Spec.Paused {
		return true
	}
	return false
}
*/

// IsPaused returns true if the object has the `paused` annotation.
func IsPaused(o metav1.Object) bool {
	return HasPausedAnnotation(o)
}

// HasPausedAnnotation returns true if the object has the `paused` annotation.
func HasPausedAnnotation(o metav1.Object) bool {
	return hasAnnotation(o, v1beta1.PausedAnnotation)
}

// hasAnnotation returns true if the object has the specified annotation.
func hasAnnotation(o metav1.Object, annotation string) bool {
	annotations := o.GetAnnotations()
	if annotations == nil {
		return false
	}
	_, ok := annotations[annotation]
	return ok
}

// HasSkipRemediation returns true if the object has the `skip-remediation` annotation.
func HasSkipRemediation(o metav1.Object) bool {
	return hasAnnotation(o, v1beta1.MinerSkipRemediationAnnotation)
}

// HasWithPrefix returns true if at least one of the annotations has the prefix specified.
func HasWithPrefix(prefix string, annotations map[string]string) bool {
	for key := range annotations {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// AddAnnotations sets the desired annotations on the object and returns true if the annotations have changed.
func AddAnnotations(o metav1.Object, desired map[string]string) bool {
	if len(desired) == 0 {
		return false
	}
	annotations := o.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
		o.SetAnnotations(annotations)
	}
	hasChanged := false
	for k, v := range desired {
		if cur, ok := annotations[k]; !ok || cur != v {
			annotations[k] = v
			hasChanged = true
		}
	}
	return hasChanged
}
