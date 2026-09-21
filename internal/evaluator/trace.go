package evaluator

import "fmt"

type Influence string

const (
	Supporting  Influence = "Supporting"
	Limiting    Influence = "Limiting"
	NonDecisive Influence = "Non-decisive"
)

// References are resolved against assessments; the provider cannot rewrite facts.
type InfluenceReference struct {
	RequirementID string
	Influence     Influence
}

type TraceEntry struct {
	RequirementID string     `json:"requirementId"`
	Requirement   string     `json:"requirement"`
	Match         MatchLevel `json:"match"`
	SourceType    string     `json:"sourceType"`
}

type DecisionTrace struct {
	Source      string       `json:"source"`
	Supporting  []TraceEntry `json:"supporting"`
	Limiting    []TraceEntry `json:"limiting"`
	NonDecisive []TraceEntry `json:"nonDecisive"`
}

func ResolveTrace(overall MatchLevel, rows []Assessment, refs []InfluenceReference) (DecisionTrace, error) {
	if _, err := ParseMatchLevel(string(overall)); err != nil {
		return DecisionTrace{}, err
	}
	t := DecisionTrace{Source: "jev-post-hoc-attribution; not a causal record of the original decision", Supporting: []TraceEntry{}, Limiting: []TraceEntry{}, NonDecisive: []TraceEntry{}}
	byID := map[string]Assessment{}
	for _, r := range rows {
		if _, exists := byID[r.ID]; exists || r.ID == "" {
			return DecisionTrace{}, fmt.Errorf("invalid or duplicate requirement identity")
		}
		byID[r.ID] = r
	}
	if len(refs) != len(rows) {
		return DecisionTrace{}, fmt.Errorf("incomplete decision trace")
	}
	roles := map[string]Influence{}
	for _, ref := range refs {
		r, exists := byID[ref.RequirementID]
		if !exists {
			return DecisionTrace{}, fmt.Errorf("unknown trace requirement")
		}
		if _, duplicate := roles[ref.RequirementID]; duplicate {
			return DecisionTrace{}, fmt.Errorf("duplicate trace reference")
		}
		switch ref.Influence {
		case Supporting:
			if r.Match == Unknown || len(r.SupportingEvidence) == 0 {
				return DecisionTrace{}, fmt.Errorf("supporting influence without assessable evidence")
			}
		case Limiting:
			if overall == Strong {
				return DecisionTrace{}, fmt.Errorf("Limiting influence is inconsistent with Strong overall")
			}
		case NonDecisive:
		default:
			return DecisionTrace{}, fmt.Errorf("invalid overall influence")
		}
		roles[ref.RequirementID] = ref.Influence
	}
	// Stable source order, regardless of provider answer order.
	for _, r := range rows {
		e := TraceEntry{r.ID, r.Text, r.Match, r.Category}
		switch roles[r.ID] {
		case Supporting:
			t.Supporting = append(t.Supporting, e)
		case Limiting:
			t.Limiting = append(t.Limiting, e)
		case NonDecisive:
			t.NonDecisive = append(t.NonDecisive, e)
		}
	}
	return t, nil
}
