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
		case "error":
			slog.Error("HTTP API error", "status", resp.StatusCode, "path", r.URL.Path, "body", resp.Body)
		case "warn":
			slog.Warn("Client error", "status", resp.StatusCode, "path", r.URL.Path, "body", resp.Body)
		}

	}
}

func writeResponse(w http.ResponseWriter, resp *HTTPResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	if resp.Body == nil {
		return
	}

	var responseBody map[string]any

	if resp.StatusCode >= 400 {
		responseBody = map[string]any{
			"status": http.StatusText(resp.StatusCode),
			"msg":    resp.Body,
		}
	} else {
		responseBody = map[string]any{
			"status": http.StatusText(resp.StatusCode),
			"data":   resp.Body,
		}
	}

	_ = json.NewEncoder(w).Encode(responseBody)
}
