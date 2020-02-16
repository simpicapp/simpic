package http

import (
	"encoding/json"
	"net/http"
)

// writeJSON marshals the given payload and writes it to the ResponseWriter.
func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
