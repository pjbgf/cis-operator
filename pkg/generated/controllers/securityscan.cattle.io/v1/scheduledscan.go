/*
Copyright 2020 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/rancher/clusterscan-operator/pkg/apis/securityscan.cattle.io/v1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ScheduledScanHandler func(string, *v1.ScheduledScan) (*v1.ScheduledScan, error)

type ScheduledScanController interface {
	generic.ControllerMeta
	ScheduledScanClient

	OnChange(ctx context.Context, name string, sync ScheduledScanHandler)
	OnRemove(ctx context.Context, name string, sync ScheduledScanHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() ScheduledScanCache
}

type ScheduledScanClient interface {
	Create(*v1.ScheduledScan) (*v1.ScheduledScan, error)
	Update(*v1.ScheduledScan) (*v1.ScheduledScan, error)
	UpdateStatus(*v1.ScheduledScan) (*v1.ScheduledScan, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1.ScheduledScan, error)
	List(opts metav1.ListOptions) (*v1.ScheduledScanList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ScheduledScan, err error)
}

type ScheduledScanCache interface {
	Get(name string) (*v1.ScheduledScan, error)
	List(selector labels.Selector) ([]*v1.ScheduledScan, error)

	AddIndexer(indexName string, indexer ScheduledScanIndexer)
	GetByIndex(indexName, key string) ([]*v1.ScheduledScan, error)
}

type ScheduledScanIndexer func(obj *v1.ScheduledScan) ([]string, error)

type scheduledScanController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewScheduledScanController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ScheduledScanController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &scheduledScanController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromScheduledScanHandlerToHandler(sync ScheduledScanHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.ScheduledScan
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.ScheduledScan))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *scheduledScanController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.ScheduledScan))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateScheduledScanDeepCopyOnChange(client ScheduledScanClient, obj *v1.ScheduledScan, handler func(obj *v1.ScheduledScan) (*v1.ScheduledScan, error)) (*v1.ScheduledScan, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *scheduledScanController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *scheduledScanController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *scheduledScanController) OnChange(ctx context.Context, name string, sync ScheduledScanHandler) {
	c.AddGenericHandler(ctx, name, FromScheduledScanHandlerToHandler(sync))
}

func (c *scheduledScanController) OnRemove(ctx context.Context, name string, sync ScheduledScanHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromScheduledScanHandlerToHandler(sync)))
}

func (c *scheduledScanController) Enqueue(name string) {
	c.controller.Enqueue("", name)
}

func (c *scheduledScanController) EnqueueAfter(name string, duration time.Duration) {
	c.controller.EnqueueAfter("", name, duration)
}

func (c *scheduledScanController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *scheduledScanController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *scheduledScanController) Cache() ScheduledScanCache {
	return &scheduledScanCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *scheduledScanController) Create(obj *v1.ScheduledScan) (*v1.ScheduledScan, error) {
	result := &v1.ScheduledScan{}
	return result, c.client.Create(context.TODO(), "", obj, result, metav1.CreateOptions{})
}

func (c *scheduledScanController) Update(obj *v1.ScheduledScan) (*v1.ScheduledScan, error) {
	result := &v1.ScheduledScan{}
	return result, c.client.Update(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *scheduledScanController) UpdateStatus(obj *v1.ScheduledScan) (*v1.ScheduledScan, error) {
	result := &v1.ScheduledScan{}
	return result, c.client.UpdateStatus(context.TODO(), "", obj, result, metav1.UpdateOptions{})
}

func (c *scheduledScanController) Delete(name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), "", name, *options)
}

func (c *scheduledScanController) Get(name string, options metav1.GetOptions) (*v1.ScheduledScan, error) {
	result := &v1.ScheduledScan{}
	return result, c.client.Get(context.TODO(), "", name, result, options)
}

func (c *scheduledScanController) List(opts metav1.ListOptions) (*v1.ScheduledScanList, error) {
	result := &v1.ScheduledScanList{}
	return result, c.client.List(context.TODO(), "", result, opts)
}

func (c *scheduledScanController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), "", opts)
}

func (c *scheduledScanController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*v1.ScheduledScan, error) {
	result := &v1.ScheduledScan{}
	return result, c.client.Patch(context.TODO(), "", name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type scheduledScanCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *scheduledScanCache) Get(name string) (*v1.ScheduledScan, error) {
	obj, exists, err := c.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.ScheduledScan), nil
}

func (c *scheduledScanCache) List(selector labels.Selector) (ret []*v1.ScheduledScan, err error) {

	err = cache.ListAll(c.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ScheduledScan))
	})

	return ret, err
}

func (c *scheduledScanCache) AddIndexer(indexName string, indexer ScheduledScanIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.ScheduledScan))
		},
	}))
}

func (c *scheduledScanCache) GetByIndex(indexName, key string) (result []*v1.ScheduledScan, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.ScheduledScan, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.ScheduledScan))
	}
	return result, nil
}

type ScheduledScanStatusHandler func(obj *v1.ScheduledScan, status v1.ScheduledScanStatus) (v1.ScheduledScanStatus, error)

type ScheduledScanGeneratingHandler func(obj *v1.ScheduledScan, status v1.ScheduledScanStatus) ([]runtime.Object, v1.ScheduledScanStatus, error)

func RegisterScheduledScanStatusHandler(ctx context.Context, controller ScheduledScanController, condition condition.Cond, name string, handler ScheduledScanStatusHandler) {
	statusHandler := &scheduledScanStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromScheduledScanHandlerToHandler(statusHandler.sync))
}

func RegisterScheduledScanGeneratingHandler(ctx context.Context, controller ScheduledScanController, apply apply.Apply,
	condition condition.Cond, name string, handler ScheduledScanGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &scheduledScanGeneratingHandler{
		ScheduledScanGeneratingHandler: handler,
		apply:                          apply,
		name:                           name,
		gvk:                            controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterScheduledScanStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type scheduledScanStatusHandler struct {
	client    ScheduledScanClient
	condition condition.Cond
	handler   ScheduledScanStatusHandler
}

func (a *scheduledScanStatusHandler) sync(key string, obj *v1.ScheduledScan) (*v1.ScheduledScan, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		var newErr error
		obj.Status = newStatus
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type scheduledScanGeneratingHandler struct {
	ScheduledScanGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *scheduledScanGeneratingHandler) Remove(key string, obj *v1.ScheduledScan) (*v1.ScheduledScan, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.ScheduledScan{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *scheduledScanGeneratingHandler) Handle(obj *v1.ScheduledScan, status v1.ScheduledScanStatus) (v1.ScheduledScanStatus, error) {
	objs, newStatus, err := a.ScheduledScanGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
