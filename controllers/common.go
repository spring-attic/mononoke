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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func FindTargetContainer(target *intstr.IntOrString, template *corev1.PodTemplateSpec) (string, int, error) {
	switch target.Type {
	case intstr.Int:
		idx := int(target.IntVal)
		if l := len(template.Spec.Containers); idx < l {
			c := template.Spec.Containers[idx]
			return c.Name, idx, nil
		}
	case intstr.String:
		name := target.StrVal
		for i, c := range template.Spec.Containers {
			if c.Name == name {
				return name, i, nil
			}
		}
	}

	return "", 0, fmt.Errorf("Unable to find container %q", target)
}
