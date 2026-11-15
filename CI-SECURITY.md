# Trusted Configuration Contract for CI

This document defines the supported contract for scanning **untrusted repository content** in CI. It is based on operator decision OP-2026-11-14-04: preserve compatible source-local discovery for ordinary local use, but use an explicitly trusted configuration for hostile CI targets.

## Threat model

An untrusted repository can commit:

- `.gitleaks.toml` (automatically discovered at the target root);
- `.gitleaksignore` fingerprints;
- misleading paths or allowlist-looking files;
- configuration-shaped files outside the intended trust boundary.

A local configuration can replace the rule set or suppress findings. This is an intentional compatibility feature, not a security guarantee. CI must not let the target select the policy that is supposed to judge it.

## Required CI pattern

Use a configuration path and, when needed, ignore/baseline files that live **outside** the untrusted checkout. Select them explicitly and disable target-local discovery:

```bash
gitleaks dir "$CHECKOUT" \
  --no-local-config \
  --config "$TRUSTED_CONFIG" \
  --gitleaks-ignore-path "$TRUSTED_IGNORE" \
  --baseline-path "$TRUSTED_BASELINE" \
  --redact \
  --report-path "$REPORT_PATH" \
  --report-format sarif
```

Rules:

1. `--config` must point to an operator-controlled file or an environment-provided trusted content source, never a path inside the checkout.
2. `--no-local-config` is required when the checkout may contain `.gitleaks.toml`.
3. `--gitleaks-ignore-path` must point to an operator-controlled file when suppressions are desired. Its contents are trust-sensitive.
4. Baseline files are findings reports and may contain plaintext secrets. Store them in a protected backend, not in the untrusted repository or public artifacts.
5. Never place a trusted config, ignore file, baseline, or key in the checkout and then rely on `--no-local-config` alone to protect the path. Grouping keys are additionally rejected inside the scan target; see [KEY-MANAGEMENT.md](KEY-MANAGEMENT.md).
6. Use `--redact` and a private/access-controlled report destination for shared CI output.
7. The default compatible behavior remains: explicit config, `GITLEAKS_CONFIG`, `GITLEAKS_CONFIG_TOML`, then target-local `.gitleaks.toml`, then the embedded default. `--no-local-config` skips only the target-local step.

## Precedence probes

The regression test `TestTrustedCIConfigurationContract` uses a synthetic hostile repository whose local `.gitleaks.toml` suppresses a synthetic secret:

- default compatible scan: no finding;
- `--no-local-config --config <trusted>`: finding recovered;
- `GITLEAKS_CONFIG=<trusted>` with `--no-local-config`: finding recovered.

This test proves precedence, not that arbitrary local configurations are safe. No upstream GitHub issue/PR was closed by this contract.

## CI workflow scope

The existing source-local discovery default is unchanged. Do not activate the upstream Gitleaks Action scanner workflow or add a new production CI requirement without operator authorization. If a future RNC workflow adopts this contract, it must pin the scanner artifact, keep the trusted config in a protected path, and publish reports only to an access-controlled destination.
