package application

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"job-fit/internal/evaluator"
	"job-fit/internal/failure"
)

type fakeEngine struct {
	cancel     context.CancelFunc
	extractErr error
}

func (f fakeEngine) Extract(c context.Context, j evaluator.Job) ([]evaluator.Requirement, error) {
	if f.cancel != nil {
		f.cancel()
	}
	if f.extractErr != nil {
		return nil, f.extractErr
	}
	return []evaluator.Requirement{{ID: "L001", Text: j.Lines[0], Category: "required"}}, nil
}
func (f fakeEngine) Assess(_ context.Context, p evaluator.Profile, _ evaluator.Job, r []evaluator.Requirement) ([]evaluator.Assessment, error) {
	evidence := p.Evidence[0]
	p.Evidence[0].Text = "mutated by fake engine"
	return []evaluator.Assessment{{Requirement: r[0], Match: evaluator.Strong, SupportingEvidence: []evaluator.SupportingEvidence{{Evidence: evidence, Relation: "direct"}}, MissingEvidence: []string{}, Unknowns: []string{}, Reasoning: "Direct."}}, nil
}
func (f fakeEngine) Overall(context.Context, evaluator.Profile, evaluator.Job, []evaluator.Assessment) (evaluator.MatchLevel, string, error) {
	return evaluator.Strong, "Direct.", nil
}
func (f fakeEngine) Trace(context.Context, evaluator.Profile, evaluator.Job, []evaluator.Assessment, evaluator.MatchLevel) ([]evaluator.InfluenceReference, error) {
	return []evaluator.InfluenceReference{{RequirementID: "L001", Influence: evaluator.Supporting}}, nil
}
func profileFile(t *testing.T) (string, []byte) {
	t.Helper()
	b := []byte("{\"evidence\":[{\"id\":\"E1\",\"text\":\"API design\"}]}\n")
	p := filepath.Join(t.TempDir(), "profile.json")
	if err := os.WriteFile(p, b, 0600); err != nil {
		t.Fatal(err)
	}
	return p, b
}
func TestSnapshotVersionAndConcurrentReuse(t *testing.T) {
	path, b := profileFile(t)
	s, err := New(path, fakeEngine{}, "jev", "fake")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte("not a profile anymore"), 0600); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := s.Evaluate(context.Background(), " API design \r\n")
			if err != nil {
				t.Error(err)
				return
			}
			if out.Metadata.ProfileVersion != fmt.Sprintf("sha256:%x", sha256.Sum256(b)) || out.Result.Requirements[0].SupportingEvidence[0].Text != "API design" {
				t.Error("snapshot mutated or version mismatch")
			}
		}()
	}
	wg.Wait()
	if _, err := New(path, fakeEngine{}, "jev", "fake"); failure.Classify(err) != failure.Configuration {
		t.Fatal("invalid startup accepted", err)
	}
}
func TestServiceAndDirectEvaluatorParity(t *testing.T) {
	path, _ := profileFile(t)
	s, err := New(path, fakeEngine{}, "jev", "fake")
	if err != nil {
		t.Fatal(err)
	}
	text := " \r\nRequired: API design\r\n\n"
	j, err := evaluator.PrepareJob(text)
	if err != nil {
		t.Fatal(err)
	}
	p, err := evaluator.LoadProfile(path)
	if err != nil {
		t.Fatal(err)
	}
	direct, err := evaluator.Evaluate(context.Background(), fakeEngine{}, p, j)
	if err != nil {
		t.Fatal(err)
	}
	shared, err := s.Evaluate(context.Background(), text)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(direct, shared.Result) {
		t.Fatal("workflow changed")
	}
	// CLI can still encode the original domain result, without an HTTP envelope.
	a, _ := json.Marshal(direct)
	b, _ := json.Marshal(shared.Result)
	if string(a) != string(b) {
		t.Fatal("CLI schema changed")
	}
}
func TestCancellationStopsNextStage(t *testing.T) {
	path, _ := profileFile(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s, err := New(path, fakeEngine{cancel: cancel}, "jev", "fake")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Evaluate(ctx, "API design")
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
func TestServiceClassifiesProviderAndInputErrors(t *testing.T) {
	path, _ := profileFile(t)
	for _, tc := range []struct {
		err  error
		want failure.Kind
	}{{errors.New("raw provider fault"), failure.Evaluation}, {failure.New(failure.Unavailable, "busy"), failure.Unavailable}} {
		s, err := New(path, fakeEngine{extractErr: tc.err}, "jev", "fake")
		if err != nil {
			t.Fatal(err)
		}
		_, err = s.Evaluate(context.Background(), "API design")
		if failure.Classify(err) != tc.want {
			t.Fatal(err)
		}
	}
}
