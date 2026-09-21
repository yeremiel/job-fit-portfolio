package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestStartupFailures(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		key  string
		code int
	}{
		{"help", []string{"--help"}, "", 0},
		{"missing profile", nil, "fake", 1},
		{"bad profile", []string{"--profile", "/private-profile-do-not-log"}, "fake", 1},
		{"missing key", []string{"--profile", "../../data/profile.json"}, "", 1},
		{"invalid timeout", []string{"--profile", "../../data/profile.json", "--evaluation-timeout", "0s"}, "fake", 1},
		{"public bind", []string{"--profile", "../../data/profile.json", "--listen", "0.0.0.0:8080"}, "fake", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var log bytes.Buffer
			code := run(context.Background(), tc.args, &log, func(string) string { return tc.key })
			if code != tc.code {
				t.Fatal(code, log.String())
			}
			if strings.Contains(log.String(), "private-profile-do-not-log") || strings.Contains(log.String(), "fake") {
				t.Fatal("unsafe log", log.String())
			}
		})
	}
}
