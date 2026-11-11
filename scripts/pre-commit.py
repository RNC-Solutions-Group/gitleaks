#!/usr/bin/env python3
"""Helper script to be used as a pre-commit hook."""
import subprocess
import sys


WARNING = """Warning: gitleaks has detected sensitive information in your changes.
To disable the gitleaks precommit hook run the following command:

    git config hooks.gitleaks false
"""


def gitleaksEnabled():
    """Determine if the pre-commit hook for gitleaks is enabled."""
    try:
        result = subprocess.run(
            ["git", "config", "--bool", "hooks.gitleaks"],
            capture_output=True,
            text=True,
            check=False,
        )
    except OSError as exc:
        print(
            "gitleaks precommit: could not read git configuration: %s" % exc,
            file=sys.stderr,
        )
        return True

    if result.returncode != 0:
        # An unset or unreadable setting keeps the hook enabled. The scanner
        # command below remains the authority on whether the commit is safe.
        return True
    return result.stdout.strip().lower() != "false"


def main():
    if not gitleaksEnabled():
        print(
            "gitleaks precommit disabled"
            "     (enable with `git config hooks.gitleaks true`)"
        )
        return 0

    try:
        result = subprocess.run(["gitleaks", "protect", "-v", "--staged"])
    except OSError as exc:
        print(
            "gitleaks precommit: failed to execute gitleaks: %s" % exc,
            file=sys.stderr,
        )
        return 1

    if result.returncode == 0:
        return 0

    if result.returncode == 1:
        print(WARNING)
    else:
        print(
            "gitleaks precommit: scanner failed with exit code %d"
            % result.returncode,
            file=sys.stderr,
        )
    return 1


if __name__ == "__main__":
    sys.exit(main())
