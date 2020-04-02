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

package opinions

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resource interface {
	metav1.ObjectMetaAccessor
	PodTemplate() *corev1.PodTemplateSpec
}

// setAnnotation sets the annotation on both the resource and the resource's
// PodTemplateSpec
func setAnnotation(r Resource, key, value string) {
	parent := r.GetObjectMeta().GetAnnotations()
	if parent == nil {
		parent = map[string]string{}
		r.GetObjectMeta().SetAnnotations(parent)
	}
	parent[key] = value

	template := r.PodTemplate().GetAnnotations()
	if template == nil {
		template = map[string]string{}
		r.PodTemplate().SetAnnotations(template)
	}
	template[key] = value
}

// setLabel sets the label on both the resource and the resource's
// PodTemplateSpec
func setLabel(r Resource, key, value string) {
	parent := r.GetObjectMeta().GetLabels()
	if parent == nil {
		parent = map[string]string{}
		r.GetObjectMeta().SetLabels(parent)
	}
	parent[key] = value

	template := r.PodTemplate().GetLabels()
	if template == nil {
		template = map[string]string{}
		r.PodTemplate().SetLabels(template)
	}
	template[key] = value
}
