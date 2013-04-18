package series

import (
	"testing"
	"time"
)

func Test_Constructor(t *testing.T) {
	point := NewDataPoint()

	if point.NumMetrics() != 0 {
		t.Fatal("expected 0 metrics")
	}

	if point.NumDimensions() != 0 {
		t.Fatal("expected 0 dimensions")
	}
}

func Test_SetTimestamp(t *testing.T) {
	point := NewDataPoint()
	ts := time.Now()
	point.SetTimestamp(ts)
	if ts != point.timestamp {
		t.Error("timestamps don't match")
	}
}

func Test_Dimensions(t *testing.T) {
	point := NewDataPoint()

	value, err := point.GetDimension("category")
	if err == nil {
		t.Fatal("expected error")
	}

	point.SetDimension("category", "A")
	value, err = point.GetDimension("category")
	if err != nil {
		t.Fatal("did not expect error")
	}

	if *value != "A" {
		t.Error("got invalid value %@ expected A", value)
	}

	if point.NumDimensions() != 1 {
		t.Fatal("expected 1 dimension")
	}
}

func Test_Metrics(t *testing.T) {
	point := NewDataPoint()

	value, err := point.GetMetric("m1")
	if err == nil {
		t.Fatal("expected error")
	}

	point.SetMetric("m1", 1.1)
	value, err = point.GetMetric("m1")
	if err != nil {
		t.Fatal("did not expect error")
	}

	if value != 1.1 {
		t.Error("got invalid value %@ expected 1.1", value)
	}

	if point.NumMetrics() != 1 {
		t.Fatal("expected 1 dimension")
	}

	point.SetMetric("m1", 10.4)
	value, err = point.GetMetric("m1")
	if err != nil {
		t.Fatal("did not expect error")
	}

	if point.NumMetrics() != 1 {
		t.Fatal("expected 1 dimension")
	}

	if value != 10.4 {
		t.Error("got invalid value %@ expected 10.4", value)
	}
}
