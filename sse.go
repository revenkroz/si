package si

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SSEWriter writes Server-Sent Events to the client.
type SSEWriter struct {
	w http.ResponseWriter
	f http.Flusher
}

// Event sends a named event with data.
func (s *SSEWriter) Event(event string, data string) error {
	if event != "" {
		if _, err := fmt.Fprintf(s.w, "event: %s\n", event); err != nil {
			return err
		}
	}
	for _, line := range strings.Split(data, "\n") {
		if _, err := fmt.Fprintf(s.w, "data: %s\n", line); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(s.w, "\n"); err != nil {
		return err
	}
	s.f.Flush()
	return nil
}

// Data sends an unnamed event with data.
func (s *SSEWriter) Data(data string) error {
	return s.Event("", data)
}

// JSON sends a named event with JSON-encoded data.
func (s *SSEWriter) JSON(event string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Event(event, string(b))
}

// ID sends an id field. The client will use this as the Last-Event-ID
// on reconnect.
func (s *SSEWriter) ID(id string) error {
	_, err := fmt.Fprintf(s.w, "id: %s\n", id)
	if err != nil {
		return err
	}
	s.f.Flush()
	return nil
}

// Retry tells the client to wait the given number of milliseconds
// before reconnecting.
func (s *SSEWriter) Retry(ms int) error {
	_, err := fmt.Fprintf(s.w, "retry: %d\n\n", ms)
	if err != nil {
		return err
	}
	s.f.Flush()
	return nil
}

// Comment sends an SSE comment (line starting with ":").
// Useful as a keep-alive ping.
func (s *SSEWriter) Comment(text string) error {
	_, err := fmt.Fprintf(s.w, ": %s\n\n", text)
	if err != nil {
		return err
	}
	s.f.Flush()
	return nil
}

// SSE sets up an SSE connection and calls fn with an SSEWriter.
// The caller should use ctx.Request.Context().Done() to detect
// client disconnection inside fn.
func (ctx *Context) SSE(fn func(w *SSEWriter)) {
	flusher, ok := ctx.Response.(http.Flusher)
	if !ok {
		http.Error(ctx.Response, "streaming not supported", http.StatusInternalServerError)
		return
	}

	ctx.Response.Header().Set("Content-Type", "text/event-stream")
	ctx.Response.Header().Set("Cache-Control", "no-cache")
	ctx.Response.Header().Set("Connection", "keep-alive")
	ctx.Response.WriteHeader(http.StatusOK)
	flusher.Flush()

	fn(&SSEWriter{
		w: ctx.Response,
		f: flusher,
	})
}
