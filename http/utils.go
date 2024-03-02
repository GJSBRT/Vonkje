package http

import (
	"reflect"
	"net/http"
	"encoding/json"
)

// SendErrorResponse sends an error response to the client
func (httpServer *HTTP) SendErrorResponse(w http.ResponseWriter, err string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: err,
	}); err != nil {
		httpServer.log.Warn("Error while encoding error response")
	}
}


// WriteJSONResponse sends a JSON response to the client
func (httpServer *HTTP) WriteJSONResponse(w http.ResponseWriter, req *http.Request, statusCode int, response interface{}) {
	shouldRespond := true
	for _, method := range []string{http.MethodHead, http.MethodDelete} {
		if req.Method == method {
			shouldRespond = false
		}
	}

	if shouldRespond {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(statusCode)

	if !shouldRespond {
		return
	}

	val := reflect.ValueOf(response)

	if val.Kind() == reflect.Slice && val.IsNil() {
		if err := json.NewEncoder(w).Encode([]string{}); err != nil {
			httpServer.log.WithError(err).Error("Error while encoding empty array response")
		}

		return
	}

	if response == nil {
		if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
			httpServer.log.WithError(err).Error("Error while encoding empty struct response")
		}

		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		httpServer.log.WithError(err).Error("Error while encoding json response")
	}
}
