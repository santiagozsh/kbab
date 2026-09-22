// Package writer manages filesystem output, backups, and dry-run rendering.
package writer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	if opts.DryRun {
		out := opts.Out

		if out == nil {
			out = os.Stdout
		}
		if _, err := fmt.Fprintln(out, "=== Dockerfile ==="); err != nil {
			return fmt.Errorf("dry-run output: %w", err)
		}
		if _, err := fmt.Fprintln(out, files.Dockerfile); err != nil {
			return fmt.Errorf("dry-run output: %w", err)
		}
		if _, err := fmt.Fprintln(out, "=== .dockerignore ==="); err != nil {
			return fmt.Errorf("dry-run output: %w", err)
		}
		if _, err := fmt.Fprintln(out, files.Dockerignore); err != nil {
			return fmt.Errorf("dry-run output: %w", err)
		}
		return nil
	}

	targetDir := opts.TargetDir

	if targetDir == "" {
		targetDir = "."
	}

	dockerfilePath := filepath.Join(targetDir, "Dockerfile")
	dockerignorePath := filepath.Join(targetDir, ".dockerignore")

	var conflicts []string
	for _, p := range []string{dockerfilePath, dockerignorePath} {
		exists, err := fileExists(p)
		if err != nil {
			return fmt.Errorf("checking %s: %w", p, err)
		}
		if exists {
			conflicts = append(conflicts, p)
		}
	}

	if len(conflicts) > 0 && !opts.Force && !opts.Backup {
		if opts.In == nil {
			return fmt.Errorf("%w: %s", ErrFileExists, conflicts[0])
		}
		out := opts.Out
		if out == nil {
			out = os.Stdout
		}
		if _, err := fmt.Fprintf(out, "Files already exist: %s. Overwrite? [y/N]: ", strings.Join(conflicts, ", ")); err != nil {
			return fmt.Errorf("prompt output: %w", err)
		}
		line, err := bufio.NewReader(opts.In).ReadString('\n')
		if err != nil && len(line) == 0 {
			return fmt.Errorf("reading answer: %w", err)
		}
		ans := strings.TrimSpace(strings.ToLower(line))

		if ans != "y" && ans != "yes" {
			return ErrAborted
		}
	}

	if opts.Backup {
		for _, p := range []string{dockerfilePath, dockerignorePath} {
			exists, err := fileExists(p)
			if err != nil {
				return fmt.Errorf("checking %s: %w", p, err)
			}
			if exists {
				if err := os.Rename(p, p+".bak"); err != nil {
					return fmt.Errorf("backup %s: %w", p, err)
				}
			}
		}
	}

	if err := os.WriteFile(dockerfilePath, []byte(files.Dockerfile), 0o644); err != nil {
		return fmt.Errorf("writing Dockerfile: %w", err)
	}
	if err := os.WriteFile(dockerignorePath, []byte(files.Dockerignore), 0o644); err != nil {
		return fmt.Errorf("writing .dockerignore: %w", err)
	}
	return nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
