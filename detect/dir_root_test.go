package detect

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
	"github.com/zricethezav/gitleaks/v8/regexp"
	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"
)

const dirRootSecret = "RNC_TEST_SECRET_aaaaaaaaaaaaaaaa"

func dirRootConfig(allowlists ...*config.Allowlist) *config.Config {
	for _, allowlist := range allowlists {
		if err := allowlist.Validate(); err != nil {
			panic(err)
		}
	}
	return &config.Config{
		Rules: map[string]config.Rule{
			"cpv8-test": {
				RuleID:      "cpv8-test",
				Description: "synthetic CP-V8 root test credential",
				Regex:       regexp.MustCompile(`RNC_TEST_SECRET_[A-Za-z0-9]{16}`),
				Keywords:    []string{"rnc_test_secret_"},
			},
		},
		Keywords:     map[string]struct{}{"rnc_test_secret_": {}},
		OrderedRules: []string{"cpv8-test"},
		Allowlists:   allowlists,
	}
}

func dirRootTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "src", "nested"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "src", "vendor"), 0o755))
	write := func(name, content string) {
		t.Helper()
		require.NoError(t, os.WriteFile(filepath.Join(root, filepath.FromSlash(name)), []byte(content), 0o600))
	}
	write("src/app.py", "Token = \""+dirRootSecret+"\"\n")
	write("src/nested/deep.py", "Token = \"RNC_TEST_SECRET_bbbbbbbbbbbbbbbb\"\n")
	write("src/vendor/lib.py", "Token = \"RNC_TEST_SECRET_dddddddddddddddd\"\n")
	write("clean.txt", "nothing to see here\n")
	return root
}

func scanDirWith(t *testing.T, detector *Detector, cfg *config.Config, source string, maxArchiveDepth int) []report.Finding {
	t.Helper()
	findings, err := detector.DetectSource(context.Background(), &sources.Files{
		Config:          cfg,
		Path:            source,
		Sema:            detector.Sema,
		MaxArchiveDepth: maxArchiveDepth,
		Informative:     detector.Informative,
	})
	require.NoError(t, err)
	return findings
}

func scanDirRoot(t *testing.T, cfg *config.Config, source string) []report.Finding {
	t.Helper()
	return scanDirWith(t, NewDetector(*cfg), cfg, source, 0)
}

func dirRelativeNames(t *testing.T, root string, findings []report.Finding) []string {
	t.Helper()
	names := make([]string, 0, len(findings))
	for _, finding := range findings {
		rel := sources.RootRelativePath(root, finding.File)
		require.NotEmpty(t, rel, "finding %q is not under root %q", finding.File, root)
		names = append(names, rel)
	}
	sort.Strings(names)
	return names
}

func TestDirRootFindingsKeepIdentityAcrossRootSpellings(t *testing.T) {
	root := dirRootTree(t)
	cfg := dirRootConfig()
	expected := []string{"src/app.py", "src/nested/deep.py", "src/vendor/lib.py"}

	assertIdentity := func(t *testing.T, expectedRoot string, findings []report.Finding) {
		t.Helper()
		assert.Equal(t, expected, dirRelativeNames(t, expectedRoot, findings))
		for _, finding := range findings {
			assert.Equal(t, fmt.Sprintf("%s:%s:%d", finding.File, finding.RuleID, finding.StartLine), finding.Fingerprint)
		}
	}

	t.Run("absolute root", func(t *testing.T) {
		assertIdentity(t, root, scanDirRoot(t, cfg, root))
	})

	t.Run("relative root", func(t *testing.T) {
		relative := relativeScanRoot(t, root)
		assertIdentity(t, relative, scanDirRoot(t, cfg, relative))
	})

	t.Run("dot root", func(t *testing.T) {
		t.Chdir(root)
		findings := scanDirRoot(t, cfg, ".")
		assertIdentity(t, ".", findings)
	})
}

