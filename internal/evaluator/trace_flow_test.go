package evaluator

import (
	"context"
	"errors"
	"testing"
)

type traceEngine struct {
	refs     []InfluenceReference
	err      error
	selected MatchLevel
}

func (e *traceEngine) Extract(context.Context, Job) ([]Requirement, error) {
	return []Requirement{{ID: "L001", Text: "Azure", Category: "preferred"}}, nil
}
func (e *traceEngine) Assess(_ context.Context, _ Profile, _ Job, reqs []Requirement) ([]Assessment, error) {
	return []Assessment{{Requirement: reqs[0], Match: Unknown}}, nil
}
func (e *traceEngine) Overall(context.Context, Profile, Job, []Assessment) (MatchLevel, string, error) {
	return Partial, "unchanged overall reasoning", nil
}
func (e *traceEngine) Trace(_ context.Context, _ Profile, _ Job, _ []Assessment, selected MatchLevel) ([]InfluenceReference, error) {
	e.selected = selected
	return e.refs, e.err
}

func TestTraceDoesNotReviseOverallAndFailsClosed(t *testing.T) {
	for _, tc := range []struct {
		name string
		refs []InfluenceReference
		err  error
		fail bool
	}{
		{"non-decisive", []InfluenceReference{{"L001", NonDecisive}}, nil, false},
		{"preferred limiting", []InfluenceReference{{"L001", Limiting}}, nil, false},
		{"contradictory", []InfluenceReference{{"L001", Supporting}}, nil, true},
		{"missing", nil, nil, true},
		{"provider failed", nil, errors.New("trace unavailable"), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := &traceEngine{refs: tc.refs, err: tc.err}
			r, err := Evaluate(context.Background(), e, Profile{Evidence: []Evidence{{ID: "E1", Text: "API"}}}, Job{Text: "Azure", Lines: []string{"Azure"}})
			if (err != nil) != tc.fail {
				t.Fatalf("unexpected error %v", err)
			}
			if e.selected != Partial {
				t.Fatal("trace did not receive actual overall")
			}
			if !tc.fail && (r.OverallMatch != Partial || r.OverallReasoning != "unchanged overall reasoning") {
				t.Fatal("trace revised overall")
			}
			if tc.fail && r.OverallMatch != "" {
				t.Fatal("invalid trace exposed as successful result")
			}
		})
	}
}
