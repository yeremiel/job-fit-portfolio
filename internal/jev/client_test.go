package jev

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"job-fit/internal/evaluator"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func reply(body string, status int) *http.Response {
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}
func testClient(t *testing.T, f transportFunc) *Client {
	t.Helper()
	c, err := New("test-secret-not-real")
	if err != nil {
		t.Fatal(err)
	}
	c.http.Transport = f
	return c
}

func TestMissingKey(t *testing.T) {
	for _, key := range []string{"", " \n"} {
		if _, err := New(key); err == nil || !strings.Contains(err.Error(), "TYPESAFE_API_KEY") {
			t.Fatalf("%v", err)
		}
	}
}

func TestRequestAndMapping(t *testing.T) {
	q := map[string]question{"match": choice("match?", "requirement", map[string]string{"Strong": "direct", "Unknown": "insufficient"})}
	c := testClient(t, func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != endpoint || r.Method != "POST" || r.Header.Get("Authorization") != "Bearer test-secret-not-real" || r.Header.Get("Content-Type") != "application/json" {
			t.Fatal("wrong request")
		}
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.Model != model || req.Questions["match"].Type != "choice" {
			t.Fatal("wrong wire shape")
		}
		return reply(`{"model":"jev-1.13.0","answers":{"match":{"type":"choice","choice":"Strong","confidence":0.5,"probabilities":{"Strong":0.6,"Unknown":0.4}}},"usage":{"input_tokens":100}}`, 200), nil
	})
	answers, err := c.choices(context.Background(), map[string]string{"document": "raw text"}, q)
	if err != nil || answers["match"] != "Strong" {
		t.Fatalf("%v %v", answers, err)
	}
	for _, body := range []string{`{`, `null`, `{}`, `{"model":"m","answers":{}}`, `{"model":"m","answers":{"match":{"type":"score","choice":"Strong"}}}`, `{"model":"m","answers":{"match":{"type":"choice","choice":"HIGH"}}}`, `{"model":"m","answers":{"match":{"type":"choice"}}}`, `{"model":"m","answers":{"extra":{"type":"choice","choice":"Strong"}}}`} {
		if _, err := mapChoices([]byte(body), q); err == nil {
			t.Errorf("accepted bad response %s", body)
		}
	}
}

func TestTransportErrorsDoNotLeakSecrets(t *testing.T) {
	q := map[string]question{"q": choice("?", nil, map[string]string{"yes": "yes"})}
	for _, status := range []int{400, 401, 403, 422, 429, 500, 529} {
		c := testClient(t, func(*http.Request) (*http.Response, error) { return reply("test-secret-not-real", status), nil })
		_, err := c.choices(context.Background(), nil, q)
		if err == nil || !strings.Contains(err.Error(), "HTTP") || strings.Contains(err.Error(), "test-secret") {
			t.Fatalf("unsafe/missing error: %v", err)
		}
	}
	c := testClient(t, func(*http.Request) (*http.Response, error) {
		return nil, errors.New("Authorization: test-secret-not-real")
	})
	if _, err := c.choices(context.Background(), nil, q); err == nil || strings.Contains(err.Error(), "test-secret") {
		t.Fatalf("%v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c = testClient(t, func(r *http.Request) (*http.Response, error) { return nil, r.Context().Err() })
	if _, err := c.choices(ctx, nil, q); !errors.Is(err, context.Canceled) {
		t.Fatalf("%v", err)
	}
}

func TestRedirectAndResponseSize(t *testing.T) {
	q := map[string]question{"q": choice("?", nil, map[string]string{"yes": "yes"})}
	calls := 0
	c := testClient(t, func(*http.Request) (*http.Response, error) {
		calls++
		r := reply("", 302)
		r.Header.Set("Location", "https://elsewhere.invalid")
		return r, nil
	})
	if _, err := c.choices(context.Background(), nil, q); err == nil || calls != 1 {
		t.Fatalf("redirect followed: %d %v", calls, err)
	}
	c = testClient(t, func(*http.Request) (*http.Response, error) {
		return reply(strings.Repeat("x", 2*1024*1024+1), 200), nil
	})
	if _, err := c.choices(context.Background(), nil, q); err == nil || !strings.Contains(err.Error(), "size") {
		t.Fatalf("%v", err)
	}
}

func TestTimeout(t *testing.T) {
	c := testClient(t, func(r *http.Request) (*http.Response, error) { <-r.Context().Done(); return nil, r.Context().Err() })
	c.http.Timeout = 5 * time.Millisecond
	_, err := c.choices(context.Background(), nil, map[string]question{"q": choice("?", nil, map[string]string{"yes": "yes"})})
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("%v", err)
	}
}

func TestRationaleMapping(t *testing.T) {
	for key, rationale := range rationales {
		row := evaluator.Assessment{Requirement: evaluator.Requirement{Text: "production systems using version 8 for one year"}, SupportingEvidence: []evaluator.SupportingEvidence{{Evidence: evaluator.Evidence{ID: "E1", Text: "exact source text"}, Relation: "direct"}}, MissingEvidence: []string{}, Unknowns: []string{}}
		if err := applyRationale(&row, key); err != nil {
			t.Fatal(err)
		}
		if row.Match != rationale.Level || row.Reasoning == "" || row.SupportingEvidence[0].Text != "exact source text" {
			t.Fatal("incorrect domain mapping")
		}
		if row.Match != evaluator.Strong && (len(row.MissingEvidence) == 0 || len(row.Unknowns) == 0) {
			t.Fatal("missing gap information")
		}
	}
	row := evaluator.Assessment{}
	if applyRationale(&row, "invalid") == nil || applyRationale(&row, "direct") == nil {
		t.Fatal("accepted unsupported assessment")
	}
	row.SupportingEvidence = []evaluator.SupportingEvidence{{Relation: "transferable"}}
	if applyRationale(&row, "direct") == nil {
		t.Fatal("Strong without direct evidence")
	}
}

func TestAssessmentOptionsRespectEvidence(t *testing.T) {
	limited := assessmentOptions(evaluator.Requirement{Text: "API design"}, []evaluator.SupportingEvidence{{Relation: "limited"}})
	if _, ok := limited["direct"]; ok {
		t.Fatal("Strong offered without direct evidence")
	}
	if _, ok := limited["insufficient"]; !ok {
		t.Fatal("Unknown must remain possible")
	}
	direct := assessmentOptions(evaluator.Requirement{Text: "API design"}, []evaluator.SupportingEvidence{{Relation: "direct"}})
	if _, ok := direct["direct"]; !ok {
		t.Fatal("Strong missing for direct evidence")
	}
}
