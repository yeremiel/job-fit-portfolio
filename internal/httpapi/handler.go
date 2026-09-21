// Package httpapi implements the v1 transport contract without evaluation policy.
package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"job-fit/internal/application"
)

const Endpoint = "/v1/job-fit/evaluate"

type Evaluator interface {
	Evaluate(context.Context, string) (application.Output, error)
}
type Handler struct {
	service Evaluator
	timeout time.Duration
	logger  *slog.Logger
}

func New(service Evaluator, timeout time.Duration, logger *slog.Logger) *Handler {
	return &Handler{service, timeout, logger}
}

type requestIDKey struct{}

func RequestID(ctx context.Context) string { s, _ := ctx.Value(requestIDKey{}).(string); return s }
func newRequestID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type responseMetadata struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
}
type envelope struct {
	Status   string           `json:"status"`
	Code     int              `json:"code"`
	Message  string           `json:"message"`
	Data     *resultDTO       `json:"data,omitempty"`
	Error    string           `json:"error,omitempty"`
	Errors   []fieldError     `json:"errors,omitempty"`
	Metadata responseMetadata `json:"metadata"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id := newRequestID()
	r = r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id))
	status := 0
	errorCode := ""
	outcome := "completed"
	// Never log arbitrary paths/queries/headers: they may carry caller secrets.
	route := "unmatched"
	if r.URL.Path == Endpoint {
		route = Endpoint
	}
	defer func() {
		if h.logger != nil {
			h.logger.Info("http request", "requestId", id, "method", r.Method, "path", route, "code", status, "error", errorCode, "outcome", outcome, "elapsed", time.Since(start).String())
		}
	}()
	write := func(e envelope) {
		e.Metadata = responseMetadata{time.Now().UTC().Format(time.RFC3339Nano), id}
		b, err := json.Marshal(e)
		if err != nil {
			e = envelope{Status: "error", Code: 500, Error: codeInternal, Message: "An internal error occurred.", Metadata: e.Metadata}
			b, _ = json.Marshal(e)
		}
		status = e.Code
		errorCode = e.Error
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(status)
		if _, err = w.Write(append(b, '\n')); err != nil {
			outcome = "write_failed"
		}
	}
	fail := func(e *apiError) {
		write(envelope{Status: "error", Code: e.status, Error: e.code, Message: e.message, Errors: e.fields})
	}
	// Catch unexpected failures without net/http logging a panic value containing inputs.
	defer func() {
		if recover() != nil {
			if status == 0 {
				fail(problem(500))
			}
			outcome = "internal_failure"
		}
	}()
	if r.URL.Path != Endpoint {
		fail(problem(404))
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		fail(problem(405))
		return
	}
	in, err := parseRequest(w, r)
	if err != nil {
		fail(err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	out, evalErr := h.service.Evaluate(ctx, in.description)
	if r.Context().Err() != nil {
		outcome = "caller_canceled"
		return
	}
	if ctx.Err() != nil {
		fail(problem(504))
		return
	}
	if evalErr != nil {
		fail(mapError(evalErr))
		return
	}
	dto, mapErr := present(out, in.job)
	if mapErr != nil {
		fail(problem(502))
		return
	}
	write(envelope{Status: "success", Code: 200, Message: "ok", Data: &dto})
}
