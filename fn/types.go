package fn

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

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

type FunctionRunner interface {
	Execute() ([]*yaml.RNode, error)
}

type RunnerBuilder interface {
	WithInput([]byte) RunnerBuilder
	WithFunctions(...Function) RunnerBuilder
	Build() (FunctionRunner, error)
}

type RunnerBuilderU interface {
	WithInput([]byte) RunnerBuilderU
	WithInputs(...unstructured.Unstructured) RunnerBuilderU
	WithFunctions(...Function) RunnerBuilderU
	Build() (FunctionRunnerU, error)
}

type FunctionRunnerU interface {
	Execute() (unstructured.UnstructuredList, error)
}
