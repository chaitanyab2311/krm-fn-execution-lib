package fn

import (
	"bytes"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"path/filepath"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/fn/runtime/runtimeutil"
	"sigs.k8s.io/kustomize/kyaml/runfn"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type executeFn struct {
	input     []*yaml.RNode
	functions []*yaml.RNode
}

func (e *executeFn) Execute() ([]*yaml.RNode, error) {
	var output []*yaml.RNode
	out, err := e.applyFn()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	output, err = GetRNode(out.String())
	return output, errors.Wrap(err)
}

func (e *executeFn) Run() (unstructured.UnstructuredList, error) {
	var output unstructured.UnstructuredList
	out, err := e.applyFn()
	if err != nil {
		return output, errors.Wrap(err)
	}
	output, err = GetUnstructured(out.String())
	return output, nil
}

func (e *executeFn) addInput(input []byte) error {
	nodes, err := GetRNode(string(input))
	if err != nil {
		return errors.Wrap(err)
	}
	e.input = append(e.input, nodes...)
	return nil
}

func (e *executeFn) addInputs(inputs ...*yaml.RNode) error {
	e.input = append(e.input, inputs...)
	return nil
}

func (e *executeFn) addFunctions(functions ...Function) error {
	functionConfig, err := getFunctionConfig(functions)
	if err != nil {
		return errors.Wrap(err)
	}
	e.functions = append(e.functions, functionConfig...)
	return nil
}

// getFunctionsToExecute parses the explicit functions to run.
func getFunctionConfig(functions []Function) ([]*yaml.RNode, error) {
	var functionConfig []*yaml.RNode
	for _, fn := range functions {
		res, err := buildFnConfigResource(fn)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		// create the function spec to set as an annotation
		var fnAnnotation *yaml.RNode
		if fn.Image != "" {
			fnAnnotation, err = getFnAnnotationForImage(fn)
		} else {
			fnAnnotation, err = getFnAnnotationForExec(fn)
		}

		if err != nil {
			return nil, errors.Wrap(err)
		}

		// set the function annotation on the function config, so that it is parsed by RunFns
		value, err := fnAnnotation.String()
		if err != nil {
			return nil, errors.Wrap(err)
		}

		if err = res.PipeE(
			yaml.LookupCreate(yaml.MappingNode, "metadata", "annotations"),
			yaml.SetField(runtimeutil.FunctionAnnotationKey, yaml.NewScalarRNode(value))); err != nil {
			return nil, errors.Wrap(err)
		}
		functionConfig = append(functionConfig, res)
	}
	return functionConfig, nil
}

func buildFnConfigResource(function Function) (*yaml.RNode, error) {
	if function.Image == "" && function.Exec == "" {
		return nil, fmt.Errorf("function must have either image or exec, none specified")
	}
	// create the function config
	rc, err := yaml.Parse(`
metadata:
  name: function-input
data: {}
`)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// default the function config kind to ConfigMap, this may be overridden
	var kind = "ConfigMap"
	var version = "v1"

	// populate the function config with data.
	dataField, err := rc.Pipe(yaml.Lookup("data"))
	for key, value := range function.ConfigMap {
		err := dataField.PipeE(
			yaml.FieldSetter{Name: key, Value: yaml.NewStringRNode(value), OverrideStyle: true})
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	if err = rc.PipeE(yaml.SetField("kind", yaml.NewScalarRNode(kind))); err != nil {
		return nil, errors.Wrap(err)
	}
	if err = rc.PipeE(yaml.SetField("apiVersion", yaml.NewScalarRNode(version))); err != nil {
		return nil, errors.Wrap(err)
	}
	return rc, nil
}

func getFnAnnotationForExec(function Function) (*yaml.RNode, error) {
	fn, err := yaml.Parse(`exec: {}`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	path, err := filepath.Abs(function.Exec)
	if err = fn.PipeE(
		yaml.Lookup("exec"),
		yaml.SetField("path", yaml.NewScalarRNode(path))); err != nil {
		return nil, errors.Wrap(err)
	}
	return fn, nil
}

func getFnAnnotationForImage(function Function) (*yaml.RNode, error) {
	if err := ValidateFunctionImageURL(function.Image); err != nil {
		return nil, errors.Wrap(err)
	}

	fn, err := yaml.Parse(`container: {}`)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if err = fn.PipeE(
		yaml.Lookup("container"),
		yaml.SetField("image", yaml.NewScalarRNode(function.Image))); err != nil {
		return nil, errors.Wrap(err)
	}
	return fn, nil
}

func (e *executeFn) applyFn() (bytes.Buffer, error) {
	out := bytes.Buffer{}

	wd, err := os.Getwd()
	if err != nil {
		return out, errors.Wrap(err)
	}

	input, err := ReadInput(e.input)
	if err != nil {
		return out, errors.Wrap(err)
	}

	err = runfn.RunFns{
		Input:      input,
		Output:     &out,
		Functions:  e.functions,
		EnableExec: true,
		WorkingDir: wd,
	}.Execute()
	return out, errors.Wrap(err)
}
