package metrics

import (
	"testing"
)

func TestGet(t *testing.T) {
	AddMetric("tests", "get_test", "A test testing the get function", []string{
		"inverter",
		"string",
	})

	metric := GetMetric("tests", "get_test")
	if metric == nil {
		t.Fatalf("Metric not found")
	}

	if metric.Name != "get_test" {
		t.Fatalf("Incorrect name")
	}
}

func TestSet(t *testing.T) {
	AddMetric("tests", "set_test", "A test testing the set function", []string{
		"inverter",
		"string",
	})

	SetMetricValue("tests", "set_test", map[string]string{
		"inverter": "1",
		"string": "1",
	}, 1.0)

	metric := GetMetric("tests", "set_test")
	if metric == nil {
		t.Fatalf("Metric not found")
	}

	if metric.Values[0].Value != 1.0 {
		t.Fatalf("Incorrect value")
	}
}
