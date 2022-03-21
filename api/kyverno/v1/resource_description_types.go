package v1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ResourceDescription contains criteria used to match resources.
type ResourceDescription struct {
	// Kinds is a list of resource kinds.
	// +optional
	Kinds []string `json:"kinds,omitempty" yaml:"kinds,omitempty"`

	// Name is the name of the resource. The name supports wildcard characters
	// "*" (matches zero or many characters) and "?" (at least one character).
	// +optional
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Names are the names of the resources. Each name supports wildcard characters
	// "*" (matches zero or many characters) and "?" (at least one character).
	// NOTE: "Name" is being deprecated in favor of "Names".
	// +optional
	Names []string `json:"names,omitempty" yaml:"names,omitempty"`

	// Namespaces is a list of namespaces names. Each name supports wildcard characters
	// "*" (matches zero or many characters) and "?" (at least one character).
	// +optional
	Namespaces []string `json:"namespaces,omitempty" yaml:"namespaces,omitempty"`

	// Annotations is a  map of annotations (key-value pairs of type string). Annotation keys
	// and values support the wildcard characters "*" (matches zero or many characters) and
	// "?" (matches at least one character).
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`

	// Selector is a label selector. Label keys and values in `matchLabels` support the wildcard
	// characters `*` (matches zero or many characters) and `?` (matches one character).
	// Wildcards allows writing label selectors like ["storage.k8s.io/*": "*"]. Note that
	// using ["*" : "*"] matches any key and value but does not match an empty label set.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty" yaml:"selector,omitempty"`

	// NamespaceSelector is a label selector for the resource namespace. Label keys and values
	// in `matchLabels` support the wildcard characters `*` (matches zero or many characters)
	// and `?` (matches one character).Wildcards allows writing label selectors like
	// ["storage.k8s.io/*": "*"]. Note that using ["*" : "*"] matches any key and value but
	// does not match an empty label set.
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty" yaml:"namespaceSelector,omitempty"`
}

// Validate implements programmatic validation
func (r *ResourceDescription) Validate(path *field.Path, namespaced bool, clusterResources sets.String) field.ErrorList {
	var errs field.ErrorList
	if namespaced {
		if len(r.Namespaces) > 0 {
			errs = append(errs, field.Forbidden(path.Child("namespaces"), "Filtering namespaces not allowed in namespaced policies"))
		}
		kindsChild := path.Child("kinds")
		for i, kind := range r.Kinds {
			if clusterResources.Has(kind) {
				errs = append(errs, field.Forbidden(kindsChild.Index(i), fmt.Sprintf("Cluster wide resource '%s' not allowed in namespaced policy", kind)))
			}
		}
	}
	return errs
}
