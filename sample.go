package main

import (
	"fmt"
	"io/ioutil"
	"krm-fn-execution-lib/fn"
)

func main() {
	input, err := ioutil.ReadFile("testdata/service.yaml")
	if err != nil {
		panic(err)
	}

	executeFn := fn.ExecuteFn{}
	output, err := executeFn.Execute(input, "testdata/fnconfig.yaml")
	if err != nil {
		panic(err)
		return
	}

	fmt.Printf("Output of command is: \n%s", string(output))
}
