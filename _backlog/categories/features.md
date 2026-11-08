# Features — Gitleaks candidate register

This is the **single active category** for feature and developer-experience proposals selected for RNC evaluation. Issues and PRs for the same work package belong together here, even when one has a different title-only classification in the original snapshot. The [closed ledger](closed.md) remains a separate exclusion register.

## Pending review

| Upstream PR | Associated Gitleaks issue | Proposed behavior | RNC status | Assessment / next action |
| --- | --- | --- | --- | --- |
| [#652](https://github.com/gitleaks/gitleaks/pull/652) — @weineran, head `9fcedcd486a4f2eab56584f66482d842511f6a3d` | **None verified** in the 2026-11-08 open-issue/PR crosswalk | Improve the error presented when installing with a Go version below the declared minimum by using Go build constraints. | **PENDING REVIEW — NOT IMPLEMENTED** | PR dates to 2021 and specifically targets **Gitleaks v7**, Go **1.16**, and `go get` installation behavior. Determine whether any equivalent failure survives in the current v8 code/build process before considering adaptation; do not cherry-pick legacy build tags without compatibility tests. |

## Evaluation gate

1. Check relevance to the current fork's module path, Go toolchain requirements, supported install commands, and release artifacts.
2. Reproduce with an appropriate old-toolchain test environment only if still supported and safe; record error quality separately from actual minimum-version enforcement.
3. If adopted, retain contributor attribution and record the RNC implementation branch, PR, commit, and regression evidence.
4. If superseded or obsolete, record the reason here and move the item to [closed](closed.md) only on an explicit RNC disposition.

The original snapshot entry remains struck through in the [open PR Markdown report](../02-open-pull-requests.md), indicating category assignment—not upstream closure. JSON and CSV exports remain untouched.
