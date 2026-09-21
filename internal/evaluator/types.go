package evaluator

import "fmt"

type MatchLevel string

const (
	Strong  MatchLevel = "Strong"
	Partial MatchLevel = "Partial"
	Weak    MatchLevel = "Weak"
	Unknown MatchLevel = "Unknown"
)

func ParseMatchLevel(s string) (MatchLevel, error) {
	switch MatchLevel(s) {
	case Strong, Partial, Weak, Unknown:
		return MatchLevel(s), nil
	default:
		return "", fmt.Errorf("invalid match level")
	}
}

type Evidence struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Profile struct {
	Evidence []Evidence `json:"evidence"`
}

type Job struct {
	Text  string   `json:"text"`
	Lines []string `json:"lines"`
}

type Requirement struct {
	ID       string `json:"id"`
	Text     string `json:"requirement"`
	Category string `json:"category"`
}

type SupportingEvidence struct {
	Evidence
	Relation string `json:"relation"`
}

type Assessment struct {
	Requirement
	Match              MatchLevel           `json:"match"`
	SupportingEvidence []SupportingEvidence `json:"supportingEvidence"`
	MissingEvidence    []string             `json:"missingEvidence"`
	Unknowns           []string             `json:"unknowns"`
	Reasoning          string               `json:"reasoning"`
}

type Summary struct {
	Strong  []string `json:"strong"`
	Partial []string `json:"partial"`
	Weak    []string `json:"weak"`
	Unknown []string `json:"unknown"`
}

type EvaluationResult struct {
	DecisionTrace     DecisionTrace `json:"decisionTrace"`
	OverallMatch      MatchLevel    `json:"overallMatch"`
	OverallReasoning  string        `json:"overallReasoning"`
	Requirements      []Assessment  `json:"requirements"`
	Summary           Summary       `json:"summary"`
	ExplanationSource string        `json:"explanationSource"`
}
