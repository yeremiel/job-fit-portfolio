package jev

import (
	"context"
	"encoding/json"
	"job-fit/internal/evaluator"
	"net/http"
	"strings"
	"testing"
)

func TestTraceMappingAndProviderErrors(t *testing.T) {
	rows := []evaluator.Assessment{{Requirement: evaluator.Requirement{ID: "L014", Text: "Azure", Category: "preferred"}, Match: evaluator.Unknown, Unknowns: []string{"Azure usage"}}}
	for _, influence := range []string{"Supporting", "Limiting", "Non-decisive", "invalid"} {
		t.Run(influence, func(t *testing.T) {
			c := testClient(t, func(r *http.Request) (*http.Response, error) {
				var req request
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatal(err)
				}
				b, _ := json.Marshal(req.State)
				for _, want := range []string{`"overallMatch":"Partial"`, `"category":"preferred"`, `"match":"Unknown"`, "Azure usage", "profile", "job", "evaluationPolicy"} {
					if !strings.Contains(string(b), want) {
						t.Fatalf("missing state %s", want)
					}
				}
				return reply(`{"model":"m","answers":{"L014":{"type":"choice","choice":"`+influence+`"}}}`, 200), nil
			})
			refs, err := c.Trace(context.Background(), evaluator.Profile{}, evaluator.Job{Text: "Azure"}, rows, evaluator.Partial)
			if influence == "invalid" {
				if err == nil {
					t.Fatal("invalid choice accepted")
				}
				return
			}
			if err != nil || len(refs) != 1 || refs[0].RequirementID != "L014" || string(refs[0].Influence) != influence {
				t.Fatalf("%+v %v", refs, err)
			}
			_, err = evaluator.ResolveTrace(evaluator.Partial, rows, refs)
			if (err != nil) != (influence == "Supporting") {
				t.Fatalf("consistency validation: %v", err)
			}
		})
	}
}