func TestDirRootGitleaksIgnoreMatchesRelativeFingerprint(t *testing.T) {
	root := dirRootTree(t)
	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitleaksignore"),
		[]byte("src/app.py:cpv8-test:1\n"), 0o600))
	cfg := dirRootConfig()

	run := func(t *testing.T, source string) {
		t.Helper()
		detector := NewDetector(*cfg)
		require.NoError(t, detector.AddGitleaksIgnore(filepath.Join(root, ".gitleaksignore")))
		findings := scanDirWith(t, detector, cfg, source, 0)
		assert.Equal(t, []string{"src/nested/deep.py", "src/vendor/lib.py"}, dirRelativeNames(t, source, findings))
	}

	t.Run("absolute root", func(t *testing.T) { run(t, root) })
	t.Run("relative root", func(t *testing.T) {
		run(t, relativeScanRoot(t, root))
	})
	t.Run("dot root", func(t *testing.T) {
		t.Chdir(root)
		run(t, ".")
	})

	t.Run("wrong line remains visible", func(t *testing.T) {
		require.NoError(t, os.WriteFile(filepath.Join(root, ".gitleaksignore"),
			[]byte("src/app.py:cpv8-test:2\n"), 0o600))
		detector := NewDetector(*cfg)
		require.NoError(t, detector.AddGitleaksIgnore(filepath.Join(root, ".gitleaksignore")))
		findings := scanDirWith(t, detector, cfg, root, 0)
		assert.Contains(t, dirRelativeNames(t, root, findings), "src/app.py")
	})
}

func TestDirRootAllowlistPathSemantics(t *testing.T) {
	root := dirRootTree(t)

	check := func(t *testing.T, cfg *config.Config, expected []string) {
		t.Helper()
		sources := map[string]string{
			"absolute": root,
			"relative": relativeScanRoot(t, root),
		}
		for name, source := range sources {
			t.Run(name, func(t *testing.T) {
				findings := scanDirRoot(t, cfg, source)
				assert.Equal(t, expected, dirRelativeNames(t, source, findings))
			})
		}
	}

	t.Run("global path allowlist", func(t *testing.T) {
		allowlist := &config.Allowlist{Paths: []*regexp.Regexp{regexp.MustCompile(`^src/vendor/`)}}
		check(t, dirRootConfig(allowlist), []string{"src/app.py", "src/nested/deep.py"})
	})

	t.Run("rule path allowlist", func(t *testing.T) {
		allowlist := &config.Allowlist{Paths: []*regexp.Regexp{regexp.MustCompile(`^src/vendor/`)}}
		cfg := dirRootConfig()
		rule := cfg.Rules["cpv8-test"]
		rule.Allowlists = []*config.Allowlist{allowlist}
		cfg.Rules["cpv8-test"] = rule
		check(t, cfg, []string{"src/app.py", "src/nested/deep.py"})
	})

	t.Run("AND path and regex", func(t *testing.T) {
		allowlist := &config.Allowlist{
			MatchCondition: config.AllowlistMatchAnd,
			Paths:          []*regexp.Regexp{regexp.MustCompile(`^src/app\.py$`)},
			Regexes:        []*regexp.Regexp{regexp.MustCompile(`RNC_TEST_SECRET_a{16}`)},
		}
		check(t, dirRootConfig(allowlist), []string{"src/nested/deep.py", "src/vendor/lib.py"})
	})

	t.Run("AND path match does not bypass nonmatching regex", func(t *testing.T) {
		allowlist := &config.Allowlist{
			MatchCondition: config.AllowlistMatchAnd,
			Paths:          []*regexp.Regexp{regexp.MustCompile(`^src/app\.py$`)},
			Regexes:        []*regexp.Regexp{regexp.MustCompile(`NEVERMATCH`)},
		}
		check(t, dirRootConfig(allowlist), []string{"src/app.py", "src/nested/deep.py", "src/vendor/lib.py"})
	})

	t.Run("OR path and regex", func(t *testing.T) {
		allowlist := &config.Allowlist{
			MatchCondition: config.AllowlistMatchOr,
			Paths:          []*regexp.Regexp{regexp.MustCompile(`^src/app\.py$`)},
			Regexes:        []*regexp.Regexp{regexp.MustCompile(`NEVERMATCH`)},
		}
		check(t, dirRootConfig(allowlist), []string{"src/nested/deep.py", "src/vendor/lib.py"})
	})

	if runtime.GOOS == "windows" {
		t.Run("native separator pattern", func(t *testing.T) {
			allowlist := &config.Allowlist{Paths: []*regexp.Regexp{regexp.MustCompile(`^src\\vendor\\`)}}
			check(t, dirRootConfig(allowlist), []string{"src/app.py", "src/nested/deep.py"})
		})
	}
}

