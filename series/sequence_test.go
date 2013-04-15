package series

import (
	"testing"
)

func Test_AddDataPoint(t *testing.T) {
	seq := NewDataPointSequence(1, 10)

	for i := 0 ; i < 20; i++ {
		point := NewDataPoint()
		seq.AddDataPoint(&point)
	}

	size := seq.Size()
	if 20 != size {
		t.Fatal("expected 20 data points in sequence got", size)
	}
}
