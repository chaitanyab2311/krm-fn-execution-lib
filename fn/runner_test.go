package fn

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
	"testing"
)

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

	temp := unstructured.Unstructured{}
	jsonValue, err := yaml.YAMLToJSON([]byte(exampleService))
	err = temp.UnmarshalJSON(jsonValue)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	runner := NewRunner().
		WithInputs(&temp).
		WithInput([]byte(exampleDeployment)).
		WithFunctions(functions...).
		WhereExecWorkingDir("/usr")

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
	assert.EqualValues(t, expectedLabels, OutRl.Items[1].(*unstructured.Unstructured).GetLabels())
}