func TestDirRootBaselineCompatibility(t *testing.T) {
	root := dirRootTree(t)
	cfg := dirRootConfig()

	removeBaseline := func() {
		t.Helper()
		err := os.Remove(filepath.Join(root, "baseline.json"))
		if err != nil && !os.IsNotExist(err) {
			require.NoError(t, err)
		}
	}
	dotFindings := func() []report.Finding {
		t.Helper()
		removeBaseline()
		t.Chdir(root)
		findings := scanDirRoot(t, cfg, ".")
		require.Len(t, findings, 3)
		return findings
	}
	absoluteFindings := func() []report.Finding {
		t.Helper()
		removeBaseline()
		findings := scanDirRoot(t, cfg, root)
		require.Len(t, findings, 3)
		return findings
	}
	writeBaseline := func(findings []report.Finding) string {
		t.Helper()
		data, err := json.Marshal(findings)
		require.NoError(t, err)
		path := filepath.Join(root, "baseline.json")
		require.NoError(t, os.WriteFile(path, data, 0o600))
		return path
	}

	t.Run("dot baseline with absolute root", func(t *testing.T) {
		baseline := writeBaseline(dotFindings())
		detector := NewDetector(*cfg)
		require.NoError(t, detector.AddBaseline(baseline, root))
		assert.Empty(t, scanDirWith(t, detector, cfg, root, 0))
	})

	t.Run("absolute baseline with dot root", func(t *testing.T) {
		baseline := writeBaseline(absoluteFindings())
		detector := NewDetector(*cfg)
		require.NoError(t, detector.AddBaseline(baseline, "."))
		t.Chdir(root)
		assert.Empty(t, scanDirWith(t, detector, cfg, ".", 0))
	})

	t.Run("new finding remains visible", func(t *testing.T) {
		baseline := writeBaseline(absoluteFindings())
		require.NoError(t, os.WriteFile(filepath.Join(root, "src", "new.py"),
			[]byte("Token = \"RNC_TEST_SECRET_eeeeeeeeeeeeeeee\"\n"), 0o600))
		detector := NewDetector(*cfg)
		require.NoError(t, detector.AddBaseline(baseline, root))
		findings := scanDirWith(t, detector, cfg, root, 0)
		require.Len(t, findings, 1)
		assert.Equal(t, "src/new.py", sources.RootRelativePath(root, findings[0].File))
	})
}

func TestDirRootRelativeBaselineWorksWithDotRoot(t *testing.T) {
	root := dirRootTree(t)
	cfg := dirRootConfig()
	parent := filepath.Dir(root)
	name := filepath.Base(root)

	t.Chdir(parent)
	findings := scanDirRoot(t, cfg, name)
	require.Len(t, findings, 3)
	for _, finding := range findings {
		assert.True(t, strings.HasPrefix(finding.File, name+"/"), finding.File)
	}
	encoded, err := json.Marshal(findings)
	require.NoError(t, err)
	baseline := filepath.Join(root, "baseline.json")
	require.NoError(t, os.WriteFile(baseline, encoded, 0o600))

	t.Chdir(root)
	detector := NewDetector(*cfg)
	require.NoError(t, detector.AddBaseline("baseline.json", "."))
	assert.Empty(t, scanDirWith(t, detector, cfg, ".", 0))
}

