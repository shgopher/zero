// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package gateway

import (
	"fmt"
	"time"

	"k8s.io/client-go/tools/cache"

	"github.com/superproj/zero/pkg/generated/clientset/versioned"
	informers "github.com/superproj/zero/pkg/generated/informers/externalversions"
	"github.com/superproj/zero/pkg/log"
)

func createInformers(stopCh <-chan struct{}, client versioned.Interface) (informers.SharedInformerFactory, error) {
	f := informers.NewSharedInformerFactory(client, time.Minute)
	msinfor := f.Apps().V1beta1().MinerSets().Informer()
	minfor := f.Apps().V1beta1().Miners().Informer()

	f.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, msinfor.HasSynced, minfor.HasSynced) {
		log.Errorf("failed to wait caches to populate")
		return nil, fmt.Errorf("failed to wait caches to populate")
	}

	return f, nil
}
