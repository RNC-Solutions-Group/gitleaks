package detect

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/cmd/scm"
	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"
)

const stagedTestPAT = "ghp_Zq8xR2mV7kLp4Wn9Tb3Yc6Hd1Fs5Ga0Ju2Ke"

func initStagedTestRepo(t *testing.T) string {
	t.Helper()
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not available")
	}

	repo := t.TempDir()
	runGit := func(args ...string) string {
		t.Helper()
		cmd := exec.Command(gitPath, append([]string{"-c", "core.fsmonitor=false", "-C", repo}, args...)...)
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = &output
		require.NoError(t, cmd.Run(), "git %v: %s", args, output.String())
		return output.String()
	}
	runGit("init", "-q")
	runGit("config", "core.fsmonitor", "false")
	runGit("config", "user.email", "cpv4@example.com")
	runGit("config", "user.name", "CPV4")
	runGit("commit", "-q", "--allow-empty", "-m", "root")
	return repo
}

func stageFile(t *testing.T, repo, name, content string) {
	t.Helper()
	path := filepath.Join(repo, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	runGitInRepo(t, repo, "add", "--", name)
}

func runGitInRepo(t *testing.T, repo string, args ...string) string {
	t.Helper()
	gitPath, err := exec.LookPath("git")
	require.NoError(t, err)
	cmd := exec.Command(gitPath, append([]string{"-c", "core.fsmonitor=false", "-C", repo}, args...)...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	require.NoError(t, cmd.Run(), "git %v: %s", args, output.String())
	return output.String()
}

func stagedDetector(t *testing.T, repo string, staged bool, maxArchiveDepth int) ([]report.Finding, error) {
	t.Helper()
	detector, err := NewDetectorDefaultConfig()
	require.NoError(t, err)
	detector.MaxArchiveDepth = maxArchiveDepth

	var gitCmd *sources.GitCmd
	if staged {
		gitCmd, err = sources.NewGitDiffCmdContext(context.Background(), repo, true)
	} else {
		gitCmd, err = sources.NewGitLogCmdContext(context.Background(), repo, "")
	}
	require.NoError(t, err)
	findings, err := detector.DetectSource(context.Background(), &sources.Git{
		Cmd:             gitCmd,
		Config:          &detector.Config,
		Remote:          sources.NewRemoteInfo(scm.NoPlatform, repo),
		Sema:            detector.Sema,
		MaxArchiveDepth: detector.MaxArchiveDepth,
	})
	return findings, err
}

func TestStagedTextMarkedBinaryIsScanned(t *testing.T) {
	t.Run("control text diff", func(t *testing.T) {
		repo := initStagedTestRepo(t)
		stageFile(t, repo, "secret.txt", "token = \""+stagedTestPAT+"\"\n")
		findings, err := stagedDetector(t, repo, true, 0)
		require.NoError(t, err)
		require.Len(t, findings, 1)
		assert.Equal(t, stagedTestPAT, findings[0].Secret)
	})

	t.Run("gitattributes -diff", func(t *testing.T) {
		repo := initStagedTestRepo(t)
		stageFile(t, repo, ".gitattributes", "secret.txt -diff\n")
		stageFile(t, repo, "secret.txt", "token = \""+stagedTestPAT+"\"\n")
		diff := runGitInRepo(t, repo, "diff", "-U0", "--no-ext-diff", "--staged", ".")
		require.Contains(t, diff, "Binary files")
		findings, err := stagedDetector(t, repo, true, 0)
		require.NoError(t, err)
		require.Len(t, findings, 1)
		assert.Equal(t, "secret.txt", findings[0].File)
	})

	t.Run("core.bigFileThreshold", func(t *testing.T) {
		repo := initStagedTestRepo(t)
		stageFile(t, repo, "secret.txt", "token = \""+stagedTestPAT+"\"\n")
		runGitInRepo(t, repo, "config", "core.bigFileThreshold", "10")
		diff := runGitInRepo(t, repo, "diff", "-U0", "--no-ext-diff", "--staged", ".")
		require.Contains(t, diff, "Binary files")
		findings, err := stagedDetector(t, repo, true, 0)
		require.NoError(t, err)
		require.Len(t, findings, 1)
		assert.Equal(t, "secret.txt", findings[0].File)
	})
}

func TestStagedActualBinaryRemainsExcluded(t *testing.T) {
	repo := initStagedTestRepo(t)
	binary := append([]byte("\x89PNG\r\n\x1a\n\x00"), []byte("token = \""+stagedTestPAT+"\"\n")...)
	stageFile(t, repo, "image.png", string(binary))
	diff := runGitInRepo(t, repo, "diff", "-U0", "--no-ext-diff", "--staged", ".")
	require.Contains(t, diff, "Binary files")
	findings, err := stagedDetector(t, repo, true, 0)
	require.NoError(t, err)
	assert.Empty(t, findings)
}

func TestStagedScanUsesIndexNotWorktree(t *testing.T) {
	t.Run("worktree is clean but index contains secret", func(t *testing.T) {
		repo := initStagedTestRepo(t)
		stageFile(t, repo, "secret.txt", "token = \""+stagedTestPAT+"\"\n")
		require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.txt"), []byte("clean worktree\n"), 0o600))
		findings, err := stagedDetector(t, repo, true, 0)
		require.NoError(t, err)
		require.Len(t, findings, 1)
		assert.Equal(t, stagedTestPAT, findings[0].Secret)
	})

	t.Run("worktree contains secret but index is clean", func(t *testing.T) {
		repo := initStagedTestRepo(t)
		stageFile(t, repo, "secret.txt", "clean index\n")
		require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.txt"), []byte("token = \""+stagedTestPAT+"\"\n"), 0o600))
		findings, err := stagedDetector(t, repo, true, 0)
		require.NoError(t, err)
		assert.Empty(t, findings)
	})
}

