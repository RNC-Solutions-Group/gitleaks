#!/usr/bin/env python3
"""Exercise the published pre-commit hook with the real pre-commit runner.

This is an integration test, not a replacement for the Go unit tests. It is
skipped by its Go wrapper when pre-commit is not installed. Each scenario uses a
fresh temporary Git repository so staged-index state cannot leak between cases.
"""
import os
import pathlib
import shutil
import subprocess
import sys
import tempfile
import unittest


ROOT = pathlib.Path(__file__).resolve().parents[1]
HOOK_ID = "gitleaks-system"


def run(command, cwd, env=None):
    return subprocess.run(
        command,
        cwd=str(cwd),
        env=env,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        text=True,
        check=False,
    )


def git(repo, *args):
    result = run(["git", "-C", str(repo), *args], cwd=repo)
    if result.returncode != 0:
        raise AssertionError("git %s failed: %s" % (" ".join(args), result.stdout))
    return result


class PreCommitHookIntegrationTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.pre_commit = [sys.executable, "-m", "pre_commit"]
        probe = run(cls.pre_commit + ["--version"], cwd=ROOT)
        if probe.returncode != 0:
            raise unittest.SkipTest("pre-commit is not installed")
        cls.gitleaks = os.environ.get("GITLEAKS_BIN")
        if not cls.gitleaks or not pathlib.Path(cls.gitleaks).is_file():
            raise unittest.SkipTest("GITLEAKS_BIN is not set to a built binary")

        # try-repo resolves hooks from a committed revision. Build a tiny
        # throwaway hook repository containing the working-tree manifest so the
        # test exercises the file being edited, not an older HEAD.
        cls.hook_temp = tempfile.TemporaryDirectory()
        cls.hook_repo = pathlib.Path(cls.hook_temp.name)
        shutil.copyfile(ROOT / ".pre-commit-hooks.yaml", cls.hook_repo / ".pre-commit-hooks.yaml")
        git(cls.hook_repo, "init", "-q", "-b", "main")
        git(cls.hook_repo, "config", "user.email", "cpv9@example.invalid")
        git(cls.hook_repo, "config", "user.name", "CP-V9 Test")
        git(cls.hook_repo, "add", ".pre-commit-hooks.yaml")
        git(cls.hook_repo, "commit", "-q", "-m", "hook manifest")

    @classmethod
    def tearDownClass(cls):
        cls.hook_temp.cleanup()

    def setUp(self):
        self.tempdir = tempfile.TemporaryDirectory()
        self.repo = pathlib.Path(self.tempdir.name)
        git(self.repo, "init", "-q", "-b", "main")
        git(self.repo, "config", "user.email", "cpv9@example.invalid")
        git(self.repo, "config", "user.name", "CP-V9 Test")
        self.env = dict(os.environ)
        binary_dir = str(pathlib.Path(self.gitleaks).parent)
        self.env["PATH"] = binary_dir + os.pathsep + self.env.get("PATH", "")
        self.env["PRE_COMMIT_HOME"] = str(self.repo / ".pre-commit-cache")
        self.write_and_stage("README.md", "base\n", commit=True)

    def tearDown(self):
        self.tempdir.cleanup()

    def write_and_stage(self, name, content, commit=False):
        path = self.repo / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding="utf-8")
        git(self.repo, "add", name)
        if commit:
            git(self.repo, "commit", "-q", "-m", "base")

    def reset(self):
        git(self.repo, "reset", "-q", "--hard", "HEAD")
        git(self.repo, "clean", "-q", "-fdx")

    def try_repo_hook(self, hook_id, *args):
        return run(
            self.pre_commit
            + ["try-repo", str(self.hook_repo), hook_id]
            + list(args),
            cwd=self.repo,
            env=self.env,
        )

    def try_repo(self, *args):
        return self.try_repo_hook(HOOK_ID, *args)

    def assert_clean(self, result):
        self.assertEqual(0, result.returncode, result.stdout)

    def test_staged_scenarios(self):
        # No staged files: pre-commit may report a skip; the important property
        # is that it does not fail or forward a filename.
        self.assert_clean(self.try_repo())

        # One staged file.
        self.write_and_stage("one.txt", "ordinary text\n")
        self.assert_clean(self.try_repo())

        # Multiple staged files.
        self.write_and_stage("two.txt", "more ordinary text\n")
        self.write_and_stage("three.txt", "still ordinary\n")
        self.assert_clean(self.try_repo())

        # A staged filename containing spaces.
        self.write_and_stage("name with spaces.txt", "ordinary text\n")
        self.assert_clean(self.try_repo())

        # --all-files still runs the staged-index scanner with no forwarded
        # filenames; it does not turn this hook into an all-files scanner.
        self.reset()
        self.assert_clean(self.try_repo("--all-files"))

    def test_clean_change_and_secret(self):
        self.write_and_stage("clean.txt", "ordinary staged change\n")
        self.assert_clean(self.try_repo())

        self.reset()
        self.write_and_stage(
            "secret.txt",
            "-----BEGIN OPENSSH PRIVATE KEY-----\n"
            + ("A" * 80)
            + "\n-----END OPENSSH PRIVATE KEY-----\n",
        )
        result = self.try_repo()
        self.assertNotEqual(0, result.returncode, result.stdout)
        self.assertIn("leaks found", result.stdout.lower())

    def test_scanner_operational_failure(self):
        self.write_and_stage("clean.txt", "ordinary staged change\n")
        (self.repo / ".gitleaks.toml").write_text(
            "[[rules]]\nid = \"bad\"\nregex = '''[unterminated'''\n",
            encoding="utf-8",
        )
        result = self.try_repo()
        self.assertNotEqual(0, result.returncode, result.stdout)
        self.assertTrue(
            "config" in result.stdout.lower() or "regex" in result.stdout.lower(),
            result.stdout,
        )

    def test_explicit_files_hook(self):
        self.write_and_stage(
            "explicit-secret.txt",
            "-----BEGIN OPENSSH PRIVATE KEY-----\n"
            + ("B" * 80)
            + "\n-----END OPENSSH PRIVATE KEY-----\n",
        )
        result = self.try_repo_hook(
            "gitleaks-files-system", "--files", "explicit-secret.txt"
        )
        self.assertNotEqual(0, result.returncode, result.stdout)
        self.assertIn("leaks found", result.stdout.lower())

        # --all-files must scan tracked content through the explicit files
        # command rather than silently falling back to the staged-index rule.
        git(self.repo, "commit", "-q", "-m", "explicit synthetic secret")
        result = self.try_repo_hook("gitleaks-files-system", "--all-files")
        self.assertNotEqual(0, result.returncode, result.stdout)
        self.assertIn("leaks found", result.stdout.lower())


if __name__ == "__main__":
    unittest.main()
