# CP2A — Review corrections and evidence hardening

**Status:** Review disposition complete — CP2 accepted with corrections  
**Assessment date:** 2026-11-09 (America/New_York)  
**Repository / branch:** `RNC-Solutions-Group/gitleaks` / `feat/issues-and-pr`  
**Assessed code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Documentation HEAD at review:** `7d36d8266d52e96eecc76a3aa6d3cf3c80b7f8ee`  
**Scope:** documentation/evidence reconciliation only; no application changes, no new defect remediation.  
**Parent report:** [`report-CP2.md`](report-CP2.md)

---

## 1. Owner-review disposition

CP2 is accepted as an executed assessment. Its measured commands, exit codes, findings, and limitations remain valid. The review corrects the **classification** of the 18 findings, adds reproducible detail for the highest-impact cases, and corrects three documentation inaccuracies:

1. **Native Go availability changed after CP0.** Go 1.27.3 for Windows is now installed at `C:\Program Files\Go` (directory created 2026-11-09 14:50; `go version go1.27.3 windows/amd64`). CP0's earlier “no native Go” observation was accurate for its 13:50–14:00 run; the toolchain was installed later, and CP2 did not re-check availability. CP2's runtime measurements remain the Linux `golang:1.27.3` container, not Windows. No Windows tests were performed.
2. **`go vet`/`lostcancel` is not declared a false positive.** The diagnostics remain a real static-analysis failure at the recorded lines. Source inspection shows `cancel()` in both branches, but that observation is not elevated to a proof that the analyzer is wrong; see §4.
3. **Race coverage and coverage totals are qualified.** A passing `-race` suite covers only exercised paths; 30.9% is the instrumented-package total, not whole-program coverage.

No new defect-remediation work was performed. CP2A performed one environment check (native Go) and re-used existing CP2 evidence; no new unbounded or network workload was run.

---

## 2. Finding classification crosswalk

Classification codes:

| Code | Meaning |
| --- | --- |
| **C1** | Confirmed functional/correctness defect |
| **C2** | Confirmed security or data-integrity risk |
| **C3** | Confirmed build/static-analysis failure |
| **C4** | Observed behavior requiring a documented contract decision |
| **C5** | Expected or intentional behavior |
| **C6** | Unverified hypothesis or incomplete reproduction |
| **C7** | Positive verification result |

