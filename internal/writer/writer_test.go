package writer_test

import (
	"os"
	"path/filepath"
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
