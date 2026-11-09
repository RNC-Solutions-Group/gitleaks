# Assessment workspace instructions

The root [`AGENTS.md`](../AGENTS.md) defines authority and safety boundaries; it always applies here.

Before starting an assigned checkpoint:

1. Read [`README.md`](README.md) for the workspace layout, evidence format, and shared definition of done.
2. Read [`current-state/assessment-roadmap.md`](current-state/assessment-roadmap.md) for mission context and checkpoint dependencies.
3. Read the specifically authorized checkpoint's `assignment-CP*.md` in full. No checkpoint is automatically authorized by completing the previous one.
4. Use [`templates/checkpoint-assignment.md`](templates/checkpoint-assignment.md) for future checkpoint assignments and [`templates/assessment-report.md`](templates/assessment-report.md) for investigation reports.

Write current-state checkpoint outputs inside their authorized `current-state/CPX-*/` folder. Update the workspace index or roadmap within `_assessment/` as needed. Keep `_backlog/` read-only; link to its canonical records rather than copying or reclassifying them.

Every material finding needs source/command evidence, a verification status, impact, and an explicit uncertainty statement when proof is incomplete. Treat the current-state baseline as a reference to challenge against code, not as an authority that overrides contrary evidence.

At the end of a checkpoint, report its deliverables, test results, blockers, and readiness for review, then **stop**. Do not advance, modify source, or conduct external PR integration experiments without separate authorization.
