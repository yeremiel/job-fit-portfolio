package jev

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"job-fit/internal/application"
	"job-fit/internal/failure"
	"job-fit/internal/httpapi"
)

func TestTypedTransportClassification(t *testing.T) {
	q := map[string]question{"q": choice("?", nil, map[string]string{"yes": "yes"})}
	for _, tc := range []struct {
		status int
		kind   failure.Kind
	}{{400, failure.Evaluation}, {401, failure.Configuration}, {403, failure.Configuration}, {422, failure.Evaluation}, {429, failure.Unavailable}, {500, failure.Evaluation}, {529, failure.Unavailable}} {
		c := testClient(t, func(*http.Request) (*http.Response, error) { return reply("private-secret", tc.status), nil })
		_, err := c.choices(context.Background(), nil, q)
		if failure.Classify(err) != tc.kind || strings.Contains(err.Error(), "private-secret") {
			t.Fatal(tc, err)
		}
	}
	c := testClient(t, func(r *http.Request) (*http.Response, error) { <-r.Context().Done(); return nil, r.Context().Err() })
	c.http.Timeout = 5 * time.Millisecond
	_, err := c.choices(context.Background(), nil, q)
	if failure.Classify(err) != failure.Timeout {
		t.Fatal(err)
	}
	c = testClient(t, func(*http.Request) (*http.Response, error) { return nil, errors.New("private-secret") })
	_, err = c.choices(context.Background(), nil, q)
	if failure.Classify(err) != failure.Evaluation {
		t.Fatal(err)
	}
}

func TestHTTPThroughApplicationCancelsJev(t *testing.T) {
	for _, disconnect := range []bool{false, true} {
		entered := make(chan struct{})
		observed := make(chan error, 1)
		c := testClient(t, func(r *http.Request) (*http.Response, error) {
			if httpapi.RequestID(r.Context()) == "" {
				t.Error("request ID not propagated")
			}
			close(entered)
			<-r.Context().Done()
			observed <- r.Context().Err()
			return nil, r.Context().Err()
		})
		name, model := c.Identity()
		service, err := application.New("../../data/profile.json", c, name, model)
		if err != nil {
			t.Fatal(err)
		}
		deadline := 20 * time.Millisecond
		if disconnect {
			deadline = time.Second
		}
		server := httptest.NewServer(httpapi.New(service, deadline, nil))
		ctx, cancel := context.WithCancel(context.Background())
		req, err := http.NewRequestWithContext(ctx, "POST", server.URL+httpapi.Endpoint, strings.NewReader(`{"job":{"description":"Required: API design"}}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		if disconnect {
			go func() { <-entered; cancel() }()
		}
		response, err := server.Client().Do(req)
		if disconnect {
			if err == nil {
				t.Error("expected canceled caller")
				response.Body.Close()
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			b, _ := io.ReadAll(response.Body)
			response.Body.Close()
			if response.StatusCode != 504 || !strings.Contains(string(b), "INTERNAL_EVALUATION_TIMEOUT") {
				t.Fatal(response.StatusCode, string(b))
			}
		}
		select {
		case cause := <-observed:
			want := context.DeadlineExceeded
			if disconnect {
				want = context.Canceled
			}
			if !errors.Is(cause, want) {
				t.Errorf("got %v want %v", cause, want)
			}
		case <-time.After(2 * time.Second):
			t.Error("Jev request was not canceled")
		}
		cancel()
		server.Close()
	}
}
