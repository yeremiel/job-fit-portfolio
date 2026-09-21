package jev

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"job-fit/internal/evaluator"
)

// All five HTTP stages use an in-memory transport. No external API or key required.
func TestVerticalSliceOffline(t *testing.T) {
	p := evaluator.Profile{Evidence: []evaluator.Evidence{{ID: "E1", Text: "API design"}, {ID: "E2", Text: "Basic cloud usage"}}}
	j := evaluator.Job{Text: "Required: API design\nPreferred: cloud production\nSource: example", Lines: []string{"Required: API design", "Preferred: cloud production", "Source: example"}}
	stages := 0
	c := testClient(t, func(r *http.Request) (*http.Response, error) {
		stages++
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		selected := map[string]string{}
		switch stages {
		case 1:
			selected = map[string]string{"L001": "required", "L002": "preferred", "L003": "ignore"}
		case 2:
			selected = map[string]string{"r0_e0": "direct", "r0_e1": "unrelated", "r1_e0": "unrelated", "r1_e1": "limited"}
		case 3:
			b, _ := json.Marshal(req.Questions["L002"].Instructions)
			if !strings.Contains(string(b), "Basic cloud usage") {
				t.Fatal("assessment did not receive evidence links")
			}
			selected = map[string]string{"L001": "direct", "L002": "large_gap"}
		case 4:
			b, _ := json.Marshal(req.State)
			if !strings.Contains(string(b), `"match":"Weak"`) {
				t.Fatal("overall did not receive actual assessments")
			}
			selected = map[string]string{"overall": "Strong"}
		case 5:
			b, _ := json.Marshal(req.State)
			if !strings.Contains(string(b), `"overallMatch":"Strong"`) || !strings.Contains(string(b), `"category":"preferred"`) {
				t.Fatal("trace missing selected overall or source type")
			}
			selected = map[string]string{"L001": "Supporting", "L002": "Non-decisive"}
		default:
			t.Fatal("unexpected call")
		}
		answers := map[string]answer{}
		for id, value := range selected {
			if _, ok := req.Questions[id].Criteria[value]; !ok {
				t.Fatalf("missing option %s:%s", id, value)
			}
			answers[id] = answer{"choice", value}
		}
		b, _ := json.Marshal(response{Model: model, Answers: answers})
		return reply(string(b), 200), nil
	})
	result, err := evaluator.Evaluate(context.Background(), c, p, j)
	if err != nil {
		t.Fatal(err)
	}
	if stages != 5 || result.OverallMatch != evaluator.Strong || len(result.Requirements) != 2 {
		t.Fatalf("%+v calls=%d", result, stages)
	}
	if result.Requirements[1].Match != evaluator.Weak || len(result.Requirements[1].MissingEvidence) == 0 || len(result.Requirements[1].Unknowns) == 0 {
		t.Fatal("gap mapping failed")
	}
	if result.Requirements[0].SupportingEvidence[0].Text != p.Evidence[0].Text {
		t.Fatal("evidence invented")
	}
	b, _ := json.Marshal(result)
	for _, forbidden := range []string{"confidence", "probabilities", "hiringProbability", "applicationRecommendation"} {
		if strings.Contains(string(b), forbidden) {
			t.Fatalf("leaked provider output: %s", forbidden)
		}
	}
	if string(b) == "" || !strings.Contains(string(b), `"unknown":[]`) {
		t.Fatal("empty arrays must be JSON arrays")
	}
}

func TestNoRequirementsStopsBeforeAssessment(t *testing.T) {
	calls := 0
	c := testClient(t, func(*http.Request) (*http.Response, error) {
		calls++
		return reply(`{"model":"m","answers":{"L001":{"type":"choice","choice":"ignore"}}}`, 200), nil
	})
	_, err := evaluator.Evaluate(context.Background(), c, evaluator.Profile{Evidence: []evaluator.Evidence{{ID: "e", Text: "API"}}}, evaluator.Job{Text: "Salary only", Lines: []string{"Salary only"}})
	if err == nil || calls != 1 || !strings.Contains(err.Error(), "no capability requirements") {
		t.Fatalf("%d %v", calls, err)
	}
}

func TestNoLinkedEvidenceRemainsUnknown(t *testing.T) {
	calls := 0
	c := testClient(t, func(*http.Request) (*http.Response, error) {
		calls++
		switch calls {
		case 1:
			return reply(`{"model":"m","answers":{"L001":{"type":"choice","choice":"required"}}}`, 200), nil
		case 2:
			return reply(`{"model":"m","answers":{"r0_e0":{"type":"choice","choice":"unrelated"}}}`, 200), nil
		case 3:
			return reply(`{"model":"m","answers":{"overall":{"type":"choice","choice":"Unknown"}}}`, 200), nil
		case 4:
			return reply(`{"model":"m","answers":{"L001":{"type":"choice","choice":"Limiting"}}}`, 200), nil
		default:
			t.Fatal("unexpected request")
			return nil, nil
		}
	})
	result, err := evaluator.Evaluate(context.Background(), c, evaluator.Profile{Evidence: []evaluator.Evidence{{ID: "e", Text: "API design"}}}, evaluator.Job{Text: "Kubernetes operations", Lines: []string{"Kubernetes operations"}})
	if err != nil {
		t.Fatal(err)
	}
	r := result.Requirements[0]
	if calls != 4 || r.Match != evaluator.Unknown || len(r.SupportingEvidence) != 0 || len(r.MissingEvidence) == 0 || len(r.Unknowns) == 0 || !strings.Contains(r.Reasoning, "application") {
		t.Fatalf("%+v calls=%d", r, calls)
	}
}