func TestDirRootBaselineNormalizationStaysUnderSource(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	baselinePath := filepath.Join(outside, "baseline.json")
	baseline := []report.Finding{
		{File: filepath.Join(outside, "other.py"), RuleID: "cpv8-test"},
		{File: filepath.Join(root, "src", "inside.py"), RuleID: "cpv8-test"},
		{File: filepath.Join(root, "bundle.zip") + "!inner/x.py", RuleID: "cpv8-test"},
	}
	encoded, err := json.Marshal(baseline)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(baselinePath, encoded, 0o600))
	before, err := os.ReadFile(baselinePath)
	require.NoError(t, err)

	detector := NewDetector(*dirRootConfig())
	require.NoError(t, detector.AddBaseline(baselinePath, root))
	require.Len(t, detector.baseline, 3)
	assert.Equal(t, filepath.Join(outside, "other.py"), detector.baseline[0].File,
		"entries outside the source must not be rewritten")
	assert.Equal(t, "src/inside.py", detector.baseline[1].File)
	assert.Equal(t, "bundle.zip!inner/x.py", detector.baseline[2].File)

	after, err := os.ReadFile(baselinePath)
	require.NoError(t, err)
	assert.Equal(t, before, after, "normalization is in memory only")
}

func TestDirRootConfigSelfExclusion(t *testing.T) {
	root := dirRootTree(t)
	configPath := filepath.Join(root, ".gitleaks.toml")
	require.NoError(t, os.WriteFile(configPath,
		[]byte("# synthetic credential in config: RNC_TEST_SECRET_cccccccccccccccc\n"), 0o600))

	cfg := dirRootConfig()
	cfg.Path = configPath

	findings := scanDirRoot(t, cfg, root)
	assert.Equal(t, []string{"src/app.py", "src/nested/deep.py", "src/vendor/lib.py"},
		dirRelativeNames(t, root, findings))
	for _, finding := range findings {
		assert.NotContains(t, finding.File, ".gitleaks.toml")
	}

	relative := relativeScanRoot(t, root)
	cfg.Path = filepath.Join(relative, ".gitleaks.toml")
	findings = scanDirRoot(t, cfg, relative)
	assert.Len(t, findings, 3)
	for _, finding := range findings {
		assert.NotContains(t, finding.File, ".gitleaks.toml")
	}
}

// relativeScanRoot returns a genuine relative scan root on every supported
// platform. Windows temporary directories can be on a different volume from the
// package, so a cwd-relative spelling is impossible there; move to the fixture's
// parent and use its one-component name instead of skipping relative-root
// coverage.
func relativeScanRoot(t *testing.T, path string) string {
	t.Helper()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	rel, err := filepath.Rel(cwd, path)
	if err == nil {
		return rel
	}
	t.Chdir(filepath.Dir(path))
	return filepath.Base(path)
}

func TestDirRootSingleFileHasNoRoot(t *testing.T) {
	root := dirRootTree(t)
	path := filepath.Join(root, "src", "app.py")
	cfg := dirRootConfig(&config.Allowlist{
		Paths: []*regexp.Regexp{regexp.MustCompile(`^src/app\.py$`)},
	})

	findings := scanDirRoot(t, cfg, path)
	require.Len(t, findings, 1, "a single-file target must not borrow a directory root")
	assert.Equal(t, filepath.ToSlash(path), findings[0].File)

	detector := NewDetector(*cfg)
	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitleaksignore"),
		[]byte("src/app.py:cpv8-test:1\n"), 0o600))
	require.NoError(t, detector.AddGitleaksIgnore(filepath.Join(root, ".gitleaksignore")))
	findings = scanDirWith(t, detector, cfg, path, 0)
	require.Len(t, findings, 1)
}

func TestDirRootArchiveMemberMatching(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "bundle.zip")
	writeZip(t, archivePath, map[string]string{
		"src/vendor/lib.py": "Token = \"" + dirRootSecret + "\"\n",
	})

	cfg := dirRootConfig()
	findings := scanDirWith(t, NewDetector(*cfg), cfg, root, 2)
	require.Len(t, findings, 1)
	assert.Equal(t, "bundle.zip!src/vendor/lib.py",
		sources.RootRelativePath(root, findings[0].File))

	allowlist := &config.Allowlist{Paths: []*regexp.Regexp{regexp.MustCompile(`^bundle\.zip!src/vendor/`)}}
	allowlisted := dirRootConfig(allowlist)
	findings = scanDirWith(t, NewDetector(*allowlisted), allowlisted, root, 2)
	assert.Empty(t, findings, "an archive member path is matched relative to the scan root")

	cfg.Path = filepath.Join(root, ".gitleaks.toml")
	require.NoError(t, os.WriteFile(cfg.Path,
		[]byte("# synthetic credential in config: RNC_TEST_SECRET_cccccccccccccccc\n"), 0o600))
	findings = scanDirWith(t, NewDetector(*cfg), cfg, root, 2)
	require.Len(t, findings, 1)
	assert.NotContains(t, findings[0].File, ".gitleaks.toml")
}

