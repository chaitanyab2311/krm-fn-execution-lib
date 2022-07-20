package fn

var runnerErr error

type Runner struct {
	executeFn executeFn
}

func NewRunner() RunnerBuilder {
	return &Runner{
		executeFn: executeFn{},
	}
}

func (r *Runner) WithInput(input []byte) RunnerBuilder {
	err := r.executeFn.addInput(input)
	if err != nil {
		runnerErr = err
	}
	return r
}

func (r *Runner) WithFunctions(function ...Function) RunnerBuilder {
	err := r.executeFn.addFunctions(function...)
	if err != nil {
		runnerErr = err
	}
	return r
}

func (r *Runner) Build() (FunctionRunner, error) {
	return &r.executeFn, runnerErr
}
