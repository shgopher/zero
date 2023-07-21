// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"

	known "github.com/superproj/zero/internal/pkg/known/apiserver"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

/*
// SetDefaults_MinerSet sets defaults for MinerSet
func SetDefaults_MinerSet(obj *MinerSet) {
	// Set MinerSetSpec.Replicas to 1 if it is not set.
	if obj.Spec.Replicas == nil {
		obj.Spec.Replicas = new(int32)
		*obj.Spec.Replicas = 1
	}

	// Set default template
	SetDefaults_MinerSpec(&obj.Spec.Template.Spec)

	// Set default DeletePolicy as Random.
	if obj.Spec.DeletePolicy == "" {
		obj.Spec.DeletePolicy = "Random"
	}
}
*/

// SetDefaults_Miner sets defaults for Miner.
func SetDefaults_Miner(obj *Miner) {
	// Miner name prefix is fixed to `mi-`
	if obj.ObjectMeta.GenerateName == "" {
		obj.ObjectMeta.GenerateName = "mi-"
	}

	SetDefaults_MinerSpec(&obj.Spec)
}

// SetDefaults_MinerSpec sets defaults for Miner spec.
func SetDefaults_MinerSpec(obj *MinerSpec) {
	if obj.MinerType == "" {
		obj.MinerType = known.DefaultNodeMinerType
	}
}

// SetDefaults_Chain sets defaults for Chain.
func SetDefaults_Chain(obj *Chain) {
	SetDefaults_ChainSpec(&obj.Spec)
}

// SetDefaults_ChainSpec sets defaults for Chain spec.
func SetDefaults_ChainSpec(obj *ChainSpec) {
	obj.BootstrapAccount = pointer.String("0x210d9eD12CEA87E33a98AA7Bcb4359eABA9e800e")
	if obj.MinerType == "" {
		obj.MinerType = known.DefaultGenesisMinerType
	}

	if obj.MinMineIntervalSeconds <= 0 {
		obj.MinMineIntervalSeconds = 12 * 60 * 60 // 12 hours
	}
}