func TestStagedBinaryMarkedUsesIndexNotWorktree(t *testing.T) {
	repo := initStagedTestRepo(t)
	stageFile(t, repo, ".gitattributes", "secret.txt -diff\n")
	stageFile(t, repo, "secret.txt", "token = \""+stagedTestPAT+"\"\n")
	require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.txt"), []byte("clean worktree\n"), 0o600))

	diff := runGitInRepo(t, repo, "diff", "-U0", "--no-ext-diff", "--staged", ".")
	require.Contains(t, diff, "Binary files")
	findings, err := stagedDetector(t, repo, true, 0)
	require.NoError(t, err)
	require.Len(t, findings, 1)
	assert.Equal(t, stagedTestPAT, findings[0].Secret)
}
func TestStagedBinaryMarkedCleanIndexIgnoresWorktreeSecret(t *testing.T) {
	repo := initStagedTestRepo(t)
	stageFile(t, repo, ".gitattributes", "secret.txt -diff\n")
	stageFile(t, repo, "secret.txt", "clean index\n")
	require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.txt"), []byte("token = \""+stagedTestPAT+"\"\n"), 0o600))

	diff := runGitInRepo(t, repo, "diff", "-U0", "--no-ext-diff", "--staged", ".")
	require.Contains(t, diff, "Binary files")
	findings, err := stagedDetector(t, repo, true, 0)
	require.NoError(t, err)
	assert.Empty(t, findings)
}
func TestStagedDeletionDoesNotReadBlob(t *testing.T) {
	repo := initStagedTestRepo(t)
	stageFile(t, repo, "secret.txt", "token = \""+stagedTestPAT+"\"\n")
	runGitInRepo(t, repo, "rm", "-q", "-f", "secret.txt")
	findings, err := stagedDetector(t, repo, true, 0)
	require.NoError(t, err)
	assert.Empty(t, findings)
}

