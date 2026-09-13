// Package detector provides workspace inspection and runtime metadata
// extraction for containerization targets.
package detector

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type ProjectType string

const (
	Go      ProjectType = "go"
	Python  ProjectType = "py"
	Node    ProjectType = "node"
	Unknown ProjectType = "unknown"
)

type PackageManager string

const (
	Npm  PackageManager = "npm"
	Pnpm PackageManager = "pnpm"
	Yarn PackageManager = "yarn"
	Bun  PackageManager = "bun"
	None PackageManager = "none"
)

type ProjectConfig struct {
	Type           ProjectType
	RuntimeVersion string
	PackageManager PackageManager
}

var (
	ErrUnsupportedProject = errors.New("unsupported project type")
	ErrInvalidConfig      = errors.New("invalid project configuration")
)

func Detect(targetDir string) (ProjectConfig, error) {
	for _, r := range defaultRule {
		path := filepath.Join(targetDir, r.filename)
		exists, err := fileExists(path)
		if err != nil {
			return ProjectConfig{}, err
		}

		if exists {
			switch r.project {
			case Go:
				version, err := parseGoMod(path)
				if err != nil {
					return ProjectConfig{}, err
				}
				return ProjectConfig{
					Type:           Go,
					RuntimeVersion: version,
					PackageManager: None,
				}, nil
			}
		}
	}
	return ProjectConfig{}, ErrUnsupportedProject
}

type rule struct {
	filename string
	project  ProjectType
}

var defaultRule = []rule{
	{filename: "go.mod", project: Go},
	{filename: "package.json", project: Node},
	{filename: "requirements.txt", project: Python},
	{filename: "pyproject.toml", project: Python},
}

func parseGoMod(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	content := string(data)

	for line := range strings.SplitSeq(content, "\n") {
		line = strings.TrimSpace(line)
		if version, ok := strings.CutPrefix(line, "go "); ok {
			return version, nil
		}

	}
	return "", ErrInvalidConfig
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}
