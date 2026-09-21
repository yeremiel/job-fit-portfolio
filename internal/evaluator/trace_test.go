package evaluator

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResolveTrace(t *testing.T) {
	rows := []Assessment{
		{Requirement: Requirement{ID: "L007", Text: "API", Category: "responsibility"}, Match: Strong, SupportingEvidence: []SupportingEvidence{{Evidence: Evidence{ID: "E1", Text: "API design"}, Relation: "direct"}}},
		{Requirement: Requirement{ID: "L014", Text: "Azure", Category: "preferred"}, Match: Unknown},
		{Requirement: Requirement{ID: "L006", Text: "Java duration", Category: "required"}, Match: Partial},
	}
	valid := []InfluenceReference{{"L006", Limiting}, {"L014", NonDecisive}, {"L007", Supporting}}
	trace, err := ResolveTrace(Partial, rows, valid)
	if err != nil {
		t.Fatal(err)
	}
	if trace.Supporting[0].RequirementID != "L007" || trace.Supporting[0].Match != Strong || trace.Supporting[0].SourceType != "responsibility" || trace.Limiting[0].Requirement != "Java duration" || trace.NonDecisive[0].Match != Unknown {
		t.Fatalf("identity not preserved: %+v", trace)
	}
	b, err := json.Marshal(EvaluationResult{DecisionTrace: trace, Requirements: rows})
	if err != nil || !strings.Contains(string(b), `"decisionTrace"`) || !strings.Contains(string(b), `"requirementId":"L007"`) || !strings.Contains(string(b), `"supportingEvidence"`) {
		t.Fatalf("serialization: %s %v", b, err)
	}
	for name, refs := range map[string][]InfluenceReference{
		"missing":                valid[:2],
		"unknown id":             {{"missing", Limiting}, {"L014", NonDecisive}, {"L007", Supporting}},
		"duplicate":              {{"L006", Limiting}, {"L006", NonDecisive}, {"L007", Supporting}},
		"invalid influence":      {{"L006", "High"}, {"L014", NonDecisive}, {"L007", Supporting}},
		"unknown supporting":     {{"L006", Limiting}, {"L014", Supporting}, {"L007", Supporting}},
		"no evidence supporting": {{"L006", Supporting}, {"L014", NonDecisive}, {"L007", Supporting}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ResolveTrace(Partial, rows, refs); err == nil {
				t.Fatal("invalid trace accepted")
			}
		})
	}
	// Preferred Unknown may be Limiting; this observation must not be hidden.
	valid[1].Influence = Limiting
	trace, err = ResolveTrace(Partial, rows, valid)
	if err != nil || len(trace.Limiting) != 2 {
		t.Fatalf("preferred limiting suppressed: %+v %v", trace, err)
	}
	b, _ = json.Marshal(trace)
	if !strings.Contains(string(b), `"nonDecisive":[]`) {
		t.Fatal("empty list must serialize as array")
	}
	rows[1].ID = rows[0].ID
	if _, err := ResolveTrace(Partial, rows, valid); err == nil {
		t.Fatal("duplicate source identity accepted")
	}
}
