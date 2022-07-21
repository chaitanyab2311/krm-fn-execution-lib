package fn

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestExecuteFnRunner(t *testing.T) {
	fnRunner, err := getFnRunner()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	OutRl, err := fnRunner.Execute()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	expectedLabels := map[string]string{
		"app-name": "my-app",
		"env":      "dev",
		"tier":     "frontend",
		"app":      "guestbook",
	}
	assert.EqualValues(t, expectedLabels, OutRl[1].GetLabels())
}

func getFnRunner() (FunctionRunner, error) {
	functions := []Function{
		{
			Exec: "../testdata/clean-metadata",
		},
		{
			Image: "gcr.io/kpt-fn/set-labels:v0.1",
			ConfigMap: map[string]string{
				"env":      "dev",
				"app-name": "my-app",
			},
		},
	}
	runner := NewRunner().WithInput([]byte(exampleService)).
		WithInput([]byte(exampleDeployment)).
		WithFunctions(functions...)

	return runner.Build()
}

func TestRunFnRunner(t *testing.T) {
	functions := []Function{
		{
			Exec: "../testdata/clean-metadata",
		},
		{
			Image: "gcr.io/kpt-fn/set-labels:v0.1",
			ConfigMap: map[string]string{
				"env":      "dev",
				"app-name": "my-app",
			},
		},
	}
	input := unstructured.Unstructured{}
	jsonValue, err := yaml.YAMLToJSON([]byte(exampleService))
	err = input.UnmarshalJSON(jsonValue)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	runner := NewRunnerU().
		WithInputs(input).
		WithInput([]byte(exampleDeployment)).
		WithFunctions(functions...)

	fnRunner, err := runner.Build()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	OutRl, err := fnRunner.Execute()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	expectedLabels := map[string]string{
		"app-name": "my-app",
		"env":      "dev",
		"tier":     "frontend",
		"app":      "guestbook",
	}
	assert.EqualValues(t, expectedLabels, OutRl.Items[1].GetLabels())
}
