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
	RailsApplicationConditionReady                              = apis.ConditionReady
	RailsApplicationConditionDeploymentReady apis.ConditionType = "DeploymentReady"
)

var railsappCondSet = apis.NewLivingConditionSet(
	RailsApplicationConditionDeploymentReady,
)

func (rs *RailsApplicationStatus) GetObservedGeneration() int64 {
	return rs.ObservedGeneration
}

func (rs *RailsApplicationStatus) IsReady() bool {
	return railsappCondSet.Manage(rs).IsHappy()
}

func (*RailsApplicationStatus) GetReadyConditionType() apis.ConditionType {
	return RailsApplicationConditionReady
}

func (rs *RailsApplicationStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return railsappCondSet.Manage(rs).GetCondition(t)
}

func (rs *RailsApplicationStatus) InitializeConditions() {
	railsappCondSet.Manage(rs).InitializeConditions()
}

func (rs *RailsApplicationStatus) MarkDeploymentNotOwned(name string) {
	railsappCondSet.Manage(rs).MarkFalse(RailsApplicationConditionDeploymentReady, "NotOwned", "There is an existing Deployer %q that the RailsApplication does not own.", name)
}

func (rs *RailsApplicationStatus) PropagateDeploymentStatus(crs *appsv1.DeploymentStatus) {
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
		// DeploymentAvailable is False while progressing, avoid reporting RailsApplicationConditionReady as False
		railsappCondSet.Manage(rs).MarkUnknown(RailsApplicationConditionDeploymentReady, progressing.Reason, progressing.Message)
		return
	}
	switch {
	case available.Status == corev1.ConditionUnknown:
		railsappCondSet.Manage(rs).MarkUnknown(RailsApplicationConditionDeploymentReady, available.Reason, available.Message)
	case available.Status == corev1.ConditionTrue:
		railsappCondSet.Manage(rs).MarkTrue(RailsApplicationConditionDeploymentReady)
	case available.Status == corev1.ConditionFalse:
		railsappCondSet.Manage(rs).MarkFalse(RailsApplicationConditionDeploymentReady, available.Reason, available.Message)
	}
}