func TestStagedRenameScansNewPath(t *testing.T) {
	repo := initStagedTestRepo(t)
	stageFile(t, repo, "old.txt", "clean\n")
	runGitInRepo(t, repo, "commit", "-q", "-m", "old")
	runGitInRepo(t, repo, "mv", "old.txt", "new.txt")
	require.NoError(t, os.WriteFile(filepath.Join(repo, "new.txt"), []byte("token = \""+stagedTestPAT+"\"\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(repo, ".gitattributes"), []byte("new.txt -diff\n"), 0o600))
	runGitInRepo(t, repo, "add", "--", "new.txt", ".gitattributes")

	findings, err := stagedDetector(t, repo, true, 0)
	require.NoError(t, err)
	require.NotEmpty(t, findings)
	assert.Equal(t, "new.txt", findings[0].File)
}

func TestStagedPathsWithSpacesAndUnicode(t *testing.T) {
	repo := initStagedTestRepo(t)
	stageFile(t, repo, ".gitattributes", "\"space name.txt\" -diff\n\"unicodé.txt\" -diff\n")
	stageFile(t, repo, "space name.txt", "token = \""+stagedTestPAT+"\"\n")
	stageFile(t, repo, "unicodé.txt", "token = \""+stagedTestPAT+"\"\n")

	findings, err := stagedDetector(t, repo, true, 0)
	require.NoError(t, err)
	require.Len(t, findings, 2)
	files := map[string]bool{}
	for _, finding := range findings {
		files[finding.File] = true
	}
	assert.True(t, files["space name.txt"])
	assert.True(t, files["unicodé.txt"])
}

func TestHistoryBinaryBehaviorRemainsUnchanged(t *testing.T) {
	repo := initStagedTestRepo(t)
	stageFile(t, repo, ".gitattributes", "secret.txt -diff\n")
	stageFile(t, repo, "secret.txt", "token = \""+stagedTestPAT+"\"\n")
	runGitInRepo(t, repo, "commit", "-q", "-m", "binary-marked text")

	findings, err := stagedDetector(t, repo, false, 0)
	require.NoError(t, err)
	assert.Empty(t, findings)
}

func TestStagedArchiveIsExtracted(t *testing.T) {
	repo := initStagedTestRepo(t)
	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)
	entry, err := zipWriter.Create("inside.txt")
	require.NoError(t, err)
	_, err = entry.Write([]byte("token = \"" + stagedTestPAT + "\"\n"))
	require.NoError(t, err)
	require.NoError(t, zipWriter.Close())
	stageFile(t, repo, "archive.zip", archive.String())

	findings, err := stagedDetector(t, repo, true, 2)
	require.NoError(t, err)
	require.NotEmpty(t, findings)
	assert.Contains(t, findings[0].File, "archive.zip!")
}

func TestUnstagedBinaryBehaviorRemainsUnchanged(t *testing.T) {
	repo := initStagedTestRepo(t)
	stageFile(t, repo, ".gitattributes", "secret.txt -diff\n")
	stageFile(t, repo, "secret.txt", "clean\n")
	runGitInRepo(t, repo, "commit", "-q", "-m", "binary attribute")
	require.NoError(t, os.WriteFile(filepath.Join(repo, "secret.txt"), []byte("token = \""+stagedTestPAT+"\"\n"), 0o600))

	detector, err := NewDetectorDefaultConfig()
	require.NoError(t, err)
	gitCmd, err := sources.NewGitDiffCmdContext(context.Background(), repo, false)
	require.NoError(t, err)
	findings, err := detector.DetectSource(context.Background(), &sources.Git{
		Cmd:             gitCmd,
		Config:          &detector.Config,
		Remote:          sources.NewRemoteInfo(scm.NoPlatform, repo),
		Sema:            detector.Sema,
		MaxArchiveDepth: 0,
	})
	require.NoError(t, err)
	assert.Empty(t, findings)
}

func TestStagedArchiveRespectsDepthPolicy(t *testing.T) {
	repo := initStagedTestRepo(t)
	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)
	entry, err := zipWriter.Create("inside.txt")
	require.NoError(t, err)
	_, err = entry.Write([]byte("token = \"" + stagedTestPAT + "\"\n"))
	require.NoError(t, err)
	require.NoError(t, zipWriter.Close())
	stageFile(t, repo, "archive.zip", archive.String())

	findings, err := stagedDetector(t, repo, true, 0)
	require.NoError(t, err)
	assert.Empty(t, findings)
}

func stagedFindingsWithMaxFileSize(t *testing.T, repo string, maxFileSize int) ([]report.Finding, error) {
	t.Helper()
	detector, err := NewDetectorDefaultConfig()
	require.NoError(t, err)
	gitCmd, err := sources.NewGitDiffCmdContext(context.Background(), repo, true)
	require.NoError(t, err)
	return detector.DetectSource(context.Background(), &sources.Git{
		Cmd:             gitCmd,
		Config:          &detector.Config,
		Remote:          sources.NewRemoteInfo(scm.NoPlatform, repo),
		Sema:            detector.Sema,
		MaxFileSize:     maxFileSize,
		MaxArchiveDepth: 0,
	})
}

func TestStagedBinaryMarkedBlobHonorsMaxFileSize(t *testing.T) {
	repo := initStagedTestRepo(t)
	stageFile(t, repo, ".gitattributes", "secret.txt -diff\n")
	stageFile(t, repo, "secret.txt", strings.Repeat("A", 512)+" token = \""+stagedTestPAT+"\"\n")

	findings, err := stagedFindingsWithMaxFileSize(t, repo, 100)
	require.NoError(t, err)
	assert.Empty(t, findings)

	findings, err = stagedFindingsWithMaxFileSize(t, repo, 4096)
	require.NoError(t, err)
	require.NotEmpty(t, findings)
	assert.Equal(t, stagedTestPAT, findings[0].Secret)
}
