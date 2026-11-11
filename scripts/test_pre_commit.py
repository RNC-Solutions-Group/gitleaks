#!/usr/bin/env python3
"""Regression tests for the documented POSIX pre-commit helper."""
import os
import pathlib
import subprocess
import sys
import tempfile
import unittest


HELPER = pathlib.Path(__file__).with_name("pre-commit.py")


@unittest.skipIf(os.name == "nt", "POSIX executable shims are used for this test")
class PreCommitHelperTest(unittest.TestCase):
    def setUp(self):
        self.tempdir = tempfile.TemporaryDirectory()
        self.bin = pathlib.Path(self.tempdir.name) / "bin"
        self.bin.mkdir()
        (self.bin / "git").write_text("#!/bin/sh\nprintf 'true\\n'\n", encoding="utf-8")
        (self.bin / "git").chmod(0o755)

    def tearDown(self):
        self.tempdir.cleanup()

    def write_gitleaks(self, exit_code):
        path = self.bin / "gitleaks"
        path.write_text("#!/bin/sh\nexit %d\n" % exit_code, encoding="utf-8")
        path.chmod(0o755)

    def run_helper(self, extra_path=None):
        env = dict(os.environ)
        env["PATH"] = str(extra_path or self.bin)
        return subprocess.run(
            [sys.executable, str(HELPER)],
            env=env,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            check=False,
        )

    def test_success(self):
        self.write_gitleaks(0)
        result = self.run_helper()
        self.assertEqual(0, result.returncode)

    def test_finding_is_failure(self):
        self.write_gitleaks(1)
        result = self.run_helper()
        self.assertEqual(1, result.returncode)
        self.assertIn("Warning: gitleaks has detected", result.stdout)

    def test_scanner_errors_are_failures(self):
        for code in (2, 7, 126):
            with self.subTest(exit_code=code):
                self.write_gitleaks(code)
                result = self.run_helper()
                self.assertEqual(1, result.returncode)
                self.assertIn("exit code %d" % code, result.stderr)

    def test_missing_scanner_is_failure(self):
        result = self.run_helper(extra_path=self.bin)
        self.assertEqual(1, result.returncode)
        self.assertIn("failed to execute gitleaks", result.stderr)

    def test_non_executable_scanner_is_failure(self):
        path = self.bin / "gitleaks"
        path.write_text("#!/bin/sh\nexit 0\n", encoding="utf-8")
        path.chmod(0o644)
        result = self.run_helper()
        self.assertEqual(1, result.returncode)
        self.assertIn("failed to execute gitleaks", result.stderr)

    def test_disabled_hook_skips_scanner(self):
        self.write_gitleaks(99)
        (self.bin / "git").write_text(
            "#!/bin/sh\nprintf 'false\\n'\n", encoding="utf-8"
        )
        (self.bin / "git").chmod(0o755)
        result = self.run_helper()
        self.assertEqual(0, result.returncode)
        self.assertIn("precommit disabled", result.stdout)


if __name__ == "__main__":
    unittest.main()
