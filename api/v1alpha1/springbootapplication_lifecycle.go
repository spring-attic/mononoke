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

package v1alpha1

import (
	"github.com/projectriff/system/pkg/apis"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	SpringBootApplicationConditionReady                              = apis.ConditionReady
	SpringBootApplicationConditionDeploymentReady apis.ConditionType = "DeploymentReady"
)

var springbootappCondSet = apis.NewLivingConditionSet(
	SpringBootApplicationConditionDeploymentReady,
)

func (rs *SpringBootApplicationStatus) GetObservedGeneration() int64 {
	return rs.ObservedGeneration
}

func (rs *SpringBootApplicationStatus) IsReady() bool {
	return springbootappCondSet.Manage(rs).IsHappy()
}

func (*SpringBootApplicationStatus) GetReadyConditionType() apis.ConditionType {
	return SpringBootApplicationConditionReady
}

func (rs *SpringBootApplicationStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return springbootappCondSet.Manage(rs).GetCondition(t)
}

func (rs *SpringBootApplicationStatus) InitializeConditions() {
	springbootappCondSet.Manage(rs).InitializeConditions()
}

func (rs *SpringBootApplicationStatus) MarkDeploymentNotOwned(name string) {
	springbootappCondSet.Manage(rs).MarkFalse(SpringBootApplicationConditionDeploymentReady, "NotOwned", "There is an existing Deployer %q that the SpringBootApplication does not own.", name)
}

func (rs *SpringBootApplicationStatus) PropagateDeploymentStatus(crs *appsv1.DeploymentStatus) {
	var available, progressing *appsv1.DeploymentCondition
	for i := range crs.Conditions {
		switch crs.Conditions[i].Type {
		case appsv1.DeploymentAvailable:
			available = &crs.Conditions[i]
		case appsv1.DeploymentProgressing:
			progressing = &crs.Conditions[i]
		}
	}
	if available == nil || progressing == nil {
		return
	}
	if progressing.Status == corev1.ConditionTrue && available.Status == corev1.ConditionFalse {
		// DeploymentAvailable is False while progressing, avoid reporting SpringBootApplicationConditionReady as False
		springbootappCondSet.Manage(rs).MarkUnknown(SpringBootApplicationConditionDeploymentReady, progressing.Reason, progressing.Message)
		return
	}
	switch {
	case available.Status == corev1.ConditionUnknown:
		springbootappCondSet.Manage(rs).MarkUnknown(SpringBootApplicationConditionDeploymentReady, available.Reason, available.Message)
	case available.Status == corev1.ConditionTrue:
		springbootappCondSet.Manage(rs).MarkTrue(SpringBootApplicationConditionDeploymentReady)
	case available.Status == corev1.ConditionFalse:
		springbootappCondSet.Manage(rs).MarkFalse(SpringBootApplicationConditionDeploymentReady, available.Reason, available.Message)
	}
}
