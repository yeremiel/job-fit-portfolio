package evaluator

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestStoredCOREStrongLimitingIsRejected(t *testing.T) {
	b, err := os.ReadFile("../../docs/development/turn-008-results/core.json")
	if err != nil {
		t.Fatal(err)
	}
	var result EvaluationResult
	if err = json.Unmarshal(b, &result); err != nil {
		t.Fatal(err)
	}
	refs := []InfluenceReference{}
	for _, group := range []struct {
		entries []TraceEntry
		role    Influence
	}{{result.DecisionTrace.Supporting, Supporting}, {result.DecisionTrace.Limiting, Limiting}, {result.DecisionTrace.NonDecisive, NonDecisive}} {
		for _, e := range group.entries {
			refs = append(refs, InfluenceReference{e.RequirementID, group.role})
		}
	}
	before, _ := json.Marshal(result.Requirements)
	if _, err := ResolveTrace(result.OverallMatch, result.Requirements, refs); err == nil {
		t.Fatal("stored Strong + Limiting violation accepted")
	}
	// A valid alternative may retain the gap as Non-decisive; never force
	// every assessment to Supporting or alter its requirement-level match.
	for i := range refs {
		if refs[i].Influence == Limiting {
			refs[i].Influence = NonDecisive
		}
	}
	trace, err := ResolveTrace(result.OverallMatch, result.Requirements, refs)
	if err != nil || len(trace.Limiting) != 0 || len(trace.NonDecisive) == 0 {
		t.Fatalf("valid alternative rejected: %+v %v", trace, err)
	}
	after, _ := json.Marshal(result.Requirements)
	if !reflect.DeepEqual(before, after) || result.OverallMatch != Strong {
		t.Fatal("validation changed assessment or overall")
	}
	if _, err := ResolveTrace(Partial, result.Requirements, []InfluenceReference{}); err == nil {
		t.Fatal("missing references accepted")
	}
}
