// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

package rest

import (
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	"github.com/superproj/zero/internal/apiserver/storage"
	serializerutil "github.com/superproj/zero/internal/pkg/util/serializer"
	chainstore "github.com/superproj/zero/internal/registry/apps/chain/storage"
	minerstore "github.com/superproj/zero/internal/registry/apps/miner/storage"
	minersetstore "github.com/superproj/zero/internal/registry/apps/minerset/storage"
	"github.com/superproj/zero/pkg/apis/apps"
	appsv1beta1 "github.com/superproj/zero/pkg/apis/apps/v1beta1"
)

// RESTStorageProvider is a struct for apps REST storage.
type RESTStorageProvider struct{}

// Implement RESTStorageProvider.
var _ storage.RESTStorageProvider = &RESTStorageProvider{}

// NewRESTStorage is a factory constructor to creates and returns the APIGroupInfo.
func (p RESTStorageProvider) NewRESTStorage(
	apiResourceConfigSource serverstorage.APIResourceConfigSource,
	restOptionsGetter generic.RESTOptionsGetter,
) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(apps.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	apiGroupInfo.NegotiatedSerializer = serializerutil.NewProtocolShieldSerializers(&legacyscheme.Codecs)

	storageMap, err := p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}
	apiGroupInfo.VersionedResourcesStorageMap[appsv1beta1.SchemeGroupVersion.Version] = storageMap

	return apiGroupInfo, nil
}

func (p RESTStorageProvider) v1beta1Storage(
	apiResourceConfigSource serverstorage.APIResourceConfigSource,
	restOptionsGetter generic.RESTOptionsGetter,
) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// chains
	if resource := "chains"; apiResourceConfigSource.ResourceEnabled(appsv1beta1.SchemeGroupVersion.WithResource(resource)) {
		chainStorage, chainStatusStorage, err := chainstore.NewREST(restOptionsGetter)
		if err != nil {
			return storage, err
		}

		storage[resource] = chainStorage
		storage[resource+"/status"] = chainStatusStorage
	}

	// miners
	if resource := "miners"; apiResourceConfigSource.ResourceEnabled(appsv1beta1.SchemeGroupVersion.WithResource(resource)) {
		minerStorage, minerStatusStorage, err := minerstore.NewREST(restOptionsGetter)
		if err != nil {
			return storage, err
		}

		storage[resource] = minerStorage
		storage[resource+"/status"] = minerStatusStorage
	}

	// minersets
	if resource := "minersets"; apiResourceConfigSource.ResourceEnabled(appsv1beta1.SchemeGroupVersion.WithResource(resource)) {
		minerSetStorage, minerSetStatusStorage, minerSetScaleStorage, err := minersetstore.NewREST(restOptionsGetter)
		if err != nil {
			return storage, err
		}

		storage[resource] = minerSetStorage
		storage[resource+"/status"] = minerSetStatusStorage
		storage[resource+"/scale"] = minerSetScaleStorage
	}

	return storage, nil
}

// GroupName return the api group name.
func (p RESTStorageProvider) GroupName() string {
	return apps.GroupName
}
