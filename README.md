# krm-fn-execution-lib
Execute KRM Function containerized images and binaries


## Library Input
Source File: 
- [executefn.go](https://github.com/MundraAnkur/krm-fn-execution-lib/blob/main/fn/executefn.go)
- [types.go](https://github.com/MundraAnkur/krm-fn-execution-lib/blob/main/fn/types.go)

Function call
```
Execute(inputResource []byte, functionConfig string) ([]byte, error)
```

##### Input
   
 | Parameter | Description |
 | --- | ----------- |
 | inputResource | Resources on which function will execute |
 | functionConfig | Path to file which contains function related configuration|
 
 Example config:
 ```
 apiVersion: v1
kind: FunctionConfig
metadata:
  name: fn-config1
function:
  image: gcr.io/kpt-fn/set-labels:v0.1
  configMap:
    app-name: todolist
    env: qa
 ```
 
 ```
 apiVersion: v1
kind: FunctionConfig
metadata:
  name: fn-exec-config1
function:
  exec: testdata/clean-metadata
 ```
  ##### Output
```
input, _ := ioutil.ReadFile("testdata/service.yaml")
executeFn := fn.ExecuteFn{}
output, err := executeFn.Execute(input, "testdata/fnconfig.yaml")
```

 | Parameter | Description |
 | --- | ----------- |
 | outputResource | Transformed resources after function execution|
