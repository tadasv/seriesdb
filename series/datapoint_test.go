package series

import (
	"testing"
	"time"
)

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

	point.SetMetric("m1", 10.4)
	value, err = point.GetMetric("m1")
	if err != nil {
		t.Fatal("did not expect error")
	}

	if value != 10.4 {
		t.Error("got invalid value %@ expected 10.4", value)
	}
}
