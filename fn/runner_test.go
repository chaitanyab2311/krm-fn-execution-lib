package fn

import (
	"github.com/stretchr/testify/assert"
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
	fnRunner, err := getFnRunner()
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	OutRl, err := fnRunner.Run()
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
