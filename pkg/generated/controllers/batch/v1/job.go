/*
Copyright The Kubernetes Authors.

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

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	v1 "k8s.io/api/batch/v1"
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

type JobHandler func(string, *v1.Job) (*v1.Job, error)

type JobController interface {
	generic.ControllerMeta
	JobClient

	OnChange(ctx context.Context, name string, sync JobHandler)
	OnRemove(ctx context.Context, name string, sync JobHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() JobCache
}

type JobClient interface {
	Create(*v1.Job) (*v1.Job, error)
	Update(*v1.Job) (*v1.Job, error)
	UpdateStatus(*v1.Job) (*v1.Job, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Job, error)
	List(namespace string, opts metav1.ListOptions) (*v1.JobList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Job, err error)
}

type JobCache interface {
	Get(namespace, name string) (*v1.Job, error)
	List(namespace string, selector labels.Selector) ([]*v1.Job, error)

	AddIndexer(indexName string, indexer JobIndexer)
	GetByIndex(indexName, key string) ([]*v1.Job, error)
}

type JobIndexer func(obj *v1.Job) ([]string, error)

type jobController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewJobController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) JobController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &jobController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromJobHandlerToHandler(sync JobHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Job
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Job))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *jobController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Job))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateJobDeepCopyOnChange(client JobClient, obj *v1.Job, handler func(obj *v1.Job) (*v1.Job, error)) (*v1.Job, error) {
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

func (c *jobController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *jobController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *jobController) OnChange(ctx context.Context, name string, sync JobHandler) {
	c.AddGenericHandler(ctx, name, FromJobHandlerToHandler(sync))
}

func (c *jobController) OnRemove(ctx context.Context, name string, sync JobHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromJobHandlerToHandler(sync)))
}

func (c *jobController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *jobController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *jobController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *jobController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *jobController) Cache() JobCache {
	return &jobCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *jobController) Create(obj *v1.Job) (*v1.Job, error) {
	result := &v1.Job{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *jobController) Update(obj *v1.Job) (*v1.Job, error) {
	result := &v1.Job{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *jobController) UpdateStatus(obj *v1.Job) (*v1.Job, error) {
	result := &v1.Job{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *jobController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *jobController) Get(namespace, name string, options metav1.GetOptions) (*v1.Job, error) {
	result := &v1.Job{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *jobController) List(namespace string, opts metav1.ListOptions) (*v1.JobList, error) {
	result := &v1.JobList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *jobController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *jobController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.Job, error) {
	result := &v1.Job{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type jobCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *jobCache) Get(namespace, name string) (*v1.Job, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.Job), nil
}

func (c *jobCache) List(namespace string, selector labels.Selector) (ret []*v1.Job, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Job))
	})

	return ret, err
}

func (c *jobCache) AddIndexer(indexName string, indexer JobIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Job))
		},
	}))
}

func (c *jobCache) GetByIndex(indexName, key string) (result []*v1.Job, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.Job, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.Job))
	}
	return result, nil
}

type JobStatusHandler func(obj *v1.Job, status v1.JobStatus) (v1.JobStatus, error)

type JobGeneratingHandler func(obj *v1.Job, status v1.JobStatus) ([]runtime.Object, v1.JobStatus, error)

func RegisterJobStatusHandler(ctx context.Context, controller JobController, condition condition.Cond, name string, handler JobStatusHandler) {
	statusHandler := &jobStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromJobHandlerToHandler(statusHandler.sync))
}

func RegisterJobGeneratingHandler(ctx context.Context, controller JobController, apply apply.Apply,
	condition condition.Cond, name string, handler JobGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &jobGeneratingHandler{
		JobGeneratingHandler: handler,
		apply:                apply,
		name:                 name,
		gvk:                  controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterJobStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type jobStatusHandler struct {
	client    JobClient
	condition condition.Cond
	handler   JobStatusHandler
}

func (a *jobStatusHandler) sync(key string, obj *v1.Job) (*v1.Job, error) {
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

type jobGeneratingHandler struct {
	JobGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *jobGeneratingHandler) Remove(key string, obj *v1.Job) (*v1.Job, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.Job{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *jobGeneratingHandler) Handle(obj *v1.Job, status v1.JobStatus) (v1.JobStatus, error) {
	objs, newStatus, err := a.JobGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
