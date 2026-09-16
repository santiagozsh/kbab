# Scope v1 to Backend Server Workload Archetype

In v1, kbab exclusively generates container setups for long-running backend
HTTP servers (APIs), deferring other workload archetypes (CLIs, background
workers, cron jobs) to future releases.

While kbab detects workspace metadata heuristically, the CLI philosophy
mandates user confirmation over silent assumptions. Focusing v1 on backend
servers keeps multi-stage templates, hot-reloading toolchains (Air, Nodemon,
Reloaders), port exposure, and security configurations coherent and
production-grade before expanding the archetype matrix.
