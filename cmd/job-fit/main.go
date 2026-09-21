package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"job-fit/internal/application"
	"job-fit/internal/evaluator"
	"job-fit/internal/jev"
)

type namedEngine interface {
	evaluator.Engine
	Identity() (string, string)
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	return runWithEngine(ctx, args, stdout, stderr, getenv, func(key string) (namedEngine, error) { return jev.New(key) })
}

func runWithEngine(ctx context.Context, args []string, stdout, stderr io.Writer, getenv func(string) string, newEngine func(string) (namedEngine, error)) int {
	flags := flag.NewFlagSet("job-fit", flag.ContinueOnError)
	flags.SetOutput(stderr)
	profilePath := flags.String("profile", "", "canonical candidate profile JSON file (required)")
	jobPath := flags.String("job", "", "raw UTF-8 JD text file (required; one item per line)")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	fail := func(err error) int { fmt.Fprintln(stderr, "job-fit:", err); return 1 }
	if *profilePath == "" || *jobPath == "" || flags.NArg() != 0 {
		return fail(errors.New("use --profile <profile.json> --job <job.txt>; positional arguments are not supported"))
	}
	snapshot, err := application.LoadSnapshot(*profilePath)
	if err != nil {
		return fail(err)
	}
	job, err := evaluator.LoadJob(*jobPath)
	if err != nil {
		return fail(err)
	}
	client, err := newEngine(getenv("TYPESAFE_API_KEY"))
	if err != nil {
		return fail(err)
	}
	name, model := client.Identity()
	service, err := application.WithSnapshot(snapshot, client, name, model)
	if err != nil {
		return fail(err)
	}
	output, err := service.Evaluate(ctx, job.Text)
	if err != nil {
		return fail(err)
	}
	result := output.Result
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fail(errors.New("cannot write JSON output"))
	}
	return 0
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	os.Exit(run(ctx, os.Args[1:], os.Stdout, os.Stderr, os.Getenv))
}