| ID | Class | Corrected characterization | Basis |
| --- | --- | --- | --- |
| CP2-001 | **C4** | `--report-path -` without explicit `--report-format` is a usability/contract question; no repository documentation was found requiring an implicit format for `-` | Observed fatal `Unknown report format:`; code infers only from extension ([`cmd/root.go:366-379`](../../../cmd/root.go#L366-L379)) |
| CP2-002 | **C4** | `--report-template` help/code mismatch is a contract/documentation inconsistency | Help says “implies”; code fatals unless format is template ([`cmd/root.go:79`](../../../cmd/root.go#L79), [`:399-402`](../../../cmd/root.go#L399-L402)) |
| CP2-003 | **C2** | Existing report target is deleted during validation, before the scan/report completes | Sentinel reproduction; [`cmd/root.go:347-358`](../../../cmd/root.go#L347-L358) |
| CP2-004 | **C1** | Path-only finding corrupted by 100% redaction | Reproduced repeated `REDACTED`; [`report/finding.go:78-86`](../../../report/finding.go#L78-L86), [`detect/utils.go:158-160`](../../../detect/utils.go#L158-L160) |
| CP2-005 | **C2** | Git subprocess failure can be reported as a successful scan with no findings | Reproduced exit 0 after Git fatal; fail-open coverage gap |
| CP2-006 | **C4** | Unreadable file is skipped with a warning and the scan succeeds; whether this is an acceptable partial-coverage contract is undocumented | Non-root reproduction; [`sources/files.go:155-162`](../../../sources/files.go#L155-L162) |
| CP2-007 | **C2** | Baseline JSON inside the source is scanned and its stored secrets can be re-reported | Reproduced 2 findings from baseline JSON; [`detect/baseline.go:52-79`](../../../detect/baseline.go#L52-L79), [`detect/detect.go:294-299`](../../../detect/detect.go#L294-L299) |
| CP2-008 | **C1** | Invalid regex panics rather than returning a controlled configuration error | Reproduced exit 2 panic; [`config/config.go:123-128`](../../../config/config.go#L123-L128) |
| CP2-009 | **C2** | `--redact` does not redact debug/trace logs; synthetic secret appears in stderr | Reproduced in debug and trace modes; CP4 owns log-security disposition |
| CP2-010 | **C3** | `go vet` fails on two test-file `lostcancel` diagnostics | Exact output recorded; not proven false positive |
| CP2-011 | **C3 + C7** | `gore2regex` **application binary builds and loads default rules** (C7), but the full tagged **package graph** fails in generator rule files (C3) | Build/test logs; [`airtable.go:30`](../../../cmd/generate/config/rules/airtable.go#L30), [`perplexity.go:16`](../../../cmd/generate/config/rules/perplexity.go#L16) |
| CP2-012 | **C4** | Finding order is nondeterministic; whether report order is contractual is undocumented | 8-run reproduction; map traversal at [`detect/detect.go:327`](../../../detect/detect.go#L327) |
| CP2-013 | **C4** | Timeout is not a hard wall for blocked stdin; deadline is observed only after the read returns | Reproduced 2,009 ms for a 1 s timeout; contract for blocking I/O undocumented |
| CP2-014 | **C1** | Git commit allowlist logs a skip but the diff is still scanned; no demonstrated false negative because the detector-level allowlist suppresses findings | Trace reproduction; [`sources/git.go:341-346`](../../../sources/git.go#L341-L346) |
| CP2-015 | **C4** | Empty CSV report is a zero-byte file (no header); output contract is undocumented | Reproduced; [`report/csv.go:15-18`](../../../report/csv.go#L15-L18) |
| CP2-016 | **C5** | `git --pre-commit` scans the unstaged diff while `--staged` scans staged changes; observed difference is intentional mode semantics, not a defect | Reproduced 0 vs 1; source [`cmd/git.go:64-77`](../../../cmd/git.go#L64-L77) |
| CP2-017 | **C6** | Truncated-archive behavior remains unverified; no diagnostic was observed, but whether it was skipped as binary or failed silently is unresolved | Reproduction produced no error/finding; no claim made |
| CP2-018 | **C7** | Positive verification: build, 275/275 tests, race, shuffle, gofmt, CLI capabilities, reports, exit codes | See report §2–§3 |

**Negative findings preserved:** CP2-014 does not establish a false negative; CP2-017 remains unverified; CP1’s map-order hypothesis was confirmed as nondeterminism, not as an incorrect finding set.

---

## 3. Reproduction evidence for priority findings

The CP2 run created these under the disposable workspace `/work/cases` and external output `/out/cases`; raw artifacts are outside the tracked repository and may not be retained. The commands below are the exact forms used (synthetic tokens are shown as `CP2TOKEN-NNNN`). Environment for all: Linux container `golang:1.27.3`, `CGO_ENABLED=1`, `--network none`, external `/gomod` and `/gocache`, disposable source `/work/src`.

### CP2-003 — report-target deletion

- **Fixture:** empty directory `/work/cases/emptydir`; sentinel file `/out/cases/target-delete.json` containing `SENTINEL`; missing template `/out/cases/does-not-exist.tmpl`.
- **Command:** `gitleaks dir --config /work/cases/token.toml --no-banner --report-format template --report-template /out/cases/does-not-exist.tmpl --report-path /out/cases/target-delete.json /work/cases/emptydir`
- **Expected:** failure on missing template without deleting an unrelated existing target.
- **Actual:** exit 1; stderr `FTL Invalid report template … no such file or directory`; sentinel file absent afterward.
- **Source anchors:** [`cmd/root.go:347-358`](../../../cmd/root.go#L347-L358) (create/remove validation), [`:391-394`](../../../cmd/root.go#L391-L394) (template load).
- **Confidence:** confirmed; evidence ID `report-delete`.
- **Limitation:** symlink-target behavior was not exercised; CP4 may assess.

### CP2-004 — path-only redaction corruption

- **Fixture:** `/work/cases/pathdir/cp2-secret.txt`; config `/work/cases/pathonly.toml` with a path-only rule (`path = '''cp2-secret\\.txt$'''`).
- **Command:** `gitleaks dir --config /work/cases/pathonly.toml --no-banner --redact --report-format json --report-path - --exit-code 0 /work/cases/pathdir`
- **Expected:** `Match: "file detected: …"` and a redacted/empty secret, preserving readable match text.
- **Actual:** exit 0; JSON `Match` = `REDACTED` inserted between every character; `Secret` = `REDACTED`.
- **Source anchors:** path-only finding has no `Secret` ([`detect/detect.go:397-405`](../../../detect/detect.go#L397-L405)); `Redact` uses `ReplaceAll` with the empty secret ([`report/finding.go:78-86`](../../../report/finding.go#L78-L86)); applied by `filter` ([`detect/utils.go:158-160`](../../../detect/utils.go#L158-L160)).
- **Confidence:** confirmed; evidence ID `pathonly-redact`.
- **Limitation:** output was captured as JSON; verbose rendering follows the same corrupted `Match` but was not separately captured.

### CP2-005 — Git failure reported as success

- **Fixture:** `/work/cases/gitrepo`, synthetic commit containing a synthetic token; `/work/cases/token.toml`.
- **Command:** `gitleaks git --config /work/cases/token.toml --log-opts "--not-a-real-option" --no-banner --report-format json --report-path - --exit-code 0 /work/cases/gitrepo`
- **Expected:** failed/partial scan with nonzero status or explicit incomplete-scan signal.
- **Actual:** stderr `[git] fatal: unrecognized argument…` and `error="stderr is not empty"`; `0 commits scanned`; stdout `[]`; **exit 0**; `no leaks found`.
- **Source anchors:** stderr listener ([`sources/git.go:226-271`](../../../sources/git.go#L226-L271)), early return on error ([`:414-420`](../../../sources/git.go#L414-L420)), callback swallows error ([`detect/detect.go:228-232`](../../../detect/detect.go#L228-L232)).
- **Confidence:** confirmed; evidence ID `git-badflag`.
- **Limitation:** malformed *diff output* with a zero-length stderr was not reproduced; distinct from #1338/#2246.

### CP2-006 — unreadable-file partial success

- **Fixture:** `/work/cases/permcase/denied.txt` (mode `000`, contains a synthetic token) and `ok.txt` (readable, contains a different synthetic token).
- **Command:** `docker run --user 1000:1000 … gitleaks dir --config /work/cases/token.toml --no-banner --report-format json --report-path - --exit-code 0 /work/cases/permcase`
- **Expected:** explicit partial-coverage status, or documented warning-only behavior.
- **Actual:** exit 0; stderr warning `skipping file: permission denied` for `denied.txt`; one finding from `ok.txt`; no partial status in the result.
- **Source anchors:** [`sources/files.go:155-162`](../../../sources/files.go#L155-L162) (open failure returns `nil`), [`cmd/directory.go:69-73`](../../../cmd/directory.go#L69-L73) (error only logged).
- **Confidence:** confirmed observed behavior; classification as contract question C4, not automatically a defect.
- **Limitation:** root runs mask permission failures; the non-root precondition was verified by `chmod 000` and `--user 1000`.

### CP2-007 — baseline self-scan

- **Fixture:** `/work/cases/baseline-case/token.txt`; first scan writes `/work/cases/baseline-case/baseline.json` inside the source.
- **Commands:**
  1. `gitleaks dir --config …/token.toml --report-path /work/cases/baseline-case/baseline.json --exit-code 0 /work/cases/baseline-case`
  2. `gitleaks dir --config …/token.toml --baseline-path /work/cases/baseline-case/baseline.json --report-path - --exit-code 0 /work/cases/baseline-case`
- **Expected:** baseline suppress original finding and baseline file excluded from scanning.
- **Actual:** second run exit 0 with 2 findings, both in `baseline.json` (JSON lines 9 and 10); original `token.txt` finding suppressed.
- **Source anchors:** baseline relative path ([`detect/baseline.go:52-79`](../../../detect/baseline.go#L52-L79)); strict path skip ([`detect/detect.go:294-299`](../../../detect/detect.go#L294-L299)).
- **Confidence:** confirmed; evidence ID `baseline-second`.
- **Limitation:** source path was absolute; other path forms were not exhaustively tested.

### CP2-008 — invalid-regex panic

- **Fixture:** `/work/cases/badregex.toml` with `regex = '''['''`.
- **Command:** `gitleaks dir --config /work/cases/badregex.toml --no-banner --report-format json --report-path - --exit-code 0 /work/cases/emptydir`
- **Expected:** actionable config error, controlled exit.
- **Actual:** panic `regexp: Compile(...): missing closing ]`; stack includes `config.(*ViperConfig).Translate` and `cmd.Config`; **exit 2**.
- **Source anchors:** [`config/config.go:123-128`](../../../config/config.go#L123-L128) (`MustCompile`), [`cmd/root.go:240`](../../../cmd/root.go#L240).
- **Confidence:** confirmed; evidence ID `bad-regex`.
- **Limitation:** panic text contains the malformed pattern but no real secret.

### CP2-009 — redaction leakage in debug/trace

- **Fixtures:** `/work/cases/encodeddir` (base64 synthetic token); `/work/cases/inline-dir/allowed.txt` (token with `gitleaks:allow`).
- **Commands:**
  - `gitleaks dir --config …/token.toml --log-level debug --redact --max-decode-depth 2 --report-path - --report-format json … /work/cases/encodeddir`
  - `gitleaks dir --config …/token.toml --log-level trace --redact --report-path - --report-format json … /work/cases/inline-dir`
- **Expected:** `--redact` prevents raw secret in findings and logs.
- **Actual:** exit 0; synthetic token present in stderr (decoder debug segment value; `skipping finding: 'gitleaks:allow'` trace). Finding output itself was redacted/empty as applicable.
- **Source anchors:** decoder debug value ([`detect/codec/decoder.go:93-101`](../../../detect/codec/decoder.go#L93-L101)); inline-allow trace ([`detect/detect.go:510-515`](../../../detect/detect.go#L510-L515)); redaction only later in `filter`.
- **Confidence:** confirmed; evidence IDs `redact-debug-leak`, `redact-trace-inline`; CP4 disposition.
- **Limitation:** logging levels are operator-chosen; default info does not emit these messages.

### CP2-013 — delayed timeout enforcement

- **Fixture:** bounded delayed stdin producer `(sleep 2; printf 'CP2TOKEN-NNNN\n')`.
- **Command:** `… gitleaks stdin --config …/token.toml --timeout 1 --report-path - --report-format json --exit-code 0`, with producer piped in.
- **Expected:** abort near the 1 s deadline.
- **Actual:** ~2,009 ms elapsed; after the read returned, stderr `FTL failed scan input from stdin` / `context deadline exceeded`; exit 1; no findings.
- **Source anchors:** timeout context ([`cmd/root.go:48-58`](../../../cmd/root.go#L48-L58)); blocking read path ([`sources/file.go:174-190`](../../../sources/file.go#L174-L190)).
- **Confidence:** confirmed timing; contract classification C4 (blocking I/O behavior not documented).
- **Limitation:** `protect` end-to-end timeout was not run; the Git subprocess still receives the command context.

---

## 4. Documentation corrections

1. **Native Go.** `C:\Program Files\Go\bin\go.exe` now reports `go1.27.3 windows/amd64`; the install directory was created 2026-11-09 14:50. CP0's earlier absence (13:50–14:00) was accurate at the time. CP2 ran at/after 15:18 but did not re-check native Go and used the Linux container for parity. No conclusion is drawn about Windows application behavior.
2. **Linux baseline preserved.** All CP2/CP3 performance results are Linux/amd64 under `go1.27.3`; native Windows was not used for measurement.
3. **`go vet`/lostcancel.** The two diagnostics at [`detect/detect_test.go:2083`](../../../detect/detect_test.go#L2083) and [`:2122`](../../../detect/detect_test.go#L2122) remain a confirmed static-analysis failure. Source shows `cancel()` in both branches, but this is not sufficient to declare the analyzer wrong; a `defer cancel()` rewrite or a documented suppression would be needed. CI does not run `go vet`.
4. **Race scope.** `go test -race ./...` passed, but that establishes race-freedom only for exercised paths; it is not a proof about unexercised code.
5. **Coverage scope.** 30.9% is the total over instrumented/tested packages; `cmd`/`main` were not instrumented without `-coverpkg`. Whole-program coverage is lower.
6. **Status.** `report-CP2.md`, `_assessment/README.md`, and the roadmap were corrected/updated as part of CP2A/CP3.

---

## 5. Remaining uncertainties

- CP2-001/002/012/013/015 need **contract decisions**, not code changes yet.
- CP2-006's warning-only partial scan may be intentional; documented policy is needed.
- CP2-014's redundant scan has no demonstrated false negative; its cost/benefit is a CP3/CP4 question.
- CP2-017 remains unverified (archive skip vs silent failure).
- `go vet` lostcancel disposition remains unresolved.
- Windows native behavior remains unexecuted despite Go now being installed.
- Performance/security consequences are CP3/CP4 scope.

---

## 6. Handoff

- **CP3:** use CP2's classification, fixture definitions, and resource controls; do not treat contract questions as confirmed defects; keep #2019 unbounded and #1572 device scenarios excluded.
- **CP4:** debug/trace redaction leakage (CP2-009), baseline self-scan (CP2-007), report-target deletion/symlink risk (CP2-003), `.gitleaksignore` untrusted input (CP1 AR-11), and `gore2regex` tagged graph (CP2-011).
- **CP2 report:** revised with the CP2A classification labels and link to this review; measured results preserved.
