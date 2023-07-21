//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.

// Code generated by defaulter-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&Chain{}, func(obj interface{}) { SetObjectDefaults_Chain(obj.(*Chain)) })
	scheme.AddTypeDefaultingFunc(&ChainList{}, func(obj interface{}) { SetObjectDefaults_ChainList(obj.(*ChainList)) })
	scheme.AddTypeDefaultingFunc(&Miner{}, func(obj interface{}) { SetObjectDefaults_Miner(obj.(*Miner)) })
	scheme.AddTypeDefaultingFunc(&MinerList{}, func(obj interface{}) { SetObjectDefaults_MinerList(obj.(*MinerList)) })
	scheme.AddTypeDefaultingFunc(&MinerSet{}, func(obj interface{}) { SetObjectDefaults_MinerSet(obj.(*MinerSet)) })
	scheme.AddTypeDefaultingFunc(&MinerSetList{}, func(obj interface{}) { SetObjectDefaults_MinerSetList(obj.(*MinerSetList)) })
	return nil
}

func SetObjectDefaults_Chain(in *Chain) {
	SetDefaults_Chain(in)
	SetDefaults_ChainSpec(&in.Spec)
}

func SetObjectDefaults_ChainList(in *ChainList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_Chain(a)
	}
}

func SetObjectDefaults_Miner(in *Miner) {
	SetDefaults_Miner(in)
	SetDefaults_MinerSpec(&in.Spec)
}

func SetObjectDefaults_MinerList(in *MinerList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_Miner(a)
	}
}

func SetObjectDefaults_MinerSet(in *MinerSet) {
	SetDefaults_MinerSpec(&in.Spec.Template.Spec)
}

func SetObjectDefaults_MinerSetList(in *MinerSetList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_MinerSet(a)
	}
}