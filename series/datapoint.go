package series

import (
	"fmt"
	"time"
)

type DataPointError struct {
	when time.Time
	what string
}

func (e DataPointError) Error() string {
	return fmt.Sprintf("%v: %v", e.when, e.what)
}

type DataPoint struct {
	dimensions map[string]string
	metrics    map[string]float64
	timestamp  time.Time
}

func NewDataPoint() DataPoint {
	dp := DataPoint{
		dimensions: make(map[string]string),
		metrics:    make(map[string]float64),
		timestamp:  time.Now(),
	}
	return dp
}

func (point *DataPoint) SetTimestamp(t time.Time) {
	point.timestamp = t
}

func (point *DataPoint) SetDimension(name string, value string) {
	point.dimensions[name] = value
}

func (point *DataPoint) GetDimension(name string) (*string, error) {
	value, ok := point.dimensions[name]
	if ok == false {
		return nil, DataPointError{time.Now(), "dimension not found"}
	}

	return &value, nil
}

func (point *DataPoint) SetMetric(name string, value float64) {
	point.metrics[name] = value
}

func (point *DataPoint) GetMetric(name string) (float64, error) {
	value, ok := point.metrics[name]
	if ok == false {
		return 0.0, DataPointError{time.Now(), "metric not found"}
	}

	return value, nil
}
