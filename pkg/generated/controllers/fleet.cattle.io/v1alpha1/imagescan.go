/*
Copyright (c) 2020 - 2023 SUSE LLC

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

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/v2/pkg/apply"
	"github.com/rancher/wrangler/v2/pkg/condition"
	"github.com/rancher/wrangler/v2/pkg/generic"
	"github.com/rancher/wrangler/v2/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ImageScanController interface for managing ImageScan resources.
type ImageScanController interface {
	generic.ControllerInterface[*v1alpha1.ImageScan, *v1alpha1.ImageScanList]
}

// ImageScanClient interface for managing ImageScan resources in Kubernetes.
type ImageScanClient interface {
	generic.ClientInterface[*v1alpha1.ImageScan, *v1alpha1.ImageScanList]
}

// ImageScanCache interface for retrieving ImageScan resources in memory.
type ImageScanCache interface {
	generic.CacheInterface[*v1alpha1.ImageScan]
}

type ImageScanStatusHandler func(obj *v1alpha1.ImageScan, status v1alpha1.ImageScanStatus) (v1alpha1.ImageScanStatus, error)

type ImageScanGeneratingHandler func(obj *v1alpha1.ImageScan, status v1alpha1.ImageScanStatus) ([]runtime.Object, v1alpha1.ImageScanStatus, error)

func RegisterImageScanStatusHandler(ctx context.Context, controller ImageScanController, condition condition.Cond, name string, handler ImageScanStatusHandler) {
	statusHandler := &imageScanStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, generic.FromObjectHandlerToHandler(statusHandler.sync))
}

func RegisterImageScanGeneratingHandler(ctx context.Context, controller ImageScanController, apply apply.Apply,
	condition condition.Cond, name string, handler ImageScanGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &imageScanGeneratingHandler{
		ImageScanGeneratingHandler: handler,
		apply:                      apply,
		name:                       name,
		gvk:                        controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterImageScanStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type imageScanStatusHandler struct {
	client    ImageScanClient
	condition condition.Cond
	handler   ImageScanStatusHandler
}

func (a *imageScanStatusHandler) sync(key string, obj *v1alpha1.ImageScan) (*v1alpha1.ImageScan, error) {
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
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type imageScanGeneratingHandler struct {
	ImageScanGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *imageScanGeneratingHandler) Remove(key string, obj *v1alpha1.ImageScan) (*v1alpha1.ImageScan, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1alpha1.ImageScan{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *imageScanGeneratingHandler) Handle(obj *v1alpha1.ImageScan, status v1alpha1.ImageScanStatus) (v1alpha1.ImageScanStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ImageScanGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
