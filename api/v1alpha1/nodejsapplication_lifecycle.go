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
	NodeJsApplicationConditionReady                              = apis.ConditionReady
	NodeJsApplicationConditionDeploymentReady apis.ConditionType = "DeploymentReady"
)

var nodejsappCondSet = apis.NewLivingConditionSet(
	NodeJsApplicationConditionDeploymentReady,
)

func (rs *NodeJsApplicationStatus) GetObservedGeneration() int64 {
	return rs.ObservedGeneration
}

func (rs *NodeJsApplicationStatus) IsReady() bool {
	return nodejsappCondSet.Manage(rs).IsHappy()
}

func (*NodeJsApplicationStatus) GetReadyConditionType() apis.ConditionType {
	return NodeJsApplicationConditionReady
}

func (rs *NodeJsApplicationStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return nodejsappCondSet.Manage(rs).GetCondition(t)
}

func (rs *NodeJsApplicationStatus) InitializeConditions() {
	nodejsappCondSet.Manage(rs).InitializeConditions()
}

func (rs *NodeJsApplicationStatus) MarkDeploymentNotOwned(name string) {
	nodejsappCondSet.Manage(rs).MarkFalse(NodeJsApplicationConditionDeploymentReady, "NotOwned", "There is an existing Deployer %q that the NodeJsApplication does not own.", name)
}

func (rs *NodeJsApplicationStatus) PropagateDeploymentStatus(crs *appsv1.DeploymentStatus) {
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
		// DeploymentAvailable is False while progressing, avoid reporting NodeJsApplicationConditionReady as False
		nodejsappCondSet.Manage(rs).MarkUnknown(NodeJsApplicationConditionDeploymentReady, progressing.Reason, progressing.Message)
		return
	}
	switch {
	case available.Status == corev1.ConditionUnknown:
		nodejsappCondSet.Manage(rs).MarkUnknown(NodeJsApplicationConditionDeploymentReady, available.Reason, available.Message)
	case available.Status == corev1.ConditionTrue:
		nodejsappCondSet.Manage(rs).MarkTrue(NodeJsApplicationConditionDeploymentReady)
	case available.Status == corev1.ConditionFalse:
		nodejsappCondSet.Manage(rs).MarkFalse(NodeJsApplicationConditionDeploymentReady, available.Reason, available.Message)
	}
}
