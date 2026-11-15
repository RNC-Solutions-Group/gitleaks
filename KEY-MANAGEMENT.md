# Grouping Key Management (OP-01)

Grouping is opt-in. The HMAC key is an operational secret. Its length is not proof of entropy: a 32-byte file containing predictable text is still weak. Generate keys with a cryptographically secure random source; do not rely on `loadGroupingKey`'s length check as entropy validation.

## Generate

POSIX:

```bash
umask 077
openssl rand -out "$HOME/.config/gitleaks/grouping.key" 32
chmod 600 "$HOME/.config/gitleaks/grouping.key"
```

Alternative with a cryptographically secure OS source:

```bash
umask 077
head -c 32 /dev/urandom > "$HOME/.config/gitleaks/grouping.key"
chmod 600 "$HOME/.config/gitleaks/grouping.key"
```

Windows PowerShell (DPAPI/ACL guidance, not a POSIX-mode claim):

```powershell
$key = [byte[]]::new(32)
$rng = [Security.Cryptography.RandomNumberGenerator]::Create()
$rng.GetBytes($key)
[IO.File]::WriteAllBytes("$env:USERPROFILE\grouping.key", $key)
icacls "$env:USERPROFILE\grouping.key" /inheritance:r /grant:r "$env:USERNAME:(R,W)"
```

Use a protected secret store/CI secret mechanism where available. Do not paste the key into command arguments, logs, issue reports, PRs, shell history or committed files.

## Use

```bash
gitleaks dir /path/to/target \
  --config /trusted/config.toml \
  --no-local-config \
  --grouped-report /protected/reports/grouped.json \
  --grouping-key-file /protected/keys/grouping.key
```

The key path and grouped destination are validated as a set with all other report destinations. The key must be outside the scan target.

## Controls and limitations

- Minimum length is 32 bytes, but entropy comes from generation, not length.
- POSIX: regular file, no symlink, not group/world-accessible (`0600` or stricter).
- Windows: regular, non-symlink file; POSIX mode bits are not authoritative, so ACLs and protected storage are the operator control.
- The loader opens the file and verifies the opened file identity matches the inspected path; replacement/symlink-swap between inspection and open is rejected.
- Hardlink aliases are not a provenance control. If a key is hardlinked into the scan target, an operator can still scan that alias; keep the key outside the tree and avoid link aliases.
- The key remains in memory for the scan/publication lifetime; Go does not guarantee zeroization of copied byte slices. Do not run grouping in a process exposed to untrusted memory dumps.
- Group IDs are deterministic for the same key, domain and secret. Rotating the key changes IDs by design.
- Do not reuse this key for fingerprint v2. Fingerprint v2 will use its own domain and key epoch; the current domain is `rnc-mod/report-grouping/v1`.
- Grouped reports are not crash-atomic with primary/informative reports. Preflight collision checks prevent destructive alias combinations, but a failure after primary publication can leave a complete primary and missing/partial grouped sidecar. The command exits nonzero and reports the primary path.
- Only the exact configured grouped destination (or a filesystem alias of it) is excluded from source scanning. A source file merely named like `grouped_001.json` is scanned normally.
- Aggregate publication is not capped at 1,000 alternatives. Preflight validates naming families structurally, and the selected candidate is rechecked against protected inputs after exclusive reservation. A directory with very many occupied names costs one reservation attempt per occupied candidate; this is linear allocation work, not a compatibility ceiling.

## Rotation

1. Stop using the old key.
2. Generate a new key and place it at a protected path.
3. Re-run grouping; IDs change.
4. Do not try to combine old/new IDs as if they were stable. Retain old reports keyed by their original key epoch if audit comparison is required.
5. Never commit either key.
