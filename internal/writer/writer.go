package writer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/santiagozsh/kbab/internal/generator"
)

// WriteOptions configures the file writing behavior.
type WriteOptions struct {
	TargetDir string
	DryRun    bool
	Force     bool
	Backup    bool
	Out       io.Writer
	In        io.Reader
}

var (
	ErrFileExists = errors.New("file already exists")
	ErrAborted    = errors.New("operation aborted by user")
)

func WriteFiles(files generator.GeneratedFiles, opts WriteOptions) error {
	targetDir := opts.TargetDir

	if targetDir == "" {
		targetDir = "."
	}

	dockerfilePath := filepath.Join(targetDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(files.Dockerfile), 0o644); err != nil {
		return fmt.Errorf("writing Dockerfile: %w", err)
	}

	dockerignorePath := filepath.Join(targetDir, ".dockerignore")
	if err := os.WriteFile(dockerignorePath, []byte(files.Dockerignore), 0o644); err != nil {
		return fmt.Errorf("writing .dockerinogre: %w", err)
	}

	return nil
}
