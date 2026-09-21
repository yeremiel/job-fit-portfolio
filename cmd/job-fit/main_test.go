package main

import (
	"bytes"
	"context"
	"encoding/json"
	"job-fit/internal/evaluator"
	"strings"
	"testing"
)

func TestCLIInputAndMissingKey(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want string
		code int
	}{
		{"help", []string{"--help"}, "profile", 0},
		{"flags", nil, "--profile", 1},
		{"missing profile", []string{"--profile", "/missing/profile", "--job", "job.txt"}, "profile file", 1},
		{"missing job", []string{"--profile", "../../data/profile.json", "--job", "/missing/job"}, "job file", 1},
		{"missing key", []string{"--profile", "../../data/profile.json", "--job", "../../samples/proptech-plus.txt"}, "TYPESAFE_API_KEY is not set", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out, errout bytes.Buffer
			code := run(context.Background(), tc.args, &out, &errout, func(key string) string {
				if key != "TYPESAFE_API_KEY" {
					t.Fatalf("unexpected env %s", key)
				}
				return ""
			})
			if code != tc.code || !strings.Contains(errout.String(), tc.want) || out.Len() != 0 {
				t.Fatalf("code=%d out=%q err=%q", code, out.String(), errout.String())
			}
		})
	}
}

type offlineEngine struct{}

func (offlineEngine) Identity() (string, string) { return "jev", "offline" }
func (offlineEngine) Extract(_ context.Context, j evaluator.Job) ([]evaluator.Requirement, error) {
	return []evaluator.Requirement{{ID: "L001", Text: j.Lines[0], Category: "context"}}, nil
}
func (offlineEngine) Assess(_ context.Context, _ evaluator.Profile, _ evaluator.Job, r []evaluator.Requirement) ([]evaluator.Assessment, error) {
	return []evaluator.Assessment{{Requirement: r[0], Match: evaluator.Unknown, SupportingEvidence: []evaluator.SupportingEvidence{}, MissingEvidence: []string{}, Unknowns: []string{"not documented"}}}, nil
}
func (offlineEngine) Overall(context.Context, evaluator.Profile, evaluator.Job, []evaluator.Assessment) (evaluator.MatchLevel, string, error) {
	return evaluator.Unknown, "unknown", nil
}
func (offlineEngine) Trace(context.Context, evaluator.Profile, evaluator.Job, []evaluator.Assessment, evaluator.MatchLevel) ([]evaluator.InfluenceReference, error) {
	return []evaluator.InfluenceReference{{RequirementID: "L001", Influence: evaluator.Limiting}}, nil
}
func TestCLISuccessKeepsDomainJSON(t *testing.T) {
	var out, errout bytes.Buffer
	code := runWithEngine(context.Background(), []string{"--profile", "../../data/profile.json", "--job", "../../samples/proptech-plus.txt"}, &out, &errout, func(string) string { return "" }, func(string) (namedEngine, error) { return offlineEngine{}, nil })
	if code != 0 || errout.Len() != 0 {
		t.Fatal(code, errout.String())
	}
	var data map[string]any
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"overallMatch", "overallReasoning", "requirements", "summary", "explanationSource", "decisionTrace"} {
		if _, ok := data[key]; !ok {
			t.Fatal("missing CLI field", key)
		}
	}
	for _, key := range []string{"status", "code", "data", "metadata"} {
		if _, ok := data[key]; ok {
			t.Fatal("HTTP field in CLI", key)
		}
	}
}
