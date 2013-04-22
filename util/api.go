package util

import (
	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strconv"
)

func ApiResponse(w http.ResponseWriter, statusCode int, statusText string, data interface{}) {
	response, _ := json.Marshal(struct {
		StatusCode int         `json:"statusCode"`
		Data       interface{} `json:"data,omitempty"`
	}{
		statusCode,
		data,
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	w.WriteHeader(statusCode)
	w.Write(response)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, errorCode int, errorMessage string, errors interface{}) {
	response, _ := json.Marshal(struct {
		Error interface{} `json:"error"`
	}{
		struct {
			ErrorCode    int         `json:"errorCode"`
			ErrorMessage string      `json:"message"`
			Errors       interface{} `json:"errors,omitempty"`
		}{
			errorCode,
			errorMessage,
			errors,
		},
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	w.WriteHeader(statusCode)
	w.Write(response)
}

func ExtractData(req *http.Request, w http.ResponseWriter) *simplejson.Json {
	contentType, ok := req.Header["Content-Type"]
	if ok == false || contentType[0] != "application/json" {
		ErrorResponse(w, 400, 400, "Content-Type must be application/json", nil)
		return nil
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ErrorResponse(w, 500, 500, "failed to read body", nil)
		return nil
	}

	json, err := simplejson.NewJson(body)
	if err != nil {
		ErrorResponse(w, 400, 400, "invalid json", nil)
		return nil
	}

	json, ok = json.CheckGet("data")
	if !ok {
		ErrorResponse(w, 400, 400, "invalid json", nil)
		return nil
	}

	return json
}
