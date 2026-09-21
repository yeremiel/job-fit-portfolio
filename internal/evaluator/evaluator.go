package evaluator

import (
	"context"
	_ "embed"
	"fmt"

	"job-fit/internal/failure"
)

//go:embed instructions.md
var Policy string

// Engine is the application boundary; no provider request/response types escape it.
type Engine interface {
	Extract(context.Context, Job) ([]Requirement, error)
	Assess(context.Context, Profile, Job, []Requirement) ([]Assessment, error)
	Overall(context.Context, Profile, Job, []Assessment) (MatchLevel, string, error)
	Trace(context.Context, Profile, Job, []Assessment, MatchLevel) ([]InfluenceReference, error)
}

func Evaluate(ctx context.Context, engine Engine, profile Profile, job Job) (EvaluationResult, error) {
	if err := profile.Validate(); err != nil {
		return EvaluationResult{}, err
	}
	if len(job.Lines) == 0 {
		return EvaluationResult{}, fmt.Errorf("job text is empty")
	}
	if err := ctx.Err(); err != nil {
		return EvaluationResult{}, err
	}
	reqs, err := engine.Extract(ctx, job)
	if err != nil {
		return EvaluationResult{}, fmt.Errorf("extract requirements: %w", err)
	}
	if len(reqs) == 0 {
		return EvaluationResult{}, failure.New(failure.NoRequirements, "no capability requirements identified in job text")
	}
	if err := ctx.Err(); err != nil {
		return EvaluationResult{}, err
	}
	rows, err := engine.Assess(ctx, profile, job, reqs)
	if err != nil {
		return EvaluationResult{}, fmt.Errorf("assess requirements: %w", err)
	}
	if len(rows) != len(reqs) {
		return EvaluationResult{}, fmt.Errorf("incomplete requirement assessment")
	}
	summary := Summary{Strong: []string{}, Partial: []string{}, Weak: []string{}, Unknown: []string{}}
	for i, r := range rows {
		if r.Requirement != reqs[i] {
			return EvaluationResult{}, fmt.Errorf("assessment does not match source requirement")
		}
		if _, err := ParseMatchLevel(string(r.Match)); err != nil {
			return EvaluationResult{}, err
		}
		switch r.Match {
		case Strong:
			summary.Strong = append(summary.Strong, r.ID)
		case Partial:
			summary.Partial = append(summary.Partial, r.ID)
		case Weak:
			summary.Weak = append(summary.Weak, r.ID)
		case Unknown:
			summary.Unknown = append(summary.Unknown, r.ID)
		}
	}
	if err := ctx.Err(); err != nil {
		return EvaluationResult{}, err
	}
	level, reason, err := engine.Overall(ctx, profile, job, rows)
	if err != nil {
		return EvaluationResult{}, fmt.Errorf("assess overall: %w", err)
	}
	if _, err := ParseMatchLevel(string(level)); err != nil {
		return EvaluationResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return EvaluationResult{}, err
	}
	refs, err := engine.Trace(ctx, profile, job, rows, level)
	if err != nil {
		return EvaluationResult{}, fmt.Errorf("attribute overall: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return EvaluationResult{}, err
	}
	trace, err := ResolveTrace(level, rows, refs)
	if err != nil {
		return EvaluationResult{}, fmt.Errorf("validate decision trace: %w", err)
	}
	return EvaluationResult{DecisionTrace: trace, OverallMatch: level, OverallReasoning: reason, Requirements: rows, Summary: summary, ExplanationSource: "application-rendered rationale from provider choices or absence of linked evidence; evidence copied from profile"}, nil
}
