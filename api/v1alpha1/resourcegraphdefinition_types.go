// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package v1alpha1

import (
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

const (
	// DefaultServiceAccountKey is the key to use for the default service account
	// in the serviceAccounts map.
	DefaultServiceAccountKey = "*"
)

// ResourceGraphDefinitionSpec defines the desired state of ResourceGraphDefinition
type ResourceGraphDefinitionSpec struct {
	// The schema of the resourcegraphdefinition, which includes the
	// apiVersion, kind, spec, status, types, and some validation
	// rules.
	//
	// +kubebuilder:validation:Required
	Schema *Schema `json:"schema,omitempty"`
	// The resources that are part of the resourcegraphdefinition.
	//
	// +kubebuilder:validation:Optional
	Resources []*Resource `json:"resources,omitempty"`
	// ServiceAccount configuration for controller impersonation.
	// Key is the namespace, value is the service account name to use.
	// Special key "*" defines the default service account for any
	// namespace not explicitly mapped.
	//
	// +kubebuilder:validation:Optional
	DefaultServiceAccounts map[string]string `json:"defaultServiceAccounts,omitempty"`
}

// Schema represents the attributes that define an instance of
// a resourcegraphdefinition.
type Schema struct {
	// The kind of the resourcegraphdefinition. This is used to generate
	// and create the CRD for the resourcegraphdefinition.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Z][a-zA-Z0-9]{0,62}$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="kind is immutable"
	Kind string `json:"kind,omitempty"`
	// The APIVersion of the resourcegraphdefinition. This is used to generate
	// and create the CRD for the resourcegraphdefinition.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^v[0-9]+(alpha[0-9]+|beta[0-9]+)?$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="apiVersion is immutable"
	APIVersion string `json:"apiVersion,omitempty"`
	// The group of the resourcegraphdefinition. This is used to set the API group
	// of the generated CRD. If omitted, it defaults to "kro.run".
	//
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="kro.run"
	Group string `json:"group,omitempty"`
	// The spec of the resourcegraphdefinition. Typically, this is the spec of
	// the CRD that the resourcegraphdefinition is managing. This is adhering
	// to the SimpleSchema spec
	Spec runtime.RawExtension `json:"spec,omitempty"`

	// Types is a map of custom type definitions. These can be used in the spec
	// of the resourcegraphdefinition. Each type definition is also adhering to
	// the SimpleSchema spec.
	Types runtime.RawExtension `json:"types,omitempty"`

	// The status of the resourcegraphdefinition. This is the status of the CRD
	// that the resourcegraphdefinition is managing. This is adhering to the
	// SimpleSchema spec.
	Status runtime.RawExtension `json:"status,omitempty"`
	// Validation is a list of validation rules that are applied to the
	// resourcegraphdefinition.
	Validation []Validation `json:"validation,omitempty"`
	// AdditionalPrinterColumns defines additional printer columns
	// that will be passed down to the created CRD. If set, no
	// default printer columns will be added to the created CRD,
	// and if default printer columns need to be retained, they
	// need to be added explicitly.
	//
	// +kubebuilder:validation:Optional
	AdditionalPrinterColumns []extv1.CustomResourceColumnDefinition `json:"additionalPrinterColumns,omitempty"`
}

type Validation struct {
	Expression string `json:"expression,omitempty"`
	Message    string `json:"message,omitempty"`
}

type ExternalRefMetadata struct {
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace,omitempty"`
}

// ExternalRef is a reference to an external resource.
// It allows the user to specify the Kind, Version, Name and Namespace of the resource
// to be read and used in the Graph.
type ExternalRef struct {
	// +kubebuilder:validation:Required
	APIVersion string `json:"apiVersion"`
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
	// +kubebuilder:validation:Required
	Metadata ExternalRefMetadata `json:"metadata"`
}

// +kubebuilder:validation:XValidation:rule="(has(self.template) && !has(self.externalRef)) || (!has(self.template) && has(self.externalRef))",message="exactly one of template or externalRef must be provided"
type Resource struct {
	// +kubebuilder:validation:Required
	ID string `json:"id,omitempty"`
	// +kubebuilder:validation:Optional
	Template runtime.RawExtension `json:"template,omitempty"`
	// +kubebuilder:validation:Optional
	ExternalRef *ExternalRef `json:"externalRef,omitempty"`
	// +kubebuilder:validation:Optional
	ReadyWhen []string `json:"readyWhen,omitempty"`
	// +kubebuilder:validation:Optional
	IncludeWhen []string `json:"includeWhen,omitempty"`
}

// ResourceGraphDefinitionState defines the state of the resource graph definition.
type ResourceGraphDefinitionState string

const (
	// ResourceGraphDefinitionStateActive represents the active state of the resource definition.
	ResourceGraphDefinitionStateActive ResourceGraphDefinitionState = "Active"
	// ResourceGraphDefinitionStateInactive represents the inactive state of the resource graph definition
	ResourceGraphDefinitionStateInactive ResourceGraphDefinitionState = "Inactive"
)

// ResourceGraphDefinitionStatus defines the observed state of ResourceGraphDefinition
type ResourceGraphDefinitionStatus struct {
	// State is the state of the resourcegraphdefinition
	State ResourceGraphDefinitionState `json:"state,omitempty"`
	// TopologicalOrder is the topological order of the resourcegraphdefinition graph
	TopologicalOrder []string `json:"topologicalOrder,omitempty"`
	// Conditions represent the latest available observations of an object's state
	Conditions Conditions `json:"conditions,omitempty"`
	// Resources represents the resources, and their information (dependencies for now)
	Resources []ResourceInformation `json:"resources,omitempty"`
}

// ResourceInformation defines the information about a resource
// in the resourcegraphdefinition
type ResourceInformation struct {
	// ID represents the id of the resources we're providing information for
	ID string `json:"id,omitempty"`
	// Dependencies represents the resource dependencies of a resource graph definition
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

// Dependency defines the dependency a resource has observed
// from the resources it points to based on expressions
type Dependency struct {
	// ID represents the id of the dependency resource
	ID string `json:"id,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="APIVERSION",type=string,priority=0,JSONPath=`.spec.schema.apiVersion`
// +kubebuilder:printcolumn:name="KIND",type=string,priority=0,JSONPath=`.spec.schema.kind`
// +kubebuilder:printcolumn:name="STATE",type=string,priority=0,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="TOPOLOGICALORDER",type=string,priority=1,JSONPath=`.status.topologicalOrder`
// +kubebuilder:printcolumn:name="AGE",type="date",priority=0,JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:shortName=rgd,scope=Cluster

// ResourceGraphDefinition is the Schema for the resourcegraphdefinitions API
type ResourceGraphDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceGraphDefinitionSpec   `json:"spec,omitempty"`
	Status ResourceGraphDefinitionStatus `json:"status,omitempty"`
}

func (o *ResourceGraphDefinition) GetConditions() []Condition {
	return o.Status.Conditions
}

func (o *ResourceGraphDefinition) SetConditions(conditions []Condition) {
	o.Status.Conditions = conditions
}

//+kubebuilder:object:root=true

// ResourceGraphDefinitionList contains a list of ResourceGraphDefinition
type ResourceGraphDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceGraphDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceGraphDefinition{}, &ResourceGraphDefinitionList{})
}
