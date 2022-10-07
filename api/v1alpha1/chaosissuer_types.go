/*
Copyright 2022 CMU-SV.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChaosIssuerSpec defines the desired state of ChaosIssuer
type ChaosIssuerSpec struct {

	// URL is the base URL for the step certificates instance.
	URL string `json:"url"`
	// A reference to a Secret in the same namespace as the referent. If the
	// referent is a ClusterIssuer, the reference instead refers to the resource
	// with the given name in the configured 'cluster resource namespace', which
	// is set as a flag on the controller component (and defaults to the
	// namespace that the controller runs in).
	AuthSecretName string `json:"authSecretName"`
}

// ChaosIssuerStatus defines the observed state of ChaosIssuer
type ChaosIssuerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	Conditions []ChaosIssuerCondition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChaosIssuer is the Schema for the chaosissuers API
type ChaosIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChaosIssuerSpec   `json:"spec,omitempty"`
	Status ChaosIssuerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChaosIssuerList contains a list of ChaosIssuer
type ChaosIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChaosIssuer `json:"items"`
}

// IssuerCondition contains condition information for an Issuer.
type IssuerCondition struct {
	// Type of the condition, known values are ('Ready').
	Type IssuerConditionType `json:"type"`

	// Status of the condition, one of ('True', 'False', 'Unknown').
	Status ConditionStatus `json:"status"`

	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	// +optional
	Message string `json:"message,omitempty"`
}

// IssuerConditionType represents an Issuer condition value.
type IssuerConditionType string

const (
	// IssuerConditionReady represents the fact that a given Issuer condition
	// is in ready state and able to issue certificates.
	// If the `status` of this condition is `False`, CertificateRequest controllers
	// should prevent attempts to sign certificates.
	IssuerConditionReady IssuerConditionType = "Ready"
)

// ConditionStatus represents a condition's status.
// +kubebuilder:validation:Enum=True;False;Unknown
type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in
// the condition; "ConditionFalse" means a resource is not in the condition;
// "ConditionUnknown" means kubernetes can't decide if a resource is in the
// condition or not. In the future, we could add other intermediate
// conditions, e.g. ConditionDegraded.
const (
	// ConditionTrue represents the fact that a given condition is true
	ConditionTrue ConditionStatus = "True"

	// ConditionFalse represents the fact that a given condition is false
	ConditionFalse ConditionStatus = "False"

	// ConditionUnknown represents the fact that a given condition is unknown
	ConditionUnknown ConditionStatus = "Unknown"
)

func init() {
	SchemeBuilder.Register(&ChaosIssuer{}, &ChaosIssuerList{})
}
