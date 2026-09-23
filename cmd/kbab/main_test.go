package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpListsCreate(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"--help"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "create") {
		t.Errorf("expected help to list create command, got:\n%s", out.String())
	}
}

func TestCreateDefaultPath(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"create"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "create: .") {
		t.Errorf("expected default path '.', got:\n%s", out.String())
	}
}
