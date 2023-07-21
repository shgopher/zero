// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.

package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using scripts/update-generated-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE. DO NOT EDIT.
var map_Lease = map[string]string{
	"":         "Lease defines a lease concept.",
	"metadata": "More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"spec":     "Specification of the Lease. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status",
}

func (Lease) SwaggerDoc() map[string]string {
	return map_Lease
}

var map_LeaseList = map[string]string{
	"":         "LeaseList is a list of Lease objects.",
	"metadata": "Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata",
	"items":    "Items is a list of schema objects.",
}

func (LeaseList) SwaggerDoc() map[string]string {
	return map_LeaseList
}

var map_LeaseSpec = map[string]string{
	"":                     "LeaseSpec is a specification of a Lease.",
	"holderIdentity":       "holderIdentity contains the identity of the holder of a current lease.",
	"leaseDurationSeconds": "leaseDurationSeconds is a duration that candidates for a lease need to wait to force acquire it. This is measure against time of last observed RenewTime.",
	"acquireTime":          "acquireTime is a time when the current lease was acquired.",
	"renewTime":            "renewTime is a time when the current holder of a lease has last updated the lease.",
	"leaseTransitions":     "leaseTransitions is the number of transitions of a lease between holders.",
}

func (LeaseSpec) SwaggerDoc() map[string]string {
	return map_LeaseSpec
}

// AUTO-GENERATED FUNCTIONS END HERE