package detect

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/sources"
)

// TestGitLargeDiffCompletes guards the streaming parser against a large patch
// without changing the default history scope. The only credential is in the
// last file, so a truncated parser would produce a false-clean result.
func TestGitLargeDiffCompletes(t *testing.T) {
	dir, run := gitTestRepo(t)
	for i := 0; i < 200; i++ {
		name := fmt.Sprintf("ordinary-%03d.txt", i)
		content := fmt.Sprintf("file=%d\n%s\n", i, "synthetic content for parser load\n")
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
	}
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "zz-last.txt"),
		[]byte("aws_access_key_id = AKIAZQ2M4X6P7R2T4V6Z\n"),
		0o600,
	))
	run("add", ".")
	run("commit", "-q", "-m", "large patch")

	findings := scanGitWithOptions(t, dir, sources.GitLogOptions{})
	require.Len(t, findings, 1)
	assert.Equal(t, "zz-last.txt", findings[0].File)
	assert.Equal(t, "aws-access-token", findings[0].RuleID)
}
