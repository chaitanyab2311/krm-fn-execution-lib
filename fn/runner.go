package fn

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
)

var runnerErr error

type Runner struct {
	executeFn executeFn
}

func NewRunner() RunnerBuilder {
	return &Runner{
		executeFn: executeFn{},
	}
}

func (r Runner) WithInput(bytes []byte) RunnerBuilder {
	err := r.executeFn.addInput(bytes)
	appendError(err)
	return r
}

func (r Runner) WithFunctions(function ...Function) RunnerBuilder {
	err := r.executeFn.addFunctions(function...)
	appendError(err)
	return r
}

func (r Runner) WithInputs(objects ...runtime.Object) RunnerBuilder {
	err := r.executeFn.addInputs(objects...)
	appendError(err)
	return r
}

func (r Runner) WhereExecWorkingDir(dir string) RunnerBuilder {
	err := r.executeFn.setExecWorkingDir(dir)
	appendError(err)
	return r
}

func (r Runner) Build() (FunctionRunner, error) {
	return &r.executeFn, runnerErr
}

func appendError(err error) {
	if err != nil {
		if runnerErr == nil {
			runnerErr = err
		} else {
			runnerErr = fmt.Errorf("%v\n%v", runnerErr, err)
		}
	}
}
