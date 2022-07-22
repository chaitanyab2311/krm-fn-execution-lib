package fn

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
)

type ResourceList struct {
	// Items is the ResourceList.items input and output value.
	//
	// e.g. given the function input:
	//
	//    items:
	//    - kind: Deployment
	//      ...
	//    - kind: Service
	//      ...
	//
	// Items will be a slice containing the Deployment and Service resources
	// Mutating functions will alter this field during processing.
	// This field is required.
	Items []runtime.Object

	// Results is ResourceList.results output value.
	// Validating functions can optionally use this field to communicate structured
	// validation error data to downstream functions.
	Results framework.Results
}

// Function specifies a KRM function to run.
type Function struct {
	// `Image` specifies the function container image.
	//	image: gcr.io/kpt-fn/set-labels
	Image string `yaml:"image,omitempty" json:"image,omitempty"`

	// Exec specifies the function binary executable.
	// The executable can be fully qualified, or it must exist in the $PATH e.g:
	//
	// 	 exec: set-namespace
	// 	 exec: /usr/local/bin/my-custom-fn
	Exec string `yaml:"exec,omitempty" json:"exec,omitempty"`

	// `ConfigMap` is a convenient way to specify a function config of kind ConfigMap.
	ConfigMap map[string]string `yaml:"configMap,omitempty" json:"configMap,omitempty"`
}

type RunnerBuilder interface {
	WithInput([]byte) RunnerBuilder
	WithInputs(...runtime.Object) RunnerBuilder
	WithFunctions(...Function) RunnerBuilder
	WhereExecWorkingDir(string) RunnerBuilder
	Build() (FunctionRunner, error)
}

type FunctionRunner interface {
	Execute() (ResourceList, error)
}
