package fn

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/fn/runtime/runtimeutil"
	"sigs.k8s.io/kustomize/kyaml/runfn"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type ExecuteFn struct {
	runfn      runfn.RunFns
	configFile ConfigFile
	Output     bytes.Buffer
}

func (e *ExecuteFn) Execute(inputResource []byte, functionConfig string) ([]byte, error) {
	// parse the function config file
	if err := e.parseConfigFile(functionConfig); err != nil {
		return nil, errors.Wrap(err)
	}
	// set function input
	input := io.Reader(bytes.NewBuffer(inputResource))

	// get the functions to run
	functions, err := e.getFunctions()
	isExec := e.configFile.Exec != ""
	wd, err := "/opt"

	err = runfn.RunFns{
		Input:      input,
		Output:     &e.Output,
		Functions:  functions,
		EnableExec: isExec,
		WorkingDir: wd,
	}.Execute()

	if err != nil {
		return nil, err
	}
	/*	if err := WriteOutput("transform", e.Output.String()); err != nil {
		return nil, err
	}*/
	return e.Output.Bytes(), nil
}

func (e *ExecuteFn) parseConfigFile(functionConfig string) error {
	configBytes, err := ioutil.ReadFile(functionConfig)
	config := ConfigFile{}

	if yaml.Unmarshal(configBytes, &config) != nil {
		return errors.Wrap(err)
	}
	e.configFile = config
	return nil
}

// getFunctions parses the commandline flags and arguments into explicit
// Functions to run.
func (e *ExecuteFn) getFunctions() ([]*yaml.RNode, error) {
	if e.configFile.Image == "" && e.configFile.Exec == "" {
		return nil, nil
	}

	res, err := e.buildFnConfigResource()
	if err != nil {
		return nil, err
	}

	// create the function spec to set as an annotation
	var fnAnnotation *yaml.RNode
	if e.configFile.Image != "" {
		fnAnnotation, err = e.getFnAnnotationForImage()
	} else {
		fnAnnotation, err = e.getFnAnnotationForExec()
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

	return []*yaml.RNode{res}, nil
}

func (e *ExecuteFn) buildFnConfigResource() (*yaml.RNode, error) {
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
	for key, value := range e.configFile.ConfigMap {
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

func (e *ExecuteFn) getFnAnnotationForExec() (*yaml.RNode, error) {
	fn, err := yaml.Parse(`exec: {}`)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	path, err := filepath.Abs(e.configFile.Exec)
	fmt.Printf("path: %s\n", path)
	if err = fn.PipeE(
		yaml.Lookup("exec"),
		yaml.SetField("path", yaml.NewScalarRNode(path))); err != nil {
		return nil, errors.Wrap(err)
	}
	return fn, nil
}

func (e *ExecuteFn) getFnAnnotationForImage() (*yaml.RNode, error) {
	if err := ValidateFunctionImageURL(e.configFile.Image); err != nil {
		return nil, err
	}

	fn, err := yaml.Parse(`container: {}`)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if err = fn.PipeE(
		yaml.Lookup("container"),
		yaml.SetField("image", yaml.NewScalarRNode(e.configFile.Image))); err != nil {
		return nil, errors.Wrap(err)
	}
	return fn, nil
}
