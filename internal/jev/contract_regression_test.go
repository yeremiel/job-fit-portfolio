package jev

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"job-fit/internal/evaluator"
)

func TestStrongTraceExcludesLimiting(t *testing.T) {
	for _, returned := range []string{"Non-decisive", "Limiting"} {
		t.Run(returned, func(t *testing.T) {
			c := testClient(t, func(r *http.Request) (*http.Response, error) {
				var req request
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatal(err)
				}
				if _, exists := req.Questions["L011"].Criteria["Limiting"]; exists {
					t.Error("Limiting offered for Strong overall")
				}
				if _, exists := req.Questions["L011"].Criteria["Non-decisive"]; !exists {
					t.Error("Non-decisive must remain available")
				}
				return reply(`{"model":"m","answers":{"L011":{"type":"choice","choice":"`+returned+`"}}}`, 200), nil
			})
			rows := []evaluator.Assessment{{Requirement: evaluator.Requirement{ID: "L011", Text: "Maintenance", Category: "responsibility"}, Match: evaluator.Partial}}
			_, err := c.Trace(context.Background(), evaluator.Profile{}, evaluator.Job{}, rows, evaluator.Strong)
			if (err != nil) != (returned == "Limiting") {
				t.Fatalf("unexpected trace validation: %v", err)
			}
		})
	}
}

func TestAssessmentQualificationGrounding(t *testing.T) {
	for _, tc := range []struct {
		name, requirement, selected string
		allowed, fail               bool
	}{
		{"implicit", "Angular / TypeScript development experience", "depth", false, false},
		{"reject ungrounded answer", "Angular / TypeScript development experience", "duration", false, true},
		{"explicit", "Java 8+ development for at least one year", "duration", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			c := testClient(t, func(r *http.Request) (*http.Response, error) {
				calls++
				var req request
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatal(err)
				}
				if calls == 1 {
					return reply(`{"model":"m","answers":{"r0_e0":{"type":"choice","choice":"limited"}}}`, 200), nil
				}
				if calls != 2 {
					t.Fatal("unexpected extra call")
				}
				instruction, _ := json.Marshal(req.Questions["L010"].Instructions)
				for _, boundary := range []string{"explicit JD condition", "directly necessary", "certification", "scale", "ownership level"} {
					if !strings.Contains(string(instruction), boundary) {
						t.Fatalf("missing qualification boundary: %s", boundary)
					}
				}
				_, exists := req.Questions["L010"].Criteria["duration"]
				if exists != tc.allowed {
					t.Errorf("duration available=%v, want %v", exists, tc.allowed)
				}
				return reply(`{"model":"m","answers":{"L010":{"type":"choice","choice":"`+tc.selected+`"}}}`, 200), nil
			})
			p := evaluator.Profile{Evidence: []evaluator.Evidence{{ID: "E5", Text: "Angular, TypeScript and Java experience. Individual durations are not documented."}}}
			rows, err := c.Assess(context.Background(), p, evaluator.Job{Text: tc.requirement + "\nOther requirement: Java version 17 for three years."}, []evaluator.Requirement{{ID: "L010", Text: tc.requirement, Category: "stack"}})
			if (err != nil) != tc.fail {
				t.Fatalf("unexpected assessment error: %v", err)
			}
			if err != nil {
				return
			}
			b, _ := json.Marshal(rows)
			if !tc.allowed && (strings.Contains(string(b), "specified version") || strings.Contains(string(b), "duration unknown")) {
				t.Fatal("synthetic qualification generated")
			}
			if tc.allowed && !strings.Contains(string(b), "specified version or duration") {
				t.Fatal("explicit qualification lost")
			}
		})
	}
}
