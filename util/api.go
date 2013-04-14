package util

import (
	"encoding/json"
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
