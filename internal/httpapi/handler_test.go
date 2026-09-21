package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"job-fit/internal/application"
	"job-fit/internal/evaluator"
	"job-fit/internal/failure"
)

type evaluateFunc func(context.Context, string) (application.Output, error)

func (f evaluateFunc) Evaluate(c context.Context, s string) (application.Output, error) {
	return f(c, s)
}
func fixture(t *testing.T) application.Output {
	t.Helper()
	rows := []evaluator.Assessment{
		{Requirement: evaluator.Requirement{ID: "L001", Text: "Required: API design", Category: "required"}, Match: evaluator.Strong, SupportingEvidence: []evaluator.SupportingEvidence{{Evidence: evaluator.Evidence{ID: "E1", Text: "API design"}, Relation: "direct"}}, Reasoning: "Direct evidence."},
		{Requirement: evaluator.Requirement{ID: "L003", Text: "Preferred: Kubernetes", Category: "preferred"}, Match: evaluator.Unknown, Unknowns: []string{"Not documented"}, Reasoning: "Unknown is not no experience."},
	}
	trace, err := evaluator.ResolveTrace(evaluator.Strong, rows, []evaluator.InfluenceReference{{RequirementID: "L001", Influence: evaluator.Supporting}, {RequirementID: "L003", Influence: evaluator.NonDecisive}})
	if err != nil {
		t.Fatal(err)
	}
	return application.Output{Result: evaluator.EvaluationResult{OverallMatch: evaluator.Strong, OverallReasoning: "Core match.", Requirements: rows, DecisionTrace: trace, ExplanationSource: "application-rendered"}, Metadata: application.Metadata{ProfileVersion: "sha256:" + strings.Repeat("a", 64), EvaluatedAt: "2026-09-21T09:00:00Z", Engine: "jev", EngineModel: "test-model"}}
}
func assertEnvelope(t *testing.T, w *httptest.ResponseRecorder, want int) map[string]any {
	t.Helper()
	var obj map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &obj); err != nil {
		t.Fatal(err, w.Body.String())
	}
	if w.Code != want || obj["code"] != float64(want) {
		t.Fatalf("status=%d want=%d body=%s", w.Code, want, w.Body.String())
	}
	metadata, ok := obj["metadata"].(map[string]any)
	if !ok {
		t.Fatal("missing metadata")
	}
	id, _ := metadata["requestId"].(string)
	if !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(id) {
		t.Fatal("invalid UUID", id)
	}
	if _, err := time.Parse(time.RFC3339Nano, metadata["timestamp"].(string)); err != nil {
		t.Fatal(err)
	}
	if want != 200 {
		if _, ok := obj["data"]; ok {
			t.Fatal("error contains data")
		}
		if obj["status"] != "error" {
			t.Fatal(obj)
		}
	} else if obj["status"] != "success" {
		t.Fatal(obj)
	}
	if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Content-Type") != "application/json" {
		t.Fatal(w.Header())
	}
	return obj
}
func TestSuccessDTOAndCorrelation(t *testing.T) {
	var logs bytes.Buffer
	var requestID string
	calls := 0
	h := New(evaluateFunc(func(ctx context.Context, text string) (application.Output, error) {
		calls++
		requestID = RequestID(ctx)
		if text != "Required: API design" {
			t.Fatalf("metadata inserted: %q", text)
		}
		return fixture(t), nil
	}), time.Second, slog.New(slog.NewJSONHandler(&logs, nil)))
	r := httptest.NewRequest("POST", Endpoint, strings.NewReader(`{"job":{"description":"Required: API design","company":" Example ","title":" Backend ","sourceUrl":"https://example.com/job"}}`))
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Set("X-Request-ID", "caller-secret")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	obj := assertEnvelope(t, w, 200)
	if calls != 1 || requestID != obj["metadata"].(map[string]any)["requestId"] {
		t.Fatal("correlation missing")
	}
	data := obj["data"].(map[string]any)
	if data["job"].(map[string]any)["company"] != "Example" {
		t.Fatal(data)
	}
	rows := data["requirements"].([]any)
	row := rows[0].(map[string]any)
	if row["sourceType"] != "required" || row["missingEvidence"] == nil || row["unknowns"] == nil {
		t.Fatal(row)
	}
	if _, exists := row["category"]; exists {
		t.Fatal("domain category leaked")
	}
	ev := row["supportingEvidence"].([]any)[0].(map[string]any)
	if ev["id"] != "E1" || ev["relation"] != "direct" {
		t.Fatal(ev)
	}
	if data["metadata"].(map[string]any)["profileVersion"] == nil || data["metadata"].(map[string]any)["engineModel"] != "test-model" {
		t.Fatal(data)
	}
	if data["decisionTrace"].(map[string]any)["source"] == nil {
		t.Fatal("missing disclosure")
	}
	if _, ok := data["summary"]; ok {
		t.Fatal("summary leaked")
	}
	if !strings.Contains(logs.String(), requestID) || strings.Contains(logs.String(), "API design") || strings.Contains(logs.String(), "caller-secret") {
		t.Fatal(logs.String())
	}
}
func TestRequestValidation(t *testing.T) {
	valid := `{"job":{"description":"API design"}}`
	cases := []struct {
		name, body, media string
		status            int
		code              string
	}{
		{"empty", "", "application/json", 400, "REQUEST_MALFORMED_JSON"},
		{"syntax", "{", "application/json", 400, "REQUEST_MALFORMED_JSON"},
		{"trailing", valid + ` {}`, "application/json", 400, "REQUEST_MALFORMED_JSON"},
		{"utf8", string([]byte{0xff}), "application/json", 400, "REQUEST_MALFORMED_JSON"},
		{"surrogate", `{"job":{"description":"\ud800"}}`, "application/json", 400, "REQUEST_MALFORMED_JSON"},
		{"low surrogate", `{"job":{"description":"\udc00"}}`, "application/json", 400, "REQUEST_MALFORMED_JSON"},
		{"array", `[]`, "application/json", 400, "REQUEST_INVALID"},
		{"missing job", `{}`, "application/json", 400, "REQUEST_INVALID"},
		{"null job", `{"job":null}`, "application/json", 400, "REQUEST_INVALID"},
		{"wrong type", `{"job":{"description":123}}`, "application/json", 400, "REQUEST_INVALID"},
		{"null optional", `{"job":{"description":"x","title":null}}`, "application/json", 400, "REQUEST_INVALID"},
		{"profile override", `{"job":{"description":"x"},"candidateProfile":{}}`, "application/json", 400, "REQUEST_INVALID"},
		{"unknown", `{"job":{"description":"x","Description":"x"}}`, "application/json", 400, "REQUEST_INVALID"},
		{"duplicate", `{"job":{"description":"x","descrip\u0074ion":"y"}}`, "application/json", 400, "REQUEST_INVALID"},
		{"big body", strings.Repeat(" ", MaxBodyBytes+1), "application/json", 413, "REQUEST_INPUT_TOO_LARGE"},
		{"big description", `{"job":{"description":"` + strings.Repeat("x", 16385) + `"}}`, "application/json", 413, "REQUEST_INPUT_TOO_LARGE"},
		{"many lines", `{"job":{"description":"` + strings.Repeat(`x\n`, 65) + `"}}`, "application/json", 413, "REQUEST_INPUT_TOO_LARGE"},
		{"big metadata", `{"job":{"description":"x","company":"` + strings.Repeat("x", 257) + `"}}`, "application/json", 413, "REQUEST_INPUT_TOO_LARGE"},
		{"media", valid, "text/plain", 415, "REQUEST_UNSUPPORTED_MEDIA_TYPE"},
		{"missing media", valid, "", 415, "REQUEST_UNSUPPORTED_MEDIA_TYPE"},
		{"charset", valid, "application/json; charset=latin1", 415, "REQUEST_UNSUPPORTED_MEDIA_TYPE"},
		{"missing description", `{"job":{}}`, "application/json", 422, "REQUEST_VALIDATION_FAILED"},
		{"blank", `{"job":{"description":" \n "}}`, "application/json", 422, "REQUEST_VALIDATION_FAILED"},
		{"blank optional", `{"job":{"description":"x","company":" "}}`, "application/json", 422, "REQUEST_VALIDATION_FAILED"},
		{"url", `{"job":{"description":"x","sourceUrl":"/relative"}}`, "application/json", 422, "REQUEST_VALIDATION_FAILED"},
		{"credentials url", `{"job":{"description":"x","sourceUrl":"https://secret@example.com"}}`, "application/json", 422, "REQUEST_VALIDATION_FAILED"},
		{"nul", `{"job":{"description":"x\u0000"}}`, "application/json", 422, "REQUEST_VALIDATION_FAILED"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := New(evaluateFunc(func(context.Context, string) (application.Output, error) {
				t.Fatal("invalid request evaluated")
				return application.Output{}, nil
			}), time.Second, nil)
			r := httptest.NewRequest("POST", Endpoint, strings.NewReader(tc.body))
			if tc.media != "" {
				r.Header.Set("Content-Type", tc.media)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			obj := assertEnvelope(t, w, tc.status)
			if obj["error"] != tc.code {
				t.Fatal(obj)
			}
			if tc.status == 422 {
				fields, ok := obj["errors"].([]any)
				if !ok || len(fields) != 1 {
					t.Fatal(obj)
				}
			} else if _, ok := obj["errors"]; ok {
				t.Fatal("unexpected field errors")
			}
		})
	}
}
func TestErrorCategoriesAndRedaction(t *testing.T) {
	secret := "Authorization: Bearer private-secret"
	cases := []struct {
		err    error
		status int
		code   string
	}{
		{errors.New(secret), 500, "INTERNAL_ERROR"},
		{failure.New(failure.Configuration, secret), 500, "INTERNAL_CONFIGURATION_ERROR"},
		{failure.New(failure.Evaluation, secret), 502, "INTERNAL_EVALUATION_FAILED"},
		{failure.New(failure.Unavailable, secret), 503, "INTERNAL_DEPENDENCY_UNAVAILABLE"},
		{failure.New(failure.Timeout, secret), 504, "INTERNAL_EVALUATION_TIMEOUT"},
		{fmt.Errorf("wrapped: %w", context.DeadlineExceeded), 504, "INTERNAL_EVALUATION_TIMEOUT"},
		{failure.New(failure.NoRequirements, secret), 422, "REQUEST_VALIDATION_FAILED"},
		{failure.New(failure.Capacity, secret), 422, "REQUEST_VALIDATION_FAILED"},
	}
	for _, tc := range cases {
		t.Run(tc.code, func(t *testing.T) {
			var logs bytes.Buffer
			h := New(evaluateFunc(func(context.Context, string) (application.Output, error) {
				return application.Output{}, fmt.Errorf("stage: %w", tc.err)
			}), time.Second, slog.New(slog.NewJSONHandler(&logs, nil)))
			r := httptest.NewRequest("POST", Endpoint, strings.NewReader(`{"job":{"description":"API design"}}`))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			obj := assertEnvelope(t, w, tc.status)
			if obj["error"] != tc.code || strings.Contains(w.Body.String()+logs.String(), "private-secret") {
				t.Fatal(obj, logs.String())
			}
		})
	}
}
func TestInvalidTraceFailsWithoutRepair(t *testing.T) {
	for _, kind := range []string{"strong-limiting", "duplicate", "missing", "wrong-match", "unknown-supporting"} {
		t.Run(kind, func(t *testing.T) {
			out := fixture(t)
			switch kind {
			case "strong-limiting":
				out.Result.DecisionTrace.Limiting = out.Result.DecisionTrace.NonDecisive
				out.Result.DecisionTrace.NonDecisive = nil
			case "duplicate":
				out.Result.DecisionTrace.Supporting = append(out.Result.DecisionTrace.Supporting, out.Result.DecisionTrace.Supporting[0])
			case "missing":
				out.Result.DecisionTrace.NonDecisive[0].RequirementID = "L999"
			case "wrong-match":
				out.Result.DecisionTrace.Supporting[0].Match = evaluator.Weak
			case "unknown-supporting":
				out.Result.DecisionTrace.Supporting = append(out.Result.DecisionTrace.Supporting, out.Result.DecisionTrace.NonDecisive...)
				out.Result.DecisionTrace.NonDecisive = nil
			}
			h := New(evaluateFunc(func(context.Context, string) (application.Output, error) { return out, nil }), time.Second, nil)
			r := httptest.NewRequest("POST", Endpoint, strings.NewReader(`{"job":{"description":"x"}}`))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			assertEnvelope(t, w, 502)
		})
	}
}
func TestTimeoutAndCallerCancellation(t *testing.T) {
	for _, cancelCaller := range []bool{false, true} {
		t.Run(fmt.Sprint(cancelCaller), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			observed := make(chan error, 1)
			h := New(evaluateFunc(func(c context.Context, _ string) (application.Output, error) {
				if cancelCaller {
					cancel()
				}
				<-c.Done()
				observed <- c.Err()
				return application.Output{}, c.Err()
			}), 10*time.Millisecond, nil)
			r := httptest.NewRequest("POST", Endpoint, strings.NewReader(`{"job":{"description":"x"}}`)).WithContext(ctx)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if cancelCaller {
				if !errors.Is(<-observed, context.Canceled) || w.Body.Len() != 0 {
					t.Fatal("cancellation ignored")
				}
			} else {
				if !errors.Is(<-observed, context.DeadlineExceeded) {
					t.Fatal("deadline missing")
				}
				assertEnvelope(t, w, 504)
			}
		})
	}
}
func TestRoutingAndOptionalMetadata(t *testing.T) {
	for _, tc := range []struct {
		method, path, encoding string
		status                 int
	}{{"GET", Endpoint, "", 405}, {"POST", "/private-secret?token=secret", "", 404}, {"POST", Endpoint + "?x=1", "", 400}, {"POST", Endpoint, "gzip", 415}, {"POST", Endpoint, "", 200}} {
		var logs bytes.Buffer
		h := New(evaluateFunc(func(context.Context, string) (application.Output, error) { return fixture(t), nil }), time.Second, slog.New(slog.NewJSONHandler(&logs, nil)))
		r := httptest.NewRequest(tc.method, tc.path, strings.NewReader(`{"job":{"description":"API \ud83d\ude00"}}`))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Content-Encoding", tc.encoding)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		obj := assertEnvelope(t, w, tc.status)
		if tc.status == 200 && len(obj["data"].(map[string]any)["job"].(map[string]any)) != 0 {
			t.Fatal("metadata must be empty object")
		}
		if tc.status == 405 && w.Header().Get("Allow") != "POST" {
			t.Fatal("missing Allow")
		}
		if strings.Contains(logs.String(), "secret") {
			t.Fatal("path leaked")
		}
	}
}
