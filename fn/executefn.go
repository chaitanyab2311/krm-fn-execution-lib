package fn

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/fn/runtime/runtimeutil"
	"sigs.k8s.io/kustomize/kyaml/runfn"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type ExecuteFn struct {
	Input          io.Reader
	Output         io.Writer
	FunctionConfig FunctionConfig
}

func (e *ExecuteFn) Execute() error {
	// get the functions to run
	functions, err := e.getFunctions()
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if e.Output == nil {
		e.Output = io.Writer(bytes.NewBuffer([]byte{}))
	}

	err = runfn.RunFns{
		Input:      e.Input,
		Output:     e.Output,
		Functions:  functions,
		EnableExec: true,
		WorkingDir: wd,
		ResultsDir: wd,
	}.Execute()

	return err
}

// getFunctions parses the explicit Functions to run.
func (e *ExecuteFn) getFunctions() ([]*yaml.RNode, error) {
	if len(e.FunctionConfig.Functions) == 0 {
		return nil, fmt.Errorf("no functions to run")
	}

	var functions []*yaml.RNode
	for _, fn := range e.FunctionConfig.Functions {
		res, err := buildFnConfigResource(fn)
		if err != nil {
			return nil, err
		}

		// create the function spec to set as an annotation
		var fnAnnotation *yaml.RNode
		if fn.Image != "" {
			fnAnnotation, err = getFnAnnotationForImage(fn)
		} else {
			fnAnnotation, err = getFnAnnotationForExec(fn)
		}

		if err != nil {
			return nil, err
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
		functions = append(functions, res)
	}

	return functions, nil
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
		return nil, err
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
			return nil, err
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
		return nil, err
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
