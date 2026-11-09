# Assessment branch instructions

These instructions govern work on `feat/issues-and-pr` in `RNC-Solutions-Group/gitleaks`. This is an **assessment and technical thought-development branch**, not an application implementation branch. Verify the active repository and branch at the start of every session. If either differs, stop and ask for direction before creating or changing files.

## Authority and write boundaries

The only in-repository documentation workspaces authorized on this branch are:

| Path | Ownership | Agent permission |
| --- | --- | --- |
| `_backlog/` | Dave and GPT maintain the upstream issue/PR evidence and dispositions | Read-only unless Dave explicitly assigns a backlog edit |
| `_assessment/` | Technical assessment, planning, comparisons, reports, and reproducible evidence | May create or revise assessment documents within the active authorized checkpoint |

All application source, tests, fixtures, scripts, build files, configuration, dependencies, CI, and other repository paths are **read-only**. Never edit, format, regenerate, delete, cherry-pick, merge, or apply patches to them as part of an assessment. Do not modify GitHub issues, PRs, comments, or other remote repositories. Do not merge, rebase, push, change branches, or create implementation branches without a separate explicit instruction.

Temporary build/test outputs and caches outside the repository may be used for authorized diagnostics. Do not commit test artifacts, raw scanner findings, caches, or secret-like fixture values. If an intended action would write outside `_assessment/` within this repository, stop that action and seek direction.

## Read and execute safely

- Read any repository files, Git history, external issue/PR links, and linked code as needed for the authorized checkpoint.
- Execute the existing source and non-destructive Go diagnostics. Prefer direct commands such as `go test -mod=readonly ./...`, `go test -mod=readonly -race ./...`, `go vet -mod=readonly ./...`, and carefully scoped benchmarks. Place profiles, coverage files, and compiled binaries in an external temporary directory.
- The sibling test fixture repository is expected at `../gitleaks-test-repo` relative to this repository root. Verify it exists; use it read-only, and never print or commit raw secret-like fixtures or findings.
- Inspect command effects before invoking scripts or Makefile targets. **Do not run `make test`, `make build`, `make format`, `make profile`, `make clean`, `go fmt`, `go mod tidy`, or `go generate`**: existing targets may format source, update dependencies, generate files, or remove output.
- Record `git status --short`, the current commit SHA, test commands, environment, and exit codes. Check status again after diagnostics. If a command changes protected paths, stop, report the changes, and do not attempt an unapproved destructive restoration.
- Missing optional tools or an unavailable sibling fixture repository should be reported as limitations, not silently treated as successful tests. Continue with safe independent checks where possible.

## Cold-start protocol

1. Verify repository identity, branch, current commit SHA, and working-tree state.
2. Read this file, `_assessment/AGENTS.md`, `_assessment/README.md`, and `_assessment/current-state/assessment-roadmap.md`.
3. Summarize the branch purpose, allowed write paths, current checkpoint state, expected deliverables, and important restrictions.
4. **Do not begin an assessment checkpoint merely because an assignment file exists.** Wait for a direct instruction such as `Execute CP0`.

When authorized, read the specific checkpoint assignment and use its scope, evidence standard, and stop conditions. The assessment workspace README owns the directory map, templates, and common definition of done.
