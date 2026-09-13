# Specification: kbab

A developer-friendly CLI tool to generate production-grade, multi-stage Dockerfiles and `.dockerignore` files for Go, Node/TypeScript, and Python applications.

## 1. Problem Statement
Application developers frequently face challenges with containerization best practices:
- Multi-stage builds are complex and prone to misconfiguration.
- Development environments require hot-reloading without rebuilding images from scratch.
- Production environments require strict security standards: unprivileged users, minimal attack surfaces (distroless/scratch), and compiler cache optimization.
- Missing `.dockerignore` files leak sensitive files (`.env`, `.git`) and degrade build cache performance.

`kbab` resolves this by inspecting the target workspace and generating an opinionated, secure, and ready-to-run container setup.

## 2. Supported Runtimes & Stage Architecture

### Runtimes
- **Go**: Parses `go.mod` for Go version; configures `air` for dev hot-reloading; compiles with `CGO_ENABLED=0` and cache mounts; targets `gcr.io/distroless/static-debian12:nonroot` for production.
- **Node/TypeScript**: Parses `package.json` to detect package managers (`npm`, `pnpm`, `yarn`, `bun`); configures development hot-reloading; builds production assets and runs under an unprivileged `node` user in `node:alpine`.
- **Python**: Detects `pyproject.toml` and `requirements.txt`; configures development reload mode; optimizes pip cache mounts; runs under a non-root user in `python:slim`.

### Multi-stage Architecture
Each generated Dockerfile includes:
- **Target `dev`**: Mounts local workspace source code, enables hot-reloading, includes necessary dev tooling.
- **Target `prod`**: Lean production build, build-tool stripping, unprivileged user execution, compiler/package cache mounts.

## 3. CLI UX & Flags
- `kbab create [path]` (defaults to `.`):
  - `--dry-run`: Prints generated `Dockerfile` and `.dockerignore` to `stdout` without modifying the filesystem.
  - `--force`: Overwrites existing files without prompting.
  - `--yes`: Accepts recommended defaults non-interactively (suited for CI/CD).
  - `--port <port>`: Overrides the exposed application port.

## 4. Package Architecture (Go Idiomatic Boundaries)
- `cmd/kbab/`: Cobra CLI commands (`root`, `create`, `version`).
- `internal/detector/`: Inspects workspace, extracts versions/package managers, returns typed `ProjectConfig`.
- `internal/templates/`: Embedded templates via `embed.FS`.
- `internal/generator/`: Parses and renders `text/template` using `ProjectConfig`.
- `internal/writer/`: Manages filesystem I/O, backup files (`.bak`), interactive prompts, and dry-run execution.
