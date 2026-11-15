#! /usr/bin/env bash
# NAME
#     profile.sh - generate comparable Gitleaks CPU/memory/execution profiles
#
# USAGE
#     profile.sh <gitleaks-path> <scan-target>
#
# OPTIONS (environment)
#     DECODE_DEPTH   recursive decode depth (default: 5)
#     HF_WARMUP      hyperfine warmup runs (default: 3)
#     HF_RUNS        hyperfine measured runs (default: 10)
#
# DESCRIPTION
#     Runs the same scanner command for benchmark and profile acquisition under
#     one timestamped output directory. All paths are quoted. A finding does not
#     terminate the script because the scan uses --exit-code=0; operational
#     failures still fail closed. Artifacts contain profiling data and metadata,
#     never a findings report.
set -euo pipefail

usage() {
	echo "usage: $0 <gitleaks-path> <scan-target>" >&2
	exit 2
}

[[ $# -eq 2 ]] || usage
gitleaks_path=$1
scan_target=$2
[[ -x "$gitleaks_path" ]] || { echo "scanner is not executable: $gitleaks_path" >&2; exit 2; }
[[ -e "$scan_target" ]] || { echo "scan target does not exist: $scan_target" >&2; exit 2; }

decode_depth=${DECODE_DEPTH:-5}
hf_warmup=${HF_WARMUP:-3}
hf_runs=${HF_RUNS:-10}
base_profile_dir="profile/$(date -u +%Y%m%dT%H%M%SZ)"
mkdir -p "$base_profile_dir"

log() { echo >&2 "$@"; }

log '========================================================================'
log 'generating profile data'
log "scanner: $gitleaks_path"
log "target:  $scan_target"
log "decode depth: $decode_depth"
log '------------------------------------------------------------------------'

{
	echo "generated_at_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
	echo "scanner=$gitleaks_path"
	echo "target=$scan_target"
	echo "decode_depth=$decode_depth"
	echo "hyperfine_warmup=$hf_warmup"
	echo "hyperfine_runs=$hf_runs"
	echo "uname=$(uname -a 2>/dev/null || true)"
	echo "go_version=$(go version 2>/dev/null || true)"
	echo "scanner_version=$("$gitleaks_path" version 2>/dev/null | head -n 1 || true)"
	if git -C "$scan_target" rev-parse HEAD >/dev/null 2>&1; then
		echo "target_git_head=$(git -C "$scan_target" rev-parse HEAD)"
		echo "target_git_shallow=$(git -C "$scan_target" rev-parse --is-shallow-repository 2>/dev/null || true)"
	fi
} > "$base_profile_dir/metadata.txt"

for scan_mode in dir git; do
	profile_dir="$base_profile_dir/$scan_mode"
	mkdir -p "$profile_dir"
	scan_cmd=("$gitleaks_path" --exit-code=0 --max-decode-depth "$decode_depth" --no-banner --log-level error "$scan_mode" "$scan_target")
	log "- mode: $scan_mode"

	if command -v hyperfine >/dev/null 2>&1; then
		export_path="$profile_dir/benchmark.json"
		log '  benchmark: hyperfine'
		hyperfine --warmup "$hf_warmup" --runs "$hf_runs" --export-json "$export_path" "${scan_cmd[@]}" >/dev/null
		log "    path: $export_path"
		if command -v yq >/dev/null 2>&1; then
			log '    results:'
			yq -P -oy '.results[] | pick(["mean","stddev","median","user","system","min","max"])' "$export_path" | sed 's/^/      /g'
		else
			log "    view: ${PAGER:-less} $export_path"
		fi
	else
		log '  benchmark: hyperfine not installed; skipped'
	fi

	log '  profile:'
	for profile_mode in cpu mem trace; do
		"${scan_cmd[@]}" --diagnostics="$profile_mode" --diagnostics-dir="$profile_dir" >"$profile_dir/${profile_mode}.log" 2>&1
		profile_file=$(find "$profile_dir" -maxdepth 1 -type f -name "${profile_mode}*" ! -name '*.log' -print -quit)
		if [[ -z "$profile_file" ]]; then
			echo "missing ${profile_mode} profile artifact in $profile_dir" >&2
			exit 1
		fi
		log "    - mode: $profile_mode"
		log "      path: $profile_file"
		if [[ "$profile_mode" == "trace" ]]; then
			log "      view: go tool trace $profile_file"
		else
			log "      view: go tool pprof -http=localhost: $gitleaks_path $profile_file"
		fi
	done
done

log '------------------------------------------------------------------------'
log "results in: $base_profile_dir"
log '========================================================================'
