package main

import (
	"github.com/tadasv/seriesdb/series"
	"log"
	"net"
)

type SeriesDB struct {
	options      *SeriesDBOptions
	httpAddress  *net.TCPAddr
	httpListener net.Listener
	exitChannel  chan int
	dataPointSeq series.DataPointSequence
}

type SeriesDBOptions struct {
	granularity int
	dataPoints  int
}

func NewSeriesDB(options *SeriesDBOptions) *SeriesDB {
	dbd := &SeriesDB{
		options: options,
	}

	dbd.dataPointSeq = series.NewDataPointSequence(options.granularity, options.dataPoints)

	return dbd
}

func (db *SeriesDB) Start() {
	httpListener, err := net.Listen("tcp", db.httpAddress.String())
	if err != nil {
		log.Fatalf("FATAL: listen (%s) failed - %s", db.httpAddress, err.Error())
	}

	db.httpListener = httpListener

	go httpServer(db.httpListener)
}

func (db *SeriesDB) Stop() {
	if db.httpListener != nil {
		db.httpListener.Close()
	}

	close(db.exitChannel)
}