func TestDirRootNestedArchivePath(t *testing.T) {
	root := t.TempDir()
	inner := zipStringBytes(t, map[string]string{
		"src/vendor/lib.py": "Token = \"" + dirRootSecret + "\"\n",
	})
	writeZipBytes(t, filepath.Join(root, "outer.zip"), map[string][]byte{"inner.zip": inner})

	cfg := dirRootConfig(&config.Allowlist{
		Paths: []*regexp.Regexp{regexp.MustCompile(`^outer\.zip!inner\.zip!src/vendor/`)},
	})
	findings := scanDirWith(t, NewDetector(*cfg), cfg, root, 3)
	assert.Empty(t, findings)
}

func TestDirRootUnicodeAndSpaces(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "sub dir", "秘密"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "sub dir", "秘密", "x y.py"),
		[]byte("Token = \""+dirRootSecret+"\"\n"), 0o600))

	cfg := dirRootConfig(&config.Allowlist{
		Paths: []*regexp.Regexp{regexp.MustCompile(`^sub dir/秘密/`)},
	})
	assert.Empty(t, scanDirRoot(t, cfg, root))
}

func TestDirRootInformativeExclusionOnce(t *testing.T) {
	root := dirRootTree(t)
	cfg := dirRootConfig(&config.Allowlist{
		Paths: []*regexp.Regexp{regexp.MustCompile(`^src/vendor($|/)`)},
	})
	collector := informative.New("directory")
	detector := NewDetector(*cfg)
	detector.Informative = collector

	findings := scanDirWith(t, detector, cfg, root, 0)
	assert.Equal(t, []string{"src/app.py", "src/nested/deep.py"}, dirRelativeNames(t, root, findings))

	doc := collector.Document("", true)
	require.Len(t, doc.Exclusions, 1)
	exclusion := doc.Exclusions[0]
	assert.Equal(t, informative.ReasonPathAllowlist, exclusion.Reason)
	assert.Equal(t, informative.KindDirectory, exclusion.Kind)
	assert.Equal(t, filepath.ToSlash(filepath.Join(root, "src", "vendor")), exclusion.Path)
	assert.Equal(t, 1, doc.Summary.ExclusionsByReason[informative.ReasonPathAllowlist])
}

func TestDirRootInformativeConfigExcludedInput(t *testing.T) {
	root := dirRootTree(t)
	cfg := dirRootConfig()
	cfg.Path = filepath.Join(root, ".gitleaks.toml")
	require.NoError(t, os.WriteFile(cfg.Path,
		[]byte("# synthetic credential in config: RNC_TEST_SECRET_cccccccccccccccc\n"), 0o600))
	collector := informative.New("directory")
	detector := NewDetector(*cfg)
	detector.Informative = collector

	findings := scanDirWith(t, detector, cfg, root, 0)
	assert.Len(t, findings, 3)
	doc := collector.Document("", true)
	require.Len(t, doc.Exclusions, 1)
	assert.Equal(t, informative.ReasonExcludedInput, doc.Exclusions[0].Reason)
	assert.Equal(t, filepath.ToSlash(filepath.Join(root, ".gitleaks.toml")), doc.Exclusions[0].Path)
}

func writeZip(t *testing.T, path string, files map[string]string) {
	t.Helper()
	data := zipStringBytes(t, files)
	require.NoError(t, os.WriteFile(path, data, 0o600))
}

func writeZipBytes(t *testing.T, path string, files map[string][]byte) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, zipBytes(t, files), 0o600))
}

func zipBytes(t *testing.T, files map[string][]byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	keys := make([]string, 0, len(files))
	for name := range files {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	for _, name := range keys {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write(files[name])
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func zipStringBytes(t *testing.T, files map[string]string) []byte {
	t.Helper()
	converted := make(map[string][]byte, len(files))
	for name, content := range files {
		converted[name] = []byte(content)
	}
	return zipBytes(t, converted)
}
