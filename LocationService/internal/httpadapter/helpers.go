package httpadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/AntYats/go_project/internal/service"
	"net/http"
	"strconv"
)

func writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}

func writeJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't marshal data: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	writeResponse(w, status, string(response))
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		writeJSONResponse(w, http.StatusNotFound, Error{Message: err.Error()})
	default:
		writeJSONResponse(w, http.StatusInternalServerError, Error{Message: err.Error()})
	}
}

func convertToFloat(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0.0
	}
	return float32(f)
}
