package hexttp

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) *HTTPResponse

func Make(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := h(w, r)
		if resp == nil || resp.StatusCode == http.StatusNoContent {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		writeResponse(w, resp)
		logResponse(r, resp)
	}
}

type SuccessBodyResponse struct {
	Status string         `json:"status"`
	Data   any            `json:"data,omitempty"`
	Meta   map[string]any `json:"meta,omitempty"`
}

type ErrorBodyResponse struct {
	Status string `json:"status"`
	Msg    any    `json:"msg"`
}

func logResponse(r *http.Request, resp *HTTPResponse) {
	requestID, _ := r.Context().Value(requestIDKey).(string)
	start, _ := r.Context().Value(requestStart).(time.Time)
	duration := time.Since(start)
	switch resp.LogLevel {
	case LogError:
		slog.Error("HTTP API error", "RequestID", requestID, "duration", duration, "status", resp.StatusCode, "path", r.URL.Path, "body", resp.Body)
	case LogWarn:
		slog.Warn("Client error", "RequestID", requestID, "duration", duration, "status", resp.StatusCode, "path", r.URL.Path, "body", resp.Body)
	}
}

func writeResponse(w http.ResponseWriter, resp *HTTPResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	if resp.Body == nil {
		return
	}

	var output any
	if resp.StatusCode >= 400 {
		output = ErrorBodyResponse{
			Status: http.StatusText(resp.StatusCode),
			Msg:    resp.Body,
		}
	} else {
		output = SuccessBodyResponse{
			Status: http.StatusText(resp.StatusCode),
			Data:   resp.Body,
		}
	}
	_ = json.NewEncoder(w).Encode(output)
}
