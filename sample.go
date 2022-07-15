package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"krm-fn-execution-lib/fn"
)

func main() {
	data, err := ioutil.ReadFile("testdata/service.yaml")
	if err != nil {
		panic(err)
	}
	input := io.Reader(bytes.NewBuffer(data))
	output := bytes.Buffer{}

	executeFn := fn.ExecuteFn{
		Input:          input,
		FunctionConfig: GetFnConfig(),
		Output:         &output,
	}
	if err = executeFn.Execute(); err != nil {
		panic(err)
	}
	fmt.Printf("Output of command is: \n%s", output.String())
}

func GetFnConfig() fn.FunctionConfig {
	functions := []fn.Function{
		{
			Exec: "testdata/clean-metadata",
		},
		{
			Image: "gcr.io/kpt-fn/set-labels:v0.1",
			ConfigMap: map[string]string{
				"env":      "dev",
				"app-name": "my-app",
			},
		},
	}
	return fn.FunctionConfig{Functions: functions}
}
