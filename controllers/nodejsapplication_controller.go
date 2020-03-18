/*
Copyright 2020 the original author or authors.

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

package controllers

import (
	"context"
	"fmt"

	"github.com/projectriff/system/pkg/controllers"
	mononokev1alpha1 "github.com/spring-cloud-incubator/mononoke/api/v1alpha1"
	"github.com/spring-cloud-incubator/mononoke/cnb"
	"github.com/spring-cloud-incubator/mononoke/opinions"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:rbac:groups=apps.mononoke.local,resources=nodejsapplications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.mononoke.local,resources=nodejsapplications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete

func NodeJsApplicationReconciler(c controllers.Config, registry cnb.Registry) *controllers.ParentReconciler {
	c.Log = c.Log.WithName("NodeJsApplication")

	return &controllers.ParentReconciler{
		Type: &mononokev1alpha1.NodeJsApplication{},
		SubReconcilers: []controllers.SubReconciler{
			NodeJsApplicationResolveImageMetadata(c, registry),
			NodeJsApplicationApplyOpinions(c),
			NodeJsApplicationChildDeploymentReconciler(c),
		},

		Config: c,
	}
}

func NodeJsApplicationResolveImageMetadata(c controllers.Config, registry cnb.Registry) controllers.SubReconciler {
	c.Log = c.Log.WithName("ResolveImageMetadata")
	return &controllers.SyncReconciler{
		Sync: func(ctx context.Context, parent *mononokev1alpha1.NodeJsApplication) error {
			// TODO(scothis) be smarter about which container to use
			applicationContainer := &parent.Spec.Template.Spec.Containers[0]

			ref := applicationContainer.Image
			img, err := registry.GetImage(ref)
			if err != nil {
				return fmt.Errorf("failed to get image %s from registry: %w", ref, err)
			}
			md, err := cnb.ParseBuildMetadata(img)
			if err != nil {
				return fmt.Errorf("failed parse cnb metadata from image %s: %w", ref, err)
			}
			controllers.StashValue(ctx, ImageMetadataStashKey, md)
			// TODO(scothis) update target container with digested image
			// applicationContainer.Image = ...
			return nil
		},

		Config: c,
	}
}

func NodeJsApplicationApplyOpinions(c controllers.Config) controllers.SubReconciler {
	c.Log = c.Log.WithName("ApplyOpinions")

	return &controllers.SyncReconciler{
		Sync: func(ctx context.Context, parent *mononokev1alpha1.NodeJsApplication) error {
			imageMetadata := controllers.RetrieveValue(ctx, ImageMetadataStashKey).(cnb.BuildMetadata)
			applied, err := opinions.NodeJs.Apply(ctx, parent.Spec.Template, imageMetadata)
			if err != nil {
				return err
			}
			parent.Status.AppliedOpinions = applied

			return nil
		},

		Config: c,
	}
}

func NodeJsApplicationChildDeploymentReconciler(c controllers.Config) controllers.SubReconciler {
	c.Log = c.Log.WithName("ChildDeployment")

	return &controllers.ChildReconciler{
		ParentType:    &mononokev1alpha1.NodeJsApplication{},
		ChildType:     &appsv1.Deployment{},
		ChildListType: &appsv1.DeploymentList{},

		DesiredChild: func(ctx context.Context, parent *mononokev1alpha1.NodeJsApplication) (*appsv1.Deployment, error) {
			labels := controllers.MergeMaps(parent.Labels, map[string]string{
				mononokev1alpha1.NodeJsApplicationLabelKey: parent.Name,
			})

			template := *parent.Spec.Template.DeepCopy()
			template.Labels = controllers.MergeMaps(template.Labels, labels)

			child := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: make(map[string]string),
					Name:        parent.Name,
					Namespace:   parent.Namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							mononokev1alpha1.NodeJsApplicationLabelKey: parent.Name,
						},
					},
					Template: template,
				},
			}

			return child, nil
		},
		ReflectChildStatusOnParent: func(parent *mononokev1alpha1.NodeJsApplication, child *appsv1.Deployment, err error) {
			if err != nil {
				if apierrs.IsAlreadyExists(err) {
					name := err.(apierrs.APIStatus).Status().Details.Name
					parent.Status.MarkDeploymentNotOwned(name)
				}
				return
			}
			if child != nil {
				parent.Status.PropagateDeploymentStatus(&child.Status)
			}
		},
		HarmonizeImmutableFields: func(current, desired *appsv1.Deployment) {
			// don't fight with an autoscaler
			desired.Spec.Replicas = current.Spec.Replicas
		},
		MergeBeforeUpdate: func(current, desired *appsv1.Deployment) {
			current.Labels = desired.Labels
			current.Spec = desired.Spec
		},
		SemanticEquals: func(a1, a2 *appsv1.Deployment) bool {
			return equality.Semantic.DeepEqual(a1.Spec, a2.Spec) &&
				equality.Semantic.DeepEqual(a1.Labels, a2.Labels)
		},

		Config:     c,
		IndexField: ".metadata.deploymentController",
		Sanitize: func(child *appsv1.Deployment) interface{} {
			return child.Spec
		},
	}
}
