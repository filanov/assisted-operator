/*
Copyright 2020.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AssistedServiceSpec defines the desired state of AssistedService
type AssistedServiceSpec struct {
	// URL is a location of an assisted-service instance
	URL   string `json:"url"`
	Token string `json:"token"`
}

type SyncState string

const (
	SyncStateOK    SyncState = "OK"
	SyncStateError SyncState = "error"
)

// AssistedServiceStatus defines the observed state of AssistedService
type AssistedServiceStatus struct {
	// SyncState holds the status of the connection to the service
	// +kubebuilder:validation:Enum="";OK;error
	SyncState SyncState `json:"syncState"`

	// LastUpdated identifies when this status was last observed.
	// +optional
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AssistedService is the Schema for the assistedservices API
type AssistedService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AssistedServiceSpec   `json:"spec,omitempty"`
	Status AssistedServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AssistedServiceList contains a list of AssistedService
type AssistedServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AssistedService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AssistedService{}, &AssistedServiceList{})
}
