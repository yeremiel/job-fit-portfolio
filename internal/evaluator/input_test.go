package evaluator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func inputFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadProfile(t *testing.T) {
	valid := `{"evidence":[{"id":"E1","text":"API design"}]}`
	p, err := LoadProfile(inputFile(t, valid))
	if err != nil || len(p.Evidence) != 1 || p.Evidence[0].Text != "API design" {
		t.Fatalf("%+v %v", p, err)
	}
	for _, content := range []string{`{`, `null`, `[]`, `{}`, `{"evidence":[]}`, valid + ` {}`, `{"evidence":[{"id":"E1","text":""}]}`, `{"evidence":[{"id":"E1","text":"API"},{"id":"E1","text":"DB"}]}`, `{"evidence":[],"resume":"ignored"}`, strings.Repeat(" ", maxInputBytes+1)} {
		if _, err := LoadProfile(inputFile(t, content)); err == nil {
			t.Errorf("expected invalid profile error for %q", content[:min(len(content), 80)])
		}
	}
	if _, err := LoadProfile(filepath.Join(t.TempDir(), "missing")); err == nil || !strings.Contains(err.Error(), "profile") {
		t.Fatalf("missing file error: %v", err)
	}
}

func TestLoadJob(t *testing.T) {
	j, err := LoadJob(inputFile(t, "\r\nRequired: Java\r\n\r\nPreferred: Azure\n"))
	if err != nil || len(j.Lines) != 2 || j.Lines[0] != "Required: Java" {
		t.Fatalf("%+v %v", j, err)
	}
	for _, content := range []string{"", " \n\t", string([]byte{0xff}), "text\x00", strings.Repeat("x\n", 65), strings.Repeat("x", maxInputBytes+1)} {
		if _, err := LoadJob(inputFile(t, content)); err == nil {
			t.Errorf("expected invalid job error")
		}
	}
	if _, err := LoadJob(filepath.Join(t.TempDir(), "missing")); err == nil || !strings.Contains(err.Error(), "job") {
		t.Fatalf("missing file error: %v", err)
	}
}

func TestParseMatchLevel(t *testing.T) {
	for _, s := range []string{"Strong", "Partial", "Weak", "Unknown"} {
		if level, err := ParseMatchLevel(s); err != nil || string(level) != s {
			t.Fatalf("%s %v", s, err)
		}
	}
	for _, s := range []string{"", "strong", "HIGH", "82%"} {
		if _, err := ParseMatchLevel(s); err == nil {
			t.Fatalf("accepted %q", s)
		}
	}
}

func TestCheckedInInputs(t *testing.T) {
	p, err := LoadProfile("../../data/profile.json")
	if err != nil {
		t.Fatal(err)
	}
	j, err := LoadJob("../../samples/proptech-plus.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Evidence) != 7 || !strings.Contains(j.Text, "Java 8") {
		t.Fatal("unexpected sample input")
	}
	for _, forbidden := range []string{"BlueMeme", "Nitori", "Prop Tech", "Overall Expected Match"} {
		if strings.Contains(Policy, forbidden) {
			t.Fatalf("baseline leaked into policy: %s", forbidden)
		}
	}
}
