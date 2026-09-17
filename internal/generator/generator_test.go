package generator_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/santiagozsh/kbab/internal/detector"
	"github.com/santiagozsh/kbab/internal/generator"
)

func TestGenerateGoProject(t *testing.T) {
	cfg := detector.ProjectConfig{
		Type:           detector.Go,
		RuntimeVersion: "1.23.0",
		PackageManager: detector.None,
	}

	opts := generator.Options{}

	files, err := generator.Generate(cfg, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(files.Dockerfile, "FROM golang:1.23.0-alpine AS base") {
		t.Errorf("expected Dockerfile to expose default port 8000, got \n%s", files.Dockerfile)
	}
	if !strings.Contains(files.Dockerfile, "EXPOSE 8080") {
		t.Errorf("expected Dockerfile to expose default port 8080, got:\n%s", files.Dockerfile)
	}
	if !strings.Contains(files.Dockerignore, "bin/") {
		t.Errorf("expected Dockerignore to ignore bin/, got:\n%s", files.Dockerignore)
	}
	if !strings.Contains(files.Dockerignore, ".git") {
		t.Errorf("expected Dockerignore to ignore .git, got:\n%s", files.Dockerignore)
	}
}

func TestGenerateNodeProject(t *testing.T) {
	cfg := detector.ProjectConfig{
		Type:           detector.Node,
		RuntimeVersion: "22",
		PackageManager: detector.Pnpm,
	}

	opts := generator.Options{}

	files, err := generator.Generate(cfg, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Node base image and default port 3000
	if !strings.Contains(files.Dockerfile, "FROM node:22-alpine AS base") {
		t.Errorf("expected Dockerfile to use node:22-alpine, got:\n%s", files.Dockerfile)
	}
	if !strings.Contains(files.Dockerfile, "EXPOSE 3000") {
		t.Errorf("expected Dockerfile to expose default port 3000, got:\n%s", files.Dockerfile)
	}

	// Verify pnpm-specific commands are rendered
	if !strings.Contains(files.Dockerfile, "RUN corepack enable pnpm") {
		t.Errorf("expected Dockerfile to enable pnpm via corepack, got:\n%s", files.Dockerfile)
	}
	if !strings.Contains(files.Dockerfile, "pnpm install --frozen-lockfile") {
		t.Errorf("expected Dockerfile to use pnpm install, got:\n%s", files.Dockerfile)
	}

	// Verify Node-specific .dockerignore
	if !strings.Contains(files.Dockerignore, "node_modules/") {
		t.Errorf("expected Dockerignore to exclude node_modules/, got:\n%s", files.Dockerignore)
	}
}

func TestGeneratePythonProject(t *testing.T) {
	cfg := detector.ProjectConfig{
		Type:           detector.Python,
		RuntimeVersion: "3.12",
		PackageManager: detector.Uv,
	}

	opts := generator.Options{
		Framework: "fastapi", // Default port for Python is 8000
	}

	files, err := generator.Generate(cfg, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 1. Verify Python base image and default port 8000
	if !strings.Contains(files.Dockerfile, "FROM python:3.12-slim AS base") {
		t.Errorf("expected Dockerfile to use python:3.12-slim, got:\n%s", files.Dockerfile)
	}
	if !strings.Contains(files.Dockerfile, "EXPOSE 8000") {
		t.Errorf("expected Dockerfile to expose default port 8000, got:\n%s", files.Dockerfile)
	}

	// 2. Verify uv tool and commands
	if !strings.Contains(files.Dockerfile, "ghcr.io/astral-sh/uv:latest") {
		t.Errorf("expected Dockerfile to copy uv binary, got:\n%s", files.Dockerfile)
	}
	if !strings.Contains(files.Dockerfile, "uv sync --frozen") {
		t.Errorf("expected Dockerfile to run uv sync, got:\n%s", files.Dockerfile)
	}

	// 3. Verify FastAPI with Uvicorn
	if !strings.Contains(files.Dockerfile, "uvicorn") {
		t.Errorf("expected Dockerfile to execute uvicorn, got:\n%s", files.Dockerfile)
	}

	// 4. Verify Python .dockerignore
	if !strings.Contains(files.Dockerignore, "__pycache__/") {
		t.Errorf("expected Dockerignore to exclude __pycache__/, got:\n%s", files.Dockerignore)
	}
	if !strings.Contains(files.Dockerignore, ".venv/") {
		t.Errorf("expected Dockerignore to exclude .venv/, got:\n%s", files.Dockerignore)
	}
}

func TestGenerateCustomOptions(t *testing.T) {
	t.Run("custom port overrides default port", func(t *testing.T) {
		cfg := detector.ProjectConfig{
			Type:           detector.Go,
			RuntimeVersion: "1.23.0",
			PackageManager: detector.None,
		}

		opts := generator.Options{
			Port: "9090", // Explicitly overrides the default 8080
		}

		files, err := generator.Generate(cfg, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(files.Dockerfile, "EXPOSE 9090") {
			t.Errorf("expected Dockerfile to expose custom port 9090, got:\n%s", files.Dockerfile)
		}
		if !strings.Contains(files.Dockerfile, "ENV PORT=9090") {
			t.Errorf("expected Dockerfile to set ENV PORT=9090, got:\n%s", files.Dockerfile)
		}
	})

	t.Run("python standard framework renders python main.py", func(t *testing.T) {
		cfg := detector.ProjectConfig{
			Type:           detector.Python,
			RuntimeVersion: "3.12",
			PackageManager: detector.Pip,
		}

		opts := generator.Options{
			Framework: "standard", // Non-fastapi
		}

		files, err := generator.Generate(cfg, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(files.Dockerfile, `CMD ["python", "main.py"]`) {
			t.Errorf("expected Dockerfile to render python main.py, got:\n%s", files.Dockerfile)
		}
	})
}

func TestGenerateErrors(t *testing.T) {
	t.Run("returns ErrUnsupportedRuntime for unknown project type", func(t *testing.T) {
		cfg := detector.ProjectConfig{
			Type:           detector.ProjectType("rust"),
			RuntimeVersion: "1.75.0",
		}

		_, err := generator.Generate(cfg, generator.Options{})
		if !errors.Is(err, generator.ErrUnsupportedRuntime) {
			t.Errorf("expected ErrUnsupportedRuntime, got: %v", err)
		}
	})

	t.Run("returns ErrMissingRuntimeVersion when version is empty", func(t *testing.T) {
		cfg := detector.ProjectConfig{
			Type:           detector.Go,
			RuntimeVersion: "", // Missing version
		}

		_, err := generator.Generate(cfg, generator.Options{})
		if !errors.Is(err, generator.ErrMissingRuntimeVersion) {
			t.Errorf("expected ErrMissingRuntimeVersion, got: %v", err)
		}
	})
}
