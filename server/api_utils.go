package main

import (
	"encoding/json"
	"net/http"
)

// JSONResponse is a helper method to write out a json response with status code
func JSONResponse(w http.ResponseWriter, v interface{}, statusCode int) {
	if v == nil {
		http.Error(w, "Response payload nil", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)

}

// JSONResponse is a helper method to easily write out a json response
// with error message and status code
func JSONErrorResponse(w http.ResponseWriter, err error, statusCode int) {

	data := map[string]interface{}{
		"error":       true,
		"message":     err.Error(),
		"status_code": statusCode,
	}
	JSONResponse(w, data, statusCode)

}
