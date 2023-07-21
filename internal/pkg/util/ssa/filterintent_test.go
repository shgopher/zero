// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package ssa

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/superproj/zero/internal/pkg/contract"
)

func Test_filterNotAllowedPaths(t *testing.T) {
	tests := []struct {
		name      string
		ctx       *FilterIntentInput
		wantValue map[string]interface{}
	}{
		{
			name: "Filters out not allowed paths",
			ctx: &FilterIntentInput{
				Path: contract.Path{},
				Value: map[string]interface{}{
					"apiVersion": "foo.bar/v1",
					"kind":       "Foo",
					"metadata": map[string]interface{}{
						"name":      "foo",
						"namespace": "bar",
						"labels": map[string]interface{}{
							"foo": "123",
						},
						"annotations": map[string]interface{}{
							"foo": "123",
						},
						"resourceVersion": "123",
					},
					"spec": map[string]interface{}{
						"foo": "123",
					},
					"status": map[string]interface{}{
						"foo": "123",
					},
				},
				ShouldFilter: IsPathNotAllowed(
					[]contract.Path{ // NOTE: we are dropping everything not in this list
						{"apiVersion"},
						{"kind"},
						{"metadata", "name"},
						{"metadata", "namespace"},
						{"metadata", "labels"},
						{"metadata", "annotations"},
						{"spec"},
					},
				),
			},
			wantValue: map[string]interface{}{
				"apiVersion": "foo.bar/v1",
				"kind":       "Foo",
				"metadata": map[string]interface{}{
					"name":      "foo",
					"namespace": "bar",
					"labels": map[string]interface{}{
						"foo": "123",
					},
					"annotations": map[string]interface{}{
						"foo": "123",
					},
					// metadata.resourceVersion filtered out
				},
				"spec": map[string]interface{}{
					"foo": "123",
				},
				// status filtered out
			},
		},
		{
			name: "Cleanup empty maps",
			ctx: &FilterIntentInput{
				Path: contract.Path{},
				Value: map[string]interface{}{
					"spec": map[string]interface{}{
						"foo": "123",
					},
				},
				ShouldFilter: IsPathNotAllowed(
					[]contract.Path{}, // NOTE: we are filtering out everything not in this list (everything)
				),
			},
			wantValue: map[string]interface{}{
				// we are filtering out spec.foo and then spec given that it is an empty map
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			FilterIntent(tt.ctx)

			g.Expect(tt.ctx.Value).To(Equal(tt.wantValue))
		})
	}
}

func Test_filterIgnoredPaths(t *testing.T) {
	tests := []struct {
		name      string
		ctx       *FilterIntentInput
		wantValue map[string]interface{}
	}{
		{
			name: "Filters out ignore paths",
			ctx: &FilterIntentInput{
				Path: contract.Path{},
				Value: map[string]interface{}{
					"spec": map[string]interface{}{
						"foo": "bar",
						"controlPlaneEndpoint": map[string]interface{}{
							"host": "foo-changed",
							"port": "123-changed",
						},
					},
				},
				ShouldFilter: IsPathIgnored(
					[]contract.Path{
						{"spec", "controlPlaneEndpoint"},
					},
				),
			},
			wantValue: map[string]interface{}{
				"spec": map[string]interface{}{
					"foo": "bar",
					// spec.controlPlaneEndpoint filtered out
				},
			},
		},
		{
			name: "Cleanup empty maps",
			ctx: &FilterIntentInput{
				Path: contract.Path{},
				Value: map[string]interface{}{
					"spec": map[string]interface{}{
						"foo": "123",
					},
				},
				ShouldFilter: IsPathIgnored(
					[]contract.Path{
						{"spec", "foo"},
					},
				),
			},
			wantValue: map[string]interface{}{
				// we are filtering out spec.foo and then spec given that it is an empty map
			},
		},
		{
			name: "Cleanup empty nested maps",
			ctx: &FilterIntentInput{
				Path: contract.Path{},
				Value: map[string]interface{}{
					"spec": map[string]interface{}{
						"bar": map[string]interface{}{
							"foo": "123",
						},
					},
				},
				ShouldFilter: IsPathIgnored(
					[]contract.Path{
						{"spec", "bar", "foo"},
					},
				),
			},
			wantValue: map[string]interface{}{
				// we are filtering out spec.bar.foo and then spec given that it is an empty map
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			FilterIntent(tt.ctx)

			g.Expect(tt.ctx.Value).To(Equal(tt.wantValue))
		})
	}
}
