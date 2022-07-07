package fn

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type ConfigFile struct {
	yaml.ResourceMeta `yaml:",inline" json:",inline"`

	Function `yaml:"function" json:"function"`
}

// Function specifies a KRM function.
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
