# kbab

A developer-friendly CLI tool that inspects a workspace and generates opinionated, production-grade containerization setups (Dockerfiles and `.dockerignore`).

## Language

**Runtime**:
The language platform and execution environment of the inspected application (`Go`, `Node`, `Python`).
_Avoid_: Language, stack, tech.

**Workload Archetype**:
The operational shape and lifecycle of the application being containerized (e.g. `API` for long-running backend HTTP servers).
_Avoid_: Project type, kind, category, mode.

**Stage**:
A distinct, named build target within a multi-stage Dockerfile (`dev`, `prod`).
_Avoid_: Environment, profile, target container.

**Detection Rule**:
A heuristic pattern (marker file or configuration key) used to identify a runtime or package manager.
_Avoid_: Check, validator, inspector.
