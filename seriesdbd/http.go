package main

import (
	"github.com/tadasv/seriesdb/series"
	"github.com/tadasv/seriesdb/util"
	"log"
	"net"
	"net/http"
	"time"
)

func httpServer(listener net.Listener) {
	log.Printf("HTTP: listening on %s", listener.Addr().String())

	handler := http.NewServeMux()
	handler.HandleFunc("/alive", aliveHandler)
	handler.HandleFunc("/datapoint", dataPointHandler)
	handler.HandleFunc("/stats", statsHandler)

	server := &http.Server{
		Handler: handler,
	}
	err := server.Serve(listener)
	if err != nil {
		log.Printf("ERROR: http.Serve() - %s", err.Error())
	}

	log.Printf("HTTP: closing %s", listener.Addr().String())
}

func aliveHandler(w http.ResponseWriter, req *http.Request) {
	util.ApiResponse(w, 200, "OK", nil)
}

func statsHandler(w http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})

	dataPointStats := make(map[string]interface{})
	dataPointStats["count"] = seriesdb.dataPointSeq.Size()

	data["datapoints"] = dataPointStats
	util.ApiResponse(w, 200, "OK", data)
}

func dataPointHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		util.ErrorResponse(w, 405, 405, "method not allowed", nil)
		return
	}

	json := util.ExtractData(req, w)
	if json == nil {
		return
	}

	var datapoints []interface{}
	var ok bool
	var err error

	json, ok = json.CheckGet("datapoints")
	if !ok {
		goto error_invalid_json
	}

	datapoints, err = json.Array()
	if err != nil {
		goto error_invalid_json
	}

	for idx := range datapoints {
		dataPointJs := json.GetIndex(idx)
		timestamp, err := dataPointJs.Get("ts").Int64()
		if err != nil {
			util.ErrorResponse(w, 400, 400, "missing timestamp", nil)
			return
		}

		dp := series.NewDataPoint()
		dp.SetTimestamp(time.Unix(timestamp, 0))

		metricsJs, ok := dataPointJs.CheckGet("m")
		if !ok {
			util.ErrorResponse(w, 400, 400, "missing metrics", nil)
			return
		}

		metrics, err := metricsJs.Map()
		if err != nil {
			util.ErrorResponse(w, 400, 400, "invalid metrics", nil)
			return
		}

		for key, _ := range metrics {
			value, err := metricsJs.Get(key).Float64()
			if err != nil {
				util.ErrorResponse(w, 400, 400, "invalid metrics", nil)
				return
			}
			dp.SetMetric(key, value)
		}

		dimensionsJs, ok := dataPointJs.CheckGet("d")
		if !ok {
			util.ErrorResponse(w, 400, 400, "missing dimensions", nil)
			return
		}

		dimensions, err := dimensionsJs.Map()
		if err != nil {
			util.ErrorResponse(w, 400, 400, "invalid dimensions", nil)
			return
		}

		for key, _ := range dimensions {
			value, err := dimensionsJs.Get(key).String()
			if err != nil {
				util.ErrorResponse(w, 400, 400, "invalid metrics", nil)
				return
			}
			dp.SetDimension(key, value)
		}

		if dp.NumMetrics() == 0 || dp.NumDimensions() == 0 {
			util.ErrorResponse(w, 400, 400, "missing dimensions or metrics", nil)
			return
		}

		seriesdb.dataPointSeq.AddDataPoint(&dp)
	}

	util.ApiResponse(w, 200, "OK", nil)
	return

error_invalid_json:
	util.ErrorResponse(w, 400, 400, "invalid json", nil)
	return
}
