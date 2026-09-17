// Package generator renders Dockerfile and dockerignore templates
package generator

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/santiagozsh/kbab/internal/detector"
	"github.com/santiagozsh/kbab/internal/templates"
)

// GeneratedFiles contains the rendered output for Dockerfile and .dockerignore.
type GeneratedFiles struct {
	Dockerfile   string
	Dockerignore string
}

// Options allows overriding default template behaviors (e.g. port, framework).
type Options struct {
	Port      string
	Framework string
}

var (
	ErrUnsupportedRuntime    = errors.New("unsupported runtime")
	ErrMissingRuntimeVersion = errors.New("missing runtime version")
)

// templateContext maps internal data to template variable names.
type templateContext struct {
	RuntimeVersion string
	PackageManager string
	Framework      string
	Port           string
}

func Generate(cfg detector.ProjectConfig, opts Options) (GeneratedFiles, error) {
	if cfg.RuntimeVersion == "" {
		return GeneratedFiles{}, ErrMissingRuntimeVersion
	}

	port := opts.Port

	framework := opts.Framework
	if framework == "" && cfg.Type == detector.Python {
		framework = "fastapi"
	}

	if port == "" {
		switch cfg.Type {
		case detector.Go:
			port = "8080"
		case detector.Node:
			port = "3000"
		case detector.Python:
			port = "8000"
		}
	}

	ctx := templateContext{
		RuntimeVersion: cfg.RuntimeVersion,
		PackageManager: string(cfg.PackageManager),
		Framework:      framework,
		Port:           port,
	}

	var dockerfileTmpl, dockerignoreTmpl string

	switch cfg.Type {
	case detector.Go:
		dockerfileTmpl = "go/Dockerfile.tmpl"
		dockerignoreTmpl = "go/dockerignore.tmpl"
	case detector.Node:
		dockerfileTmpl = "node/Dockerfile.tmpl"
		dockerignoreTmpl = "node/dockerignore.tmpl"
	case detector.Python:
		dockerfileTmpl = "py/Dockerfile.tmpl"
		dockerignoreTmpl = "py/dockerignore.tmpl"
	default:
		return GeneratedFiles{}, fmt.Errorf("%w: %s", ErrUnsupportedRuntime, cfg.Type)
	}

	dockerfile, err := renderTemplate(dockerfileTmpl, ctx)
	if err != nil {
		return GeneratedFiles{}, fmt.Errorf("rendering dockerfile: %w", err)
	}

	dockerignore, err := renderTemplate(dockerignoreTmpl, ctx)
	if err != nil {
		return GeneratedFiles{}, fmt.Errorf("rendering dockerignore: %w", err)
	}

	return GeneratedFiles{
		Dockerfile:   dockerfile,
		Dockerignore: dockerignore,
	}, nil
}

func renderTemplate(tmplPath string, ctx templateContext) (string, error) {
	tmpl, err := template.ParseFS(templates.FS, tmplPath)
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", tmplPath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("executing template %s: %w", tmplPath, err)
	}
	return buf.String(), nil
}
