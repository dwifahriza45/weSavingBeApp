package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int               `json:"code"`
	Status  string            `json:"status"`
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Error   bool              `json:"error"`
	Valid   map[string]string `json:"valid,omitempty"`
}

func JSON(w http.ResponseWriter, code int, status string, message string, error bool, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Status:  status,
		Message: message,
		Error:   error,
		Data:    data,
	})
}

func JSONError(w http.ResponseWriter, code int, status string, msg string, error bool) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Status:  status,
		Message: msg,
		Error:   error,
	})
}

func JSONErrorWithData(
	w http.ResponseWriter,
	code int,
	status string,
	msg string,
	error bool,
	valid map[string]string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Status:  status,
		Message: msg,
		Error:   error,
		Valid:   valid,
	})
}
