package jev

import (
	"os"
	"strings"
	"testing"

	"job-fit/internal/evaluator"
)

// These live verification inputs do not use duration eligibility. Negative
// wording hardening therefore does not alter their assessment choice sets.
func TestVerificationInputsHaveNoDurationOption(t *testing.T) {
	for _, name := range []string{"core", "nitori"} {
		b, err := os.ReadFile("../../samples/" + name + ".txt")
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(b), "\n") {
			if hasExplicitVersionOrDuration(line) {
				t.Fatalf("unexpected duration eligibility in %s: %s", name, line)
			}
		}
	}
}

func TestQualificationEligibility(t *testing.T) {
	for _, tc := range []struct {
		text string
		want bool
	}{
		{"Angular / TypeScript development experience", false},
		{"AWS certification", false},
		{"Ownership of large-scale production services", false},
		{"No minimum years of experience required", false},
		{"Experience with version 8 is not required", false},
		{"Experience with version 5.0 is not required", false},
		{"Version 8 experience is not mandatory", false},
		{"Version or duration unknown", false},
		{"Java 8+ development for at least one year", true},
		{"At least three years of AWS infrastructure construction and operations", true},
		{"Java 8 or later", true},
		{"TypeScript version 5", true},
		{"Python 3.12 experience", true},
		{"At least 18 months of experience", true},
		{"Years of experience with Angular", true},
		{"Java開発3年以上", true},
		{"Java 개발 3년 이상", true},
		{"Java経験年数不問", false},
		{"경력 기간 무관", false},
		{"No specific version required; at least two years of development", true},
	} {
		t.Run(tc.text, func(t *testing.T) {
			if got := hasExplicitVersionOrDuration(tc.text); got != tc.want {
				t.Fatalf("eligible=%v want %v", got, tc.want)
			}
		})
	}
}

func TestRationaleRejectsUngroundedDurationWithoutMutatingAssessment(t *testing.T) {
	row := evaluator.Assessment{Requirement: evaluator.Requirement{Text: "Angular / TypeScript development experience"}, SupportingEvidence: []evaluator.SupportingEvidence{{Evidence: evaluator.Evidence{ID: "E5", Text: "Angular and TypeScript. Duration not documented."}, Relation: "limited"}}}
	if err := applyRationale(&row, "duration"); err == nil {
		t.Fatal("unguarded synthetic gap")
	}
	if row.Match != "" || row.Reasoning != "" || len(row.MissingEvidence) != 0 || len(row.Unknowns) != 0 {
		t.Fatal("rejected rationale mutated result")
	}
	row.Text = "Java 8+ development for at least one year"
	if err := applyRationale(&row, "duration"); err != nil {
		t.Fatal(err)
	}
	if row.Match != evaluator.Partial || len(row.MissingEvidence) != 1 {
		t.Fatal("explicit qualification not evaluated")
	}
}
