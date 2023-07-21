// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package admit

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/admission"
	admissiontesting "k8s.io/apiserver/pkg/admission/testing"
	api "k8s.io/kubernetes/pkg/apis/core"
)

func TestAdmissionNonNilAttribute(t *testing.T) {
	handler := admissiontesting.WithReinvocationTesting(t, NewAlwaysAdmit().(*alwaysAdmit))
	err := handler.Admit(
		context.TODO(),
		admission.NewAttributesRecord(
			nil,
			nil,
			api.Kind("kind").WithVersion("version"),
			"namespace",
			"name",
			api.Resource("resource").WithVersion("version"),
			"subresource",
			admission.Create,
			&metav1.CreateOptions{},
			false,
			nil,
		),
		nil,
	)
	if err != nil {
		t.Errorf("Unexpected error returned from admission handler")
	}
}

func TestAdmissionNilAttribute(t *testing.T) {
	handler := NewAlwaysAdmit()
	err := handler.(*alwaysAdmit).Admit(context.TODO(), nil, nil)
	if err != nil {
		t.Errorf("Unexpected error returned from admission handler")
	}
}

func TestHandles(t *testing.T) {
	handler := NewAlwaysAdmit()
	tests := []admission.Operation{admission.Create, admission.Connect, admission.Update, admission.Delete}

	for _, test := range tests {
		if !handler.Handles(test) {
			t.Errorf("Expected handling all operations, including: %v", test)
		}
	}
}
