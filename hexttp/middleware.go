package hexttp

import (
	"encoding/json"
	"log/slog"
	"net/http"
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

		switch resp.LogLevel {
		case LogError:
			slog.Error("HTTP API error", "status", resp.StatusCode, "path", r.URL.Path, "body", resp.Body)
		case LogWarn:
			slog.Warn("Client error", "status", resp.StatusCode, "path", r.URL.Path, "body", resp.Body)
		}

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
