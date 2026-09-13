// Package detector provides workspace inspection and runtime metadata
// extraction for containerization targets.
package detector

import (
	"encoding/json"
	"errors"
	"fmt"
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
	// python packages
	Pip    PackageManager = "pip"
	Poetry PackageManager = "poetry"
	Uv     PackageManager = "uv"

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

var defaultRules = []rule{
	{filename: "go.mod", project: Go},
	{filename: "package.json", project: Node},
	{filename: "requirements.txt", project: Python},
	{filename: "pyproject.toml", project: Python},
}

func Detect(targetDir string) (ProjectConfig, error) {
	for _, r := range defaultRules {
		path := filepath.Join(targetDir, r.filename)
		exists, err := fileExists(path)
		if err != nil {
			return ProjectConfig{}, err
		}

		if exists {
			switch r.project {
			case Go:
				return parseGoProject(targetDir)
			case Node:
				return parseNodeProject(targetDir)
			case Python:
				return parsePythonProject(targetDir)
			}
		}
	}
	return ProjectConfig{}, ErrUnsupportedProject
}

type rule struct {
	filename string
	project  ProjectType
}

func parseGoProject(targetDir string) (ProjectConfig, error) {
	path := filepath.Join(targetDir, "go.mod")
	data, err := os.ReadFile(path)
	if err != nil {
		return ProjectConfig{}, fmt.Errorf("reading go.mod: %w", err)
	}
	content := string(data)

	for line := range strings.SplitSeq(content, "\n") {
		line = strings.TrimSpace(line)
		if version, ok := strings.CutPrefix(line, "go "); ok {
			return ProjectConfig{
				Type:           Go,
				RuntimeVersion: version,
				PackageManager: None,
			}, nil
		}

	}
	return ProjectConfig{}, fmt.Errorf("go version directive not found: %w", ErrInvalidConfig)
}

type packageJSON struct {
	PackageManager string `json:"packageManager"`
}

func parseNodeProject(targetDir string) (ProjectConfig, error) {
	path := filepath.Join(targetDir, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return ProjectConfig{}, fmt.Errorf("reading package.json: %w", err)
	}
	var pkg packageJSON
	err = json.Unmarshal(data, &pkg)
	if err != nil {
		return ProjectConfig{}, fmt.Errorf("parsing package.json: %w", err)
	}

	pm := Npm // default fallback if pm is not explicitly found

	switch {
	case strings.HasPrefix(pkg.PackageManager, "pnpm"):
		pm = Pnpm

	case strings.HasPrefix(pkg.PackageManager, "yarn"):
		pm = Yarn

	case strings.HasPrefix(pkg.PackageManager, "bun"):
		pm = Bun

	default:
		if exists, _ := fileExists(filepath.Join(targetDir, "pnpm-lock.yaml")); exists {
			pm = Pnpm
		} else if exists, _ := fileExists(filepath.Join(targetDir, "yarn.lock")); exists {
			pm = Yarn
		} else if exists, _ := fileExists(filepath.Join(targetDir, "bun.lockb")); exists {
			pm = Bun
		}
	}

	return ProjectConfig{
		Type:           Node,
		RuntimeVersion: "22",
		PackageManager: pm,
	}, nil
}

func parsePythonProject(targetDir string) (ProjectConfig, error) {
	pm := Pip

	if exists, _ := fileExists(filepath.Join(targetDir, "uv.lock")); exists {
		pm = Uv
	} else if exists, _ := fileExists(filepath.Join(targetDir, "poetry.lock")); exists {
		pm = Poetry
	}

	return ProjectConfig{
		Type:           Python,
		RuntimeVersion: "3.12",
		PackageManager: pm,
	}, nil
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
