// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package initializer

import (
	"k8s.io/apiserver/pkg/admission"

	clientset "github.com/superproj/zero/pkg/generated/clientset/versioned"
	informers "github.com/superproj/zero/pkg/generated/informers/externalversions"
)

// WantsInternalMinerInformerFactory defines a function which sets InformerFactory for admission plugins that need it.
type WantsInternalMinerInformerFactory interface {
	admission.InitializationValidator
	SetInternalMinerInformerFactory(informers.SharedInformerFactory)
}

// WantsExternalMinerClientSet defines a function which sets external ClientSet for admission plugins that need it.
type WantsExternalMinerClientSet interface {
	admission.InitializationValidator
	SetExternalMinerClientSet(clientset.Interface)
}
