package main

import (
	simplejson "github.com/bitly/go-simplejson"
	"github.com/tadasv/seriesdb/series"
	"github.com/tadasv/seriesdb/util"
	"io/ioutil"
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

func dataPointHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		util.ErrorResponse(w, 405, 405, "method not allowed", nil)
		return
	}

	contentType, ok := req.Header["Content-Type"]
	if ok == false || contentType[0] != "application/json" {
		util.ErrorResponse(w, 400, 400, "Content-Type must be application/json", nil)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.ErrorResponse(w, 500, 500, "failed to read body", nil)
		return
	}

	var datapoints []interface{}

	json, err := simplejson.NewJson(body)
	if err != nil {
		goto error_invalid_json
	}

	json, ok = json.CheckGet("data")
	if !ok {
		goto error_invalid_json
	}

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
