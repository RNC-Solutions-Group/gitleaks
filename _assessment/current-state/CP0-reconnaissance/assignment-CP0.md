# CP0 — Repository and environment reconnaissance

**Status:** Ready; not executed  
**Authority:** Execute only when explicitly instructed `Execute CP0`  
**Scope:** Read-only repository/tooling inspection and safe diagnostics; documentation only in `_assessment/current-state/CP0-reconnaissance/`.

## Objective

Establish the exact code baseline, repository layout, toolchain and dependency inventory, test infrastructure, and safe execution methodology that CP1–CP5 will rely on. Do not attempt the full architecture analysis or change the program.

## Required investigations

1. **Identity and lineage:** Confirm repository name, current branch, clean/dirty status, commit SHA, default branch, available remotes, and any meaningful fork/upstream relationship. Do not fetch or switch branches if that would violate an active safety instruction; record unavailable lineage instead.
2. **Structure and entry points:** Inventory major Go packages, CLI entry points, generated/embedded assets, test directories, scripts, Makefile, CI workflows, and existing docs. Identify probable subsystems without presenting speculative behavior as fact.
3. **Toolchain and dependencies:** Record Go version from `go.mod`, installed Go version, OS/architecture, module state, relevant scanners/profilers, and tool availability. Note build tags and platform-specific behavior where apparent.
4. **Safe command matrix:** Inspect `Makefile` and scripts for side effects. Establish which exact direct Go test/vet/build/benchmark commands can be run without formatting, dependency changes, code generation, or writing committed artifacts.
5. **Baseline diagnostics:** Execute reasonable non-mutating smoke and Go tests (for example `go test -mod=readonly ./...`) and lightweight static checks if available. Record duration, exit code, failures, and prerequisites. Do **not** run `make test` or `make build` because they mutate source/dependencies.
6. **Fixture repository:** Check whether `../gitleaks-test-repo` exists, its nature and safe read-only use, and whether it contains documented cases usable in later capability/performance work. Do not expose secret-like fixture contents.
7. **Scope boundary:** Note any existing files changed by diagnostics and whether the branch source remains unchanged. Identify blockers or missing prerequisites for CP1–CP5.

## Big questions

- Which exact source commit should future findings cite, independently of subsequent documentation commits?
- What can be built, tested, and profiled reproducibly in the available environment?
- Which existing commands have side effects or hidden prerequisites?
- What is the smallest accurate component/entry-point map needed before CP1?
- Which unknowns affect the credibility of later performance and security findings?

## Deliverables

Create `_assessment/current-state/CP0-reconnaissance/report-CP0.md` using the shared report template, containing:

- A repository/branch/SHA/environment manifest.
- A package and entry-point inventory with source references.
- Dependencies and tool-availability summary.
- A **safe command matrix** indicating allowed, excluded, and unavailable operations with reasons.
- Actual test/diagnostic results, including failures or skipped checks.
- The sibling fixture-repository readiness assessment.
- A limitations/risks register and CP1 readiness recommendation.
- Before/after Git working-tree status and an explicit statement of protected-file modifications (if any).

Add smaller evidence documents inside this CP0 directory only if needed for clarity; link them from the report. Update `_assessment/README.md` when the report exists.

## Hard stops and continuation rules

**Stop the affected action immediately** if repository/branch identity differs from the assignment, a command would modify protected files, credentials or raw secret-like values would enter a report, or an unexpected source mutation occurs. Report what happened; never silently reset/delete changes.

An unavailable optional tool, missing sibling fixture repository, or failing existing test is **not automatically a reason to abandon CP0**. Record it as an unavailable check, failed baseline, or limitation and continue other safe investigations. Do not describe an unrun test as passing.

## Definition of done

Meet the shared [`../../README.md`](../../README.md#shared-definition-of-done) checklist, fully address the deliverables above, and explicitly identify each remaining unknown. At completion, summarize results and **stop for Dave's review**. CP1 must not start automatically.
