package templates_test

import (
	"io/fs"
	"testing"

	"github.com/santiagozsh/kbab/internal/templates"
)

func TestFSContainsAllTemplates(t *testing.T) {
	expectedFiles := []string{
		"go/Dockerfile.tmpl",
		"go/dockerignore.tmpl",
		"node/Dockerfile.tmpl",
		"node/dockerignore.tmpl",
		"py/Dockerfile.tmpl",
		"py/dockerignore.tmpl",
	}

	for _, file := range expectedFiles {
		t.Run(file, func(t *testing.T) {
			_, err := fs.Stat(templates.FS, file)
			if err != nil {
				t.Fatalf("expected %s to be embedded in FS, got error: %v", file, err)
			}
		})
	}
}
