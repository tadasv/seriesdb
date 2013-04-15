package series

import (
	"github.com/tadasv/timealign"
)

type DataPointSequence struct {
	granularity   int                    // granularity in minutes
	numDataPoints int                    // number of points in time to keep
	dataPoints    map[int64][]*DataPoint // TODO needs to be optimized
}

func NewDataPointSequence(granularity int, dataPoints int) DataPointSequence {
	seq := DataPointSequence{
		granularity:   granularity,
		numDataPoints: dataPoints,
		dataPoints:    make(map[int64][]*DataPoint),
	}
	return seq
}

func (ds *DataPointSequence) AddDataPoint(dataPoint *DataPoint) {
	key := timealign.AlignToMinutes(dataPoint.timestamp, ds.granularity).Unix()

	points, ok := ds.dataPoints[key]
	if ok == false {
		points = make([]*DataPoint, 1)
		points[0] = dataPoint
		ds.dataPoints[key] = points
		return
	}

	ds.dataPoints[key] = append(points, dataPoint)
}

func (ds *DataPointSequence) Size() int {
	total := 0

	for _, value := range ds.dataPoints {
		total += len(value)
	}

	return total
}
