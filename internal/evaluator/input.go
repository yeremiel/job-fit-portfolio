package evaluator

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"job-fit/internal/failure"
)

// These are request-size safeguards, not evaluation thresholds.
const maxInputBytes = 16 * 1024

func readInput(path, kind string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s file: %w", kind, err)
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, maxInputBytes+1))
	if err != nil {
		return nil, fmt.Errorf("cannot read %s file", kind)
	}
	if len(b) > maxInputBytes {
		return nil, fmt.Errorf("%s exceeds 16 KiB input limit", kind)
	}
	if !utf8.Valid(b) || bytes.ContainsRune(b, '\x00') {
		return nil, fmt.Errorf("%s must be UTF-8 text", kind)
	}
	return b, nil
}

func LoadProfile(path string) (Profile, error) {
	p, _, err := LoadProfileSnapshot(path)
	return p, err
}

// LoadProfileSnapshot hashes exactly the bytes parsed, with a single file read.
func LoadProfileSnapshot(path string) (Profile, string, error) {
	b, err := readInput(path, "profile")
	if err != nil {
		return Profile{}, "", err
	}
	p, err := parseProfile(b)
	if err != nil {
		return Profile{}, "", err
	}
	return p, fmt.Sprintf("sha256:%x", sha256.Sum256(b)), nil
}

func parseProfile(b []byte) (Profile, error) {
	var p Profile
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if d.Decode(&p) != nil {
		return Profile{}, fmt.Errorf("invalid profile JSON: expected an object with an evidence array of id/text objects")
	}
	if d.Decode(new(any)) != io.EOF {
		return Profile{}, fmt.Errorf("invalid profile JSON: expected one object")
	}
	if err := p.Validate(); err != nil {
		return Profile{}, err
	}
	return p, nil
}

func (p Profile) Validate() error {
	if len(p.Evidence) == 0 || len(p.Evidence) > 32 {
		return fmt.Errorf("profile must contain 1–32 evidence items")
	}
	seen := map[string]bool{}
	for _, e := range p.Evidence {
		if strings.TrimSpace(e.ID) == "" || strings.TrimSpace(e.Text) == "" {
			return fmt.Errorf("profile evidence requires non-empty id and text")
		}
		if seen[e.ID] {
			return fmt.Errorf("profile evidence IDs must be unique")
		}
		seen[e.ID] = true
	}
	return nil
}

func LoadJob(path string) (Job, error) {
	b, err := readInput(path, "job")
	if err != nil {
		return Job{}, err
	}
	return PrepareJob(string(b))
}

// PrepareJob is shared by file and HTTP inputs; it does not summarize or parse JD content.
func PrepareJob(text string) (Job, error) {
	if len(text) > maxInputBytes {
		return Job{}, failure.New(failure.TooLarge, "job exceeds 16 KiB input limit")
	}
	if !utf8.ValidString(text) || strings.ContainsRune(text, '\x00') {
		return Job{}, failure.New(failure.Invalid, "job must be UTF-8 text without NUL")
	}
	j := Job{Text: strings.TrimSpace(text)}
	for _, line := range strings.Split(j.Text, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			j.Lines = append(j.Lines, line)
		}
	}
	if len(j.Lines) == 0 {
		return Job{}, failure.New(failure.Invalid, "job text is empty")
	}
	if len(j.Lines) > 64 {
		return Job{}, failure.New(failure.TooLarge, "job exceeds 64 non-empty lines; prepare a concise JD with one item per line")
	}
	return j, nil
}
