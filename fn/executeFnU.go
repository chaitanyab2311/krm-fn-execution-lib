package fn

import (
	"bytes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/runfn"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type executeFnU struct {
	input     []unstructured.Unstructured
	functions []*yaml.RNode
}

func (e *executeFnU) Execute() (unstructured.UnstructuredList, error) {
	var output unstructured.UnstructuredList
	out, err := e.applyFn()
	if err != nil {
		return output, errors.Wrap(err)
	}
	output, err = GetUnstructured(out.String())
	return output, err
}

func (e *executeFnU) addInput(input []byte) error {
	nodes, err := GetUnstructured(string(input))
	if err != nil {
		return errors.Wrap(err)
	}
	e.input = append(e.input, nodes.Items...)
	return nil
}

func (e *executeFnU) addInputs(inputs ...unstructured.Unstructured) error {
	e.input = append(e.input, inputs...)
	return nil
}

func (e *executeFnU) addFunctions(functions ...Function) error {
	functionConfig, err := getFunctionConfig(functions)
	if err != nil {
		return errors.Wrap(err)
	}
	e.functions = append(e.functions, functionConfig...)
	return nil
}

func (e *executeFnU) applyFn() (bytes.Buffer, error) {
	out := bytes.Buffer{}

	wd, err := os.Getwd()
	if err != nil {
		return out, errors.Wrap(err)
	}

	input, err := ReadUnstructuredInput(e.input)
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
