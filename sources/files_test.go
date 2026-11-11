package sources

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
)

func TestScanTargetsHonorsAndAllowlist(t *testing.T) {
	dir := t.TempDir()
	secretPath := filepath.Join(dir, "secret.txt")
	cleanPath := filepath.Join(dir, "clean.txt")
	require.NoError(t, os.WriteFile(secretPath, []byte("aws_access_key_id = \"AKIAQ7Z3R5T7Y9U2I4O6\"\n"), 0o600))
	require.NoError(t, os.WriteFile(cleanPath, []byte("nothing\n"), 0o600))

	allowlist := &config.Allowlist{
		MatchCondition: config.AllowlistMatchAnd,
		Paths:          []*regexp.Regexp{regexp.MustCompile(`secret\.txt$`)},
		Regexes:        []*regexp.Regexp{regexp.MustCompile(`NEVERMATCH`)},
	}
	require.NoError(t, allowlist.Validate())
	cfg := &config.Config{Allowlists: []*config.Allowlist{allowlist}}

	files := &Files{Config: cfg, Path: dir}
	var scanned []string
	err := files.scanTargets(context.Background(), func(target ScanTarget, err error) error {
		require.NoError(t, err)
		scanned = append(scanned, target.Path)
		return nil
	})
	require.NoError(t, err)
	require.Contains(t, scanned, secretPath)
	require.Contains(t, scanned, cleanPath)
}
