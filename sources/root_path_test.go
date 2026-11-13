package sources

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/fatih/semgroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
)

func TestRootRelativePath(t *testing.T) {
	temp := t.TempDir()
	tests := map[string]struct {
		root, path, want string
	}{
		"no root":                          {root: "", path: "src/x.py", want: ""},
		"relative root":                    {root: "repo", path: filepath.Join("repo", "src", "x.py"), want: "src/x.py"},
		"dot root":                         {root: ".", path: filepath.Join("src", "x.py"), want: "src/x.py"},
		"absolute root":                    {root: temp, path: filepath.Join(temp, "src", "x.py"), want: "src/x.py"},
		"slash separated path under root":  {root: "repo", path: "repo/src/x.py", want: "src/x.py"},
		"the root itself":                  {root: "repo", path: "repo", want: ""},
		"outside the root":                 {root: "repo/src", path: "repo/other/x.py", want: ""},
		"parent prefix is outside":         {root: "repo", path: "repo-other/x.py", want: ""},
		"parent traversal is outside":      {root: "repo", path: "repo/../other/x.py", want: ""},
		"root with dots":                   {root: "repo/src/..", path: "repo/file.py", want: "file.py"},
		"space and unicode":                {root: "repo dir", path: "repo dir/秘密/x y.py", want: "秘密/x y.py"},
		"archive member keeps inner path":  {root: "repo", path: "repo/bundle.tar!inner/x.py", want: "bundle.tar!inner/x.py"},
		"archive separator in root":        {root: "repo!dir", path: "repo!dir/bundle.tar!inner/x.py", want: "bundle.tar!inner/x.py"},
		"archive separator in directories": {root: "repo", path: "repo/a!b/c!d/x.py", want: "a!b/c!d/x.py"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, RootRelativePath(tt.root, tt.path))
		})
	}

	if runtime.GOOS == "windows" {
		t.Run("windows volume", func(t *testing.T) {
			root := `C:\Dev\repo`
			assert.Equal(t, "src/x.py", RootRelativePath(root, `C:/Dev/repo/src/x.py`))
			assert.Equal(t, "src/x.py", RootRelativePath(root, `c:\dev\repo\src\x.py`))
			assert.Empty(t, RootRelativePath(root, `D:\Dev\repo\src\x.py`))
		})
	}
}

func TestFilesRootRelativeAllowlistAndFragmentRoot(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "src", "vendor"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "src", "nested"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "src", "vendor", "lib.py"), []byte("vendor\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "src", "nested", "app.py"), []byte("app\n"), 0o600))

	allowlist := &config.Allowlist{Paths: []*regexp.Regexp{regexp.MustCompile(`^src/vendor/`)}}
	require.NoError(t, allowlist.Validate())
	cfg := &config.Config{Allowlists: []*config.Allowlist{allowlist}}

	sourcesToScan := map[string]string{"absolute root": root}
	if relative, err := relativeToCwd(root); err == nil {
		sourcesToScan["relative root"] = relative
	}
	for name, source := range sourcesToScan {
		t.Run(name, func(t *testing.T) {
			files := &Files{Config: cfg, Path: source, Sema: semgroup.NewGroup(context.Background(), 1)}
			var scanned []string
			err := files.scanTargets(context.Background(), func(target ScanTarget, err error) error {
				require.NoError(t, err)
				scanned = append(scanned, RootRelativePath(source, target.Path))
				return nil
			})
			require.NoError(t, err)
			assert.Equal(t, []string{"src/nested/app.py"}, scanned)
		})
	}

	files := &Files{Path: root, Sema: semgroup.NewGroup(context.Background(), 1)}
	var fragments []Fragment
	err := files.Fragments(context.Background(), func(fragment Fragment, err error) error {
		require.NoError(t, err)
		fragments = append(fragments, fragment)
		return nil
	})
	require.NoError(t, err)
	require.Len(t, fragments, 2)
	for _, fragment := range fragments {
		assert.Equal(t, root, fragment.Root)
		assert.NotEmpty(t, RootRelativePath(fragment.Root, fragment.FilePath))
	}

	single := &Files{Path: filepath.Join(root, "src", "nested", "app.py"), Sema: semgroup.NewGroup(context.Background(), 1)}
	fragments = nil
	err = single.Fragments(context.Background(), func(fragment Fragment, err error) error {
		require.NoError(t, err)
		fragments = append(fragments, fragment)
		return nil
	})
	require.NoError(t, err)
	require.Len(t, fragments, 1)
	assert.Empty(t, fragments[0].Root)
}

func TestFilesRootRelativeAllowlistKeepsAndSemantics(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "src"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "src", "app.py"), []byte("app\n"), 0o600))

	allowlist := &config.Allowlist{
		MatchCondition: config.AllowlistMatchAnd,
		Paths:          []*regexp.Regexp{regexp.MustCompile(`^src/app\.py$`)},
		Regexes:        []*regexp.Regexp{regexp.MustCompile(`NEVERMATCH`)},
	}
	require.NoError(t, allowlist.Validate())
	cfg := &config.Config{Allowlists: []*config.Allowlist{allowlist}}

	files := &Files{Config: cfg, Path: root, Sema: semgroup.NewGroup(context.Background(), 1)}
	var scanned []string
	err := files.scanTargets(context.Background(), func(target ScanTarget, err error) error {
		require.NoError(t, err)
		scanned = append(scanned, target.Path)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{filepath.Join(root, "src", "app.py")}, scanned,
		"a relative path match must not bypass the required AND content condition")
}

func TestFilesRootRelativeAllowlistWindowsSeparators(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows path spelling")
	}
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "src", "vendor"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "src", "vendor", "lib.py"), []byte("vendor\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "src", "app.py"), []byte("app\n"), 0o600))

	allowlist := &config.Allowlist{Paths: []*regexp.Regexp{regexp.MustCompile(`^src\\vendor\\`)}}
	require.NoError(t, allowlist.Validate())
	cfg := &config.Config{Allowlists: []*config.Allowlist{allowlist}}

	files := &Files{Config: cfg, Path: root, Sema: semgroup.NewGroup(context.Background(), 1)}
	var scanned []string
	err := files.scanTargets(context.Background(), func(target ScanTarget, err error) error {
		require.NoError(t, err)
		scanned = append(scanned, filepath.Base(target.Path))
		return nil
	})
	require.NoError(t, err)
	sort.Strings(scanned)
	assert.Equal(t, []string{"app.py"}, scanned)
}

func TestFilesArchiveMembersCarryRoot(t *testing.T) {
	root := t.TempDir()
	archive := makeZip(t, "src/vendor/lib.py", "vendor\n")
	require.NoError(t, os.WriteFile(filepath.Join(root, "bundle.zip"), archive, 0o600))

	files := &Files{Path: root, MaxArchiveDepth: 2, Sema: semgroup.NewGroup(context.Background(), 1)}
	var fragments []Fragment
	err := files.Fragments(context.Background(), func(fragment Fragment, err error) error {
		require.NoError(t, err)
		fragments = append(fragments, fragment)
		return nil
	})
	require.NoError(t, err)
	require.Len(t, fragments, 1)
	assert.Equal(t, root, fragments[0].Root)
	assert.Equal(t, "bundle.zip!src/vendor/lib.py", RootRelativePath(fragments[0].Root, fragments[0].FilePath))
}

func relativeToCwd(path string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Rel(cwd, path)
}
