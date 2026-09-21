// Package application shares the evaluation workflow across local and HTTP callers.
package application

import (
	"context"
	"errors"
	"time"

	"job-fit/internal/evaluator"
	"job-fit/internal/failure"
)

type Metadata struct {
	ProfileVersion string `json:"profileVersion"`
	EvaluatedAt    string `json:"evaluatedAt"`
	Engine         string `json:"engine"`
	EngineModel    string `json:"engineModel"`
}
type Output struct {
	Result   evaluator.EvaluationResult
	Metadata Metadata
}
type Service struct {
	engine                     evaluator.Engine
	profile                    evaluator.Profile
	version, engineName, model string
}

// Snapshot is opaque to callers and contains the profile bytes' validated content/version.
type Snapshot struct {
	profile evaluator.Profile
	version string
}

func LoadSnapshot(path string) (Snapshot, error) {
	p, version, err := evaluator.LoadProfileSnapshot(path)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{p, version}, nil
}
func New(profilePath string, engine evaluator.Engine, name, model string) (*Service, error) {
	snapshot, err := LoadSnapshot(profilePath)
	if err != nil {
		return nil, failure.New(failure.Configuration, "cannot load canonical profile file")
	}
	return WithSnapshot(snapshot, engine, name, model)
}
func WithSnapshot(snapshot Snapshot, engine evaluator.Engine, name, model string) (*Service, error) {
	if engine == nil || name == "" || model == "" || snapshot.version == "" {
		return nil, failure.New(failure.Configuration, "evaluation engine or profile is not configured")
	}
	return &Service{engine: engine, profile: snapshot.profile, version: snapshot.version, engineName: name, model: model}, nil
}

func (s *Service) Evaluate(ctx context.Context, text string) (Output, error) {
	if err := ctx.Err(); err != nil {
		return Output{}, err
	}
	job, err := evaluator.PrepareJob(text)
	if err != nil {
		return Output{}, err
	}
	// Isolate the stored snapshot even from an Engine implementation that mutates slices.
	p := evaluator.Profile{Evidence: append([]evaluator.Evidence(nil), s.profile.Evidence...)}
	result, err := evaluator.Evaluate(ctx, s.engine, p, job)
	if err != nil {
		var categorized *failure.Error
		if errors.As(err, &categorized) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return Output{}, err
		}
		return Output{}, failure.New(failure.Evaluation, "evaluation failed validation or provider processing")
	}
	if err := ctx.Err(); err != nil {
		return Output{}, err
	}
	return Output{result, Metadata{s.version, time.Now().UTC().Format(time.RFC3339Nano), s.engineName, s.model}}, nil
}
