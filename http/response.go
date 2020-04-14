package http

import (
	"encoding/json"
	"net/http"
)

type errorResp struct {
	Description string `json:"description"`
	Code        int    `json:"code"`
}

func writeJSON(w http.ResponseWriter, r interface{}, status int) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	if status != 0 {
		w.WriteHeader(status)
	}
	body, _ := json.MarshalIndent(r, "", "  ")
	w.Write(body)
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	if status != 0 {
		w.WriteHeader(status)
	}
	r := errorResp{
		Description: message,
		Code:        status,
	}
	body, _ := json.MarshalIndent(r, "", "  ")
	w.Write(body)
}
