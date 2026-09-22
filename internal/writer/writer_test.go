package writer_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/santiagozsh/kbab/internal/generator"
	"github.com/santiagozsh/kbab/internal/writer"
)

func TestWriteFilesCleanDirectory(t *testing.T) {
	tempDir := t.TempDir()

	files := generator.GeneratedFiles{
		Dockerfile:   "FROM golang:1.23-alpine AS base\n",
		Dockerignore: ".git\nbin/\n",
	}

	opts := writer.WriteOptions{
		TargetDir: tempDir,
	}

	err := writer.WriteFiles(files, opts)
	if err != nil {
		t.Fatalf("unexpected errro :%v", err)
	}

	dockerfilePath := filepath.Join(tempDir, "Dockerfile")
	dockerfileContent, err := os.ReadFile(dockerfilePath)
	if err != nil {
		t.Fatalf("failed to read created Dockerfile: %v", err)
	}

	if string(dockerfileContent) != files.Dockerfile {
		t.Errorf("expected Dockerfile content %q, got %q", files.Dockerfile, string(dockerfileContent))
	}
}

func TestWriteFilesDryRun(t *testing.T) {
	tempDir := t.TempDir()
	var out bytes.Buffer

	files := generator.GeneratedFiles{
		Dockerfile:   "FROM golang:1.23-alpine AS base\n",
		Dockerignore: ".git\n",
	}

	opts := writer.WriteOptions{
		TargetDir: tempDir,
		DryRun:    true,
		Out:       &out,
	}

	if err := writer.WriteFiles(files, opts); err != nil {
		t.Fatalf("unexpected error in dry-run: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Dockerfile") || !strings.Contains(output, files.Dockerfile) {
		t.Errorf("expected dry-run output to include Dockerfile, got:\n%s", output)
	}
	if !strings.Contains(output, ".dockerignore") || !strings.Contains(output, files.Dockerignore) {
		t.Errorf("expected dry-run output to include .dockerignore, got:\n%s", output)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "Dockerfile")); !os.IsNotExist(err) {
		t.Errorf("expected Dockerfile NOT to exist on disk during dry-run")
	}
	if _, err := os.Stat(filepath.Join(tempDir, ".dockerignore")); !os.IsNotExist(err) {
		t.Errorf("expected .dockerignore NOT to exist on disk during dry-run")
	}
}

func TestWriteFilesExistingFile(t *testing.T) {
	tempDir := t.TempDir()

	old := "OLD CONTENT\n"
	if err := os.WriteFile(filepath.Join(tempDir, "Dockerfile"), []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}

	files := generator.GeneratedFiles{
		Dockerfile:   "NEW CONTENT\n",
		Dockerignore: ".git\n",
	}

	// Without Force -> must refuse
	err := writer.WriteFiles(files, writer.WriteOptions{TargetDir: tempDir})
	if !errors.Is(err, writer.ErrFileExists) {
		t.Fatalf("expected ErrFileExists, got: %v", err)
	}

	// With Force -> must overwrite
	err = writer.WriteFiles(files, writer.WriteOptions{TargetDir: tempDir, Force: true})
	if err != nil {
		t.Fatalf("unexpected error with Force: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(tempDir, "Dockerfile"))
	if string(got) != files.Dockerfile {
		t.Errorf("expected overwrite to %q, got %q", files.Dockerfile, string(got))
	}
}

func TestWriteFilesBackup(t *testing.T) {
	tempDir := t.TempDir()
	dockerfilePath := filepath.Join(tempDir, "Dockerfile")

	old := "OLD CONTENT\n"
	if err := os.WriteFile(dockerfilePath, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}

	files := generator.GeneratedFiles{
		Dockerfile:   "NEW CONTENT\n",
		Dockerignore: ".git\n",
	}

	opts := writer.WriteOptions{TargetDir: tempDir, Force: true, Backup: true}
	if err := writer.WriteFiles(files, opts); err != nil {
		t.Fatalf("unexpected error with Backup: %v", err)
	}

	bak, err := os.ReadFile(dockerfilePath + ".bak")
	if err != nil {
		t.Fatalf("expected .bak to exist: %v", err)
	}
	if string(bak) != old {
		t.Errorf("expected .bak %q, got %q", old, string(bak))
	}

	got, _ := os.ReadFile(dockerfilePath)
	if string(got) != files.Dockerfile {
		t.Errorf("expected new %q, got %q", files.Dockerfile, string(got))
	}
}

func TestWriteFilesPrompt(t *testing.T) {
	files := generator.GeneratedFiles{Dockerfile: "NEW\n", Dockerignore: ".git\n"}

	t.Run("y proceeds", func(t *testing.T) {
		tempDir := t.TempDir()
		p := filepath.Join(tempDir, "Dockerfile")
		if err := os.WriteFile(p, []byte("OLD\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		opts := writer.WriteOptions{
			TargetDir: tempDir,
			In:        strings.NewReader("y\n"),
			Out:       &out,
		}
		if err := writer.WriteFiles(files, opts); err != nil {
			t.Fatalf("expected proceed on y, got: %v", err)
		}
		got, _ := os.ReadFile(p)
		if string(got) != "NEW\n" {
			t.Errorf("expected overwrite, got %q", string(got))
		}
	})

	t.Run("n aborts", func(t *testing.T) {
		tempDir := t.TempDir()
		p := filepath.Join(tempDir, "Dockerfile")
		if err := os.WriteFile(p, []byte("OLD\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		opts := writer.WriteOptions{
			TargetDir: tempDir,
			In:        strings.NewReader("n\n"),
			Out:       &bytes.Buffer{},
		}
		err := writer.WriteFiles(files, opts)
		if !errors.Is(err, writer.ErrAborted) {
			t.Fatalf("expected ErrAborted, got: %v", err)
		}
		got, _ := os.ReadFile(p)
		if string(got) != "OLD\n" {
			t.Errorf("expected file untouched, got %q", string(got))
		}
	})
}
