// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/superproj/zero/pkg/apis/apps/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMiners implements MinerInterface
type FakeMiners struct {
	Fake *FakeAppsV1beta1
	ns   string
}

var minersResource = schema.GroupVersionResource{Group: "apps.zero.io", Version: "v1beta1", Resource: "miners"}

var minersKind = schema.GroupVersionKind{Group: "apps.zero.io", Version: "v1beta1", Kind: "Miner"}

// Get takes name of the miner, and returns the corresponding miner object, and an error if there is any.
func (c *FakeMiners) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Miner, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(minersResource, c.ns, name), &v1beta1.Miner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Miner), err
}

// List takes label and field selectors, and returns the list of Miners that match those selectors.
func (c *FakeMiners) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.MinerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(minersResource, minersKind, c.ns, opts), &v1beta1.MinerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.MinerList{ListMeta: obj.(*v1beta1.MinerList).ListMeta}
	for _, item := range obj.(*v1beta1.MinerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested miners.
func (c *FakeMiners) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(minersResource, c.ns, opts))

}

// Create takes the representation of a miner and creates it.  Returns the server's representation of the miner, and an error, if there is any.
func (c *FakeMiners) Create(ctx context.Context, miner *v1beta1.Miner, opts v1.CreateOptions) (result *v1beta1.Miner, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(minersResource, c.ns, miner), &v1beta1.Miner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Miner), err
}

// Update takes the representation of a miner and updates it. Returns the server's representation of the miner, and an error, if there is any.
func (c *FakeMiners) Update(ctx context.Context, miner *v1beta1.Miner, opts v1.UpdateOptions) (result *v1beta1.Miner, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(minersResource, c.ns, miner), &v1beta1.Miner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Miner), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMiners) UpdateStatus(ctx context.Context, miner *v1beta1.Miner, opts v1.UpdateOptions) (*v1beta1.Miner, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(minersResource, "status", c.ns, miner), &v1beta1.Miner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Miner), err
}

// Delete takes name of the miner and deletes it. Returns an error if one occurs.
func (c *FakeMiners) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(minersResource, c.ns, name, opts), &v1beta1.Miner{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMiners) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(minersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.MinerList{})
	return err
}

// Patch applies the patch and returns the patched miner.
func (c *FakeMiners) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Miner, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(minersResource, c.ns, name, pt, data, subresources...), &v1beta1.Miner{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Miner), err
}
