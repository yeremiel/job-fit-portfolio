package jev

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"job-fit/internal/failure"
)

const endpoint = "https://api.typesafe.ai/v1/systemone"
const model = "jev-1.13.0"

type Client struct {
	key  string
	http *http.Client
}

// New takes a value read from TYPESAFE_API_KEY by the CLI. It never loads files.
func New(apiKey string) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, failure.New(failure.Configuration, "TYPESAFE_API_KEY is not set")
	}
	return &Client{key: strings.TrimSpace(apiKey), http: &http.Client{
		Timeout:       60 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}}, nil
}

type question struct {
	Type         string            `json:"type"`
	Instructions any               `json:"instructions"`
	Criteria     map[string]string `json:"criteria"`
}

type request struct {
	Model     string              `json:"model"`
	State     any                 `json:"state"`
	Questions map[string]question `json:"questions"`
}

type answer struct {
	Type   string `json:"type"`
	Choice string `json:"choice"`
}

type response struct {
	Model   string            `json:"model"`
	Answers map[string]answer `json:"answers"`
}

// choices discards provider probability/confidence fields: they are not match scores.
func (c *Client) choices(ctx context.Context, state any, questions map[string]question) (map[string]string, error) {
	body, err := json.Marshal(request{model, state, questions})
	if err != nil {
		return nil, failure.New(failure.Internal, "cannot encode Jev request")
	}
	if len(body) > 192*1024 {
		return nil, failure.New(failure.Capacity, "Jev request exceeds local size limit; shorten inputs")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, failure.New(failure.Internal, "cannot create Jev request")
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("Jev request interrupted: %w", ctx.Err())
		}
		return nil, transportError(ctx, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// Never print response bodies or transport errors: either can echo credentials.
		kind := failure.Evaluation
		hint := "check service availability"
		switch resp.StatusCode {
		case 401, 403:
			kind = failure.Configuration
			hint = "check TYPESAFE_API_KEY and API access"
		case 400, 422:
			hint = "request rejected; check API compatibility and input size"
		case 429, 529:
			kind = failure.Unavailable
			hint = "service busy or rate limited; retry later"
		}
		return nil, failure.New(kind, fmt.Sprintf("Jev HTTP %d: %s", resp.StatusCode, hint))
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024+1))
	if err != nil {
		return nil, transportError(ctx, err)
	}
	if len(b) > 2*1024*1024 {
		return nil, errors.New("Jev response exceeds local size limit")
	}
	return mapChoices(b, questions)
}

func mapChoices(b []byte, questions map[string]question) (map[string]string, error) {
	var r response
	if json.Unmarshal(b, &r) != nil || strings.TrimSpace(r.Model) == "" || len(r.Answers) != len(questions) {
		return nil, errors.New("invalid Jev response: missing model or answers")
	}
	result := make(map[string]string, len(questions))
	for id, q := range questions {
		a, ok := r.Answers[id]
		if !ok || a.Type != "choice" {
			return nil, errors.New("invalid Jev response: missing or unexpected answer type")
		}
		if _, ok := q.Criteria[a.Choice]; !ok {
			return nil, errors.New("invalid Jev response: choice is not an allowed option")
		}
		result[id] = a.Choice
	}
	return result, nil
}

// Identity describes the configured model used in all requests, not a provider claim.
func (c *Client) Identity() (string, string) { return "jev", model }

func transportError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return fmt.Errorf("Jev request interrupted: %w", ctx.Err())
	}
	var n net.Error
	if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &n) && n.Timeout()) {
		return failure.New(failure.Timeout, "Jev request timeout")
	}
	return failure.New(failure.Evaluation, "Jev request failed: network error")
}
