package main

import (
	"flag"
	"fmt"
	"github.com/tadasv/seriesdb/util"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	showVersion = flag.Bool("version", false, "print version information")
	httpAddress = flag.String("http-address", "0.0.0.0:19000", "address of the seriesdb HTTP interface")
	granularity = flag.Int("granularity", 1, "granularity of the data stored in memory in minutes")
	dataPoints  = flag.Int("data-points", 60*24, "number of datapoints to keep")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("seriesdb v%s\n", util.BINARY_VERSION)
		return
	}

	httpAddr, err := net.ResolveTCPAddr("tcp", *httpAddress)
	if err != nil {
		log.Fatal(err)
	}

	signalChannel := make(chan os.Signal, 1)
	exitChannel := make(chan int)
	go func() {
		<-signalChannel
		log.Printf("Received SIGINT")
		exitChannel <- 1
	}()
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("seriesdb v%s", util.BINARY_VERSION)

	dbd := NewSeriesDB(&SeriesDBOptions{
		granularity: *granularity,
		dataPoints:  *dataPoints,
	})
	dbd.httpAddress = httpAddr
	dbd.exitChannel = exitChannel

	dbd.Start()
	// Wait till we get data in the exit channel
	<-exitChannel
	dbd.Stop()
}
