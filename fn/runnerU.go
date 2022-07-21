package fn

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

var runnerUErr error

type RunnerU struct {
	executeFn executeFnU
}

func (r RunnerU) WithInput(bytes []byte) RunnerBuilderU {
	err := r.executeFn.addInput(bytes)
	if err != nil {
		runnerUErr = err
	}
	return r
}

func (r RunnerU) WithFunctions(function ...Function) RunnerBuilderU {
	err := r.executeFn.addFunctions(function...)
	if err != nil {
		runnerUErr = err
	}
	return r
}

func (r RunnerU) Build() (FunctionRunnerU, error) {
	return &r.executeFn, runnerUErr
}

func (r RunnerU) WithInputs(unstructured ...unstructured.Unstructured) RunnerBuilderU {
	err := r.executeFn.addInputs(unstructured...)
	if err != nil {
		runnerUErr = err
	}
	return r
}

func NewRunnerU() RunnerBuilderU {
	return &RunnerU{
		executeFn: executeFnU{},
	}
}
