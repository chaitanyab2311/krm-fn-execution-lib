# krm-fn-execution-lib
Execute KRM Function containerized images and binaries


## Library Input
Source File: 
- [executefn.go](https://github.com/MundraAnkur/krm-fn-execution-lib/blob/main/fn/executefn.go)
- [types.go](https://github.com/MundraAnkur/krm-fn-execution-lib/blob/main/fn/types.go)

Function call
```
executeFn := fn.ExecuteFn{
		Input:          inputResource,
		FunctionConfig: fnConfig,
		Output:         &outputResource,
	}

err = executeFn.Execute()
```

##### Input
   
 | Parameter | Description |
 | --- | ----------- |
 | inputResource | Resources on which function will execute |
 | fnConfig | function configuration|
 
 Example config:
 ```
function := fn.Function {
            Image: "gcr.io/kpt-fn/set-labels:v0.1",
            ConfigMap: map[string]string{
               "env":      "dev",
               "app-name": "my-app",
            },
         }
 ```
 
 ```
 function := fn.Function { Exec: "testdata/clean-metadata"}
 ```
  ##### Output
 | Parameter | Description |
 | --- | ----------- |
 | outputResource | Transformed resources after function execution|
