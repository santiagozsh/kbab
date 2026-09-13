package detector_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/santiagozsh/kbab/internal/detector"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string]string // filename -> file content
		wantConfig  detector.ProjectConfig
		wantErr     error
	}{
		{
			name: "Go project with valid go.mod",
			files: map[string]string{
				"go.mod": "module github.com/example/app\n\ngo 1.23.0\n",
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Go,
				RuntimeVersion: "1.23.0",
				PackageManager: detector.None,
			},
			wantErr: nil,
		},
		{
			name: "Go project with missing go version directive",
			files: map[string]string{
				"go.mod": "module github.com/example/app\n",
			},
			wantConfig: detector.ProjectConfig{},
			wantErr:    detector.ErrInvalidConfig,
		},
		{
			name: "Node project with pnpm in package.json",
			files: map[string]string{
				"package.json": `{"name": "frontend", "packageManager": "pnpm@9.1.0"}`,
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Node,
				RuntimeVersion: "22",
				PackageManager: detector.Pnpm,
			},
			wantErr: nil,
		},
		{
			name: "Node project with yarn in package.json",
			files: map[string]string{
				"package.json": `{"name": "frontend", "packageManager": "yarn@4.0.0"}`,
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Node,
				RuntimeVersion: "22",
				PackageManager: detector.Yarn,
			},
			wantErr: nil,
		},
		{
			name: "Node project with bun in package.json",
			files: map[string]string{
				"package.json": `{"name": "frontend", "packageManager": "bun@1.1.0"}`,
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Node,
				RuntimeVersion: "22",
				PackageManager: detector.Bun,
			},
			wantErr: nil,
		},
		{
			name: "Node project with pnpm-lock.yaml fallback",
			files: map[string]string{
				"package.json":   `{"name": "frontend"}`,
				"pnpm-lock.yaml": "lockfileVersion: '9.0'",
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Node,
				RuntimeVersion: "22",
				PackageManager: detector.Pnpm,
			},
			wantErr: nil,
		},
		{
			name: "Node project with yarn.lock fallback",
			files: map[string]string{
				"package.json": `{"name": "frontend"}`,
				"yarn.lock":    "# yarn lockfile v1",
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Node,
				RuntimeVersion: "22",
				PackageManager: detector.Yarn,
			},
			wantErr: nil,
		},
		{
			name: "Node project with bun.lockb fallback",
			files: map[string]string{
				"package.json": `{"name": "frontend"}`,
				"bun.lockb":    "",
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Node,
				RuntimeVersion: "22",
				PackageManager: detector.Bun,
			},
			wantErr: nil,
		},
		{
			name: "Node project defaulting to npm",
			files: map[string]string{
				"package.json": `{"name": "frontend"}`,
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Node,
				RuntimeVersion: "22",
				PackageManager: detector.Npm,
			},
			wantErr: nil,
		},
		{
			name: "Python project with uv.lock",
			files: map[string]string{
				"pyproject.toml": "[project]\nname = 'backend'\n",
				"uv.lock":        "version = 1\n",
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Python,
				RuntimeVersion: "3.12",
				PackageManager: detector.Uv,
			},
			wantErr: nil,
		},
		{
			name: "Python project with poetry.lock",
			files: map[string]string{
				"pyproject.toml": "[tool.poetry]\nname = 'backend'\n",
				"poetry.lock":    "",
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Python,
				RuntimeVersion: "3.12",
				PackageManager: detector.Poetry,
			},
			wantErr: nil,
		},
		{
			name: "Python project with requirements.txt",
			files: map[string]string{
				"requirements.txt": "fastapi>=0.110.0\nuvicorn>=0.28.0\n",
			},
			wantConfig: detector.ProjectConfig{
				Type:           detector.Python,
				RuntimeVersion: "3.12",
				PackageManager: detector.Pip,
			},
			wantErr: nil,
		},
		{
			name:       "Empty directory returns ErrUnsupportedProject",
			files:      map[string]string{},
			wantConfig: detector.ProjectConfig{},
			wantErr:    detector.ErrUnsupportedProject,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			for name, content := range tt.files {
				filePath := filepath.Join(dir, name)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("failed creating test file %s: %v", name, err)
				}
			}

			cfg, err := detector.Detect(dir)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Detect() error = %v, wantErr = %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Detect() unexpected error = %v", err)
			}

			if cfg != tt.wantConfig {
				t.Errorf("Detect() = %+v, want %+v", cfg, tt.wantConfig)
			}
		})
	}
}
