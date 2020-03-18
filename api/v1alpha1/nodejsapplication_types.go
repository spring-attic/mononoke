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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	NodeJsApplicationLabelKey = GroupVersion.Group + "/nodejs-application"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NodeJsApplicationSpec defines the desired state of NodeJsApplication
type NodeJsApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Template pod
	// +optional
	Template *corev1.PodTemplateSpec `json:"template,omitempty"`
}

// NodeJsApplicationStatus defines the observed state of NodeJsApplication
type NodeJsApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	apis.Status `json:",inline"`

	// AppliedOpinions lists opinions applied to the application
	AppliedOpinions []string `json:"appliedOpinions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:categories="mononoke"
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Reason",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].reason`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// NodeJsApplication is the Schema for the nodejsapplications API
type NodeJsApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeJsApplicationSpec   `json:"spec,omitempty"`
	Status NodeJsApplicationStatus `json:"status,omitempty"`
}

func (a *NodeJsApplication) GetStatus() apis.ResourceStatus {
	return &a.Status
}

// +kubebuilder:object:root=true

// NodeJsApplicationList contains a list of NodeJsApplication
type NodeJsApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeJsApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeJsApplication{}, &NodeJsApplicationList{})
}
