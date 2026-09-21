package httpapi

import (
	"fmt"
	"job-fit/internal/application"
	"job-fit/internal/evaluator"
	"reflect"
)

type evidenceDTO struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Relation string `json:"relation"`
}
type requirementDTO struct {
	ID                 string               `json:"id"`
	Requirement        string               `json:"requirement"`
	SourceType         string               `json:"sourceType"`
	Match              evaluator.MatchLevel `json:"match"`
	SupportingEvidence []evidenceDTO        `json:"supportingEvidence"`
	MissingEvidence    []string             `json:"missingEvidence"`
	Unknowns           []string             `json:"unknowns"`
	Reasoning          string               `json:"reasoning"`
}
type traceDTO struct {
	Source      string   `json:"source"`
	Supporting  []string `json:"supporting"`
	Limiting    []string `json:"limiting"`
	NonDecisive []string `json:"nonDecisive"`
}
type resultDTO struct {
	Job               JobMetadata          `json:"job"`
	OverallMatch      evaluator.MatchLevel `json:"overallMatch"`
	OverallReasoning  string               `json:"overallReasoning"`
	Requirements      []requirementDTO     `json:"requirements"`
	DecisionTrace     traceDTO             `json:"decisionTrace"`
	ExplanationSource string               `json:"explanationSource"`
	Metadata          application.Metadata `json:"metadata"`
}

func present(out application.Output, job JobMetadata) (resultDTO, error) {
	r := out.Result
	dto := resultDTO{Job: job, OverallMatch: r.OverallMatch, OverallReasoning: r.OverallReasoning, Requirements: []requirementDTO{}, DecisionTrace: traceDTO{r.DecisionTrace.Source, []string{}, []string{}, []string{}}, ExplanationSource: r.ExplanationSource, Metadata: out.Metadata}
	if len(r.Requirements) == 0 {
		return dto, fmt.Errorf("missing result requirements")
	}
	byID := map[string]evaluator.Assessment{}
	for _, row := range r.Requirements {
		if _, err := evaluator.ParseMatchLevel(string(row.Match)); err != nil {
			return dto, err
		}
		switch row.Category {
		case "required", "preferred", "responsibility", "stack", "context":
		default:
			return dto, fmt.Errorf("invalid source type")
		}
		byID[row.ID] = row
		ev := []evidenceDTO{}
		for _, e := range row.SupportingEvidence {
			switch e.Relation {
			case "direct", "transferable", "limited":
			default:
				return dto, fmt.Errorf("invalid evidence relation")
			}
			ev = append(ev, evidenceDTO{e.ID, e.Text, e.Relation})
		}
		dto.Requirements = append(dto.Requirements, requirementDTO{row.ID, row.Text, row.Category, row.Match, ev, append([]string{}, row.MissingEvidence...), append([]string{}, row.Unknowns...), row.Reasoning})
	}
	refs := []evaluator.InfluenceReference{}
	for _, group := range []struct {
		entries []evaluator.TraceEntry
		role    evaluator.Influence
		ids     *[]string
	}{{r.DecisionTrace.Supporting, evaluator.Supporting, &dto.DecisionTrace.Supporting}, {r.DecisionTrace.Limiting, evaluator.Limiting, &dto.DecisionTrace.Limiting}, {r.DecisionTrace.NonDecisive, evaluator.NonDecisive, &dto.DecisionTrace.NonDecisive}} {
		for _, entry := range group.entries {
			row, ok := byID[entry.RequirementID]
			if !ok || entry.Requirement != row.Text || entry.Match != row.Match || entry.SourceType != row.Category {
				return dto, fmt.Errorf("inconsistent trace entry")
			}
			refs = append(refs, evaluator.InfluenceReference{RequirementID: entry.RequirementID, Influence: group.role})
			*group.ids = append(*group.ids, entry.RequirementID)
		}
	}
	validated, err := evaluator.ResolveTrace(r.OverallMatch, r.Requirements, refs)
	if err != nil {
		return dto, err
	}
	// Reject out-of-order entries rather than silently repair them at the HTTP boundary.
	for i, a := range [][]evaluator.TraceEntry{validated.Supporting, validated.Limiting, validated.NonDecisive} {
		b := [][]evaluator.TraceEntry{r.DecisionTrace.Supporting, r.DecisionTrace.Limiting, r.DecisionTrace.NonDecisive}[i]
		if len(a) != len(b) || (len(a) > 0 && !reflect.DeepEqual(a, b)) {
			return dto, fmt.Errorf("invalid trace order")
		}
	}
	return dto, nil
}
