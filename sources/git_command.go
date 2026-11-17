package sources

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/gitleaks/go-gitdiff/gitdiff"

	"github.com/zricethezav/gitleaks/v8/logging"
)

// forceNoColor makes --no-color the effective color setting without treating
// the option as a pathspec. A later --no-color wins over user-provided
// --color/--color=... options; when the user supplied a "--" separator, it is
// inserted immediately before the separator rather than after it.
// forceGitOutputFlags pins the machine-readable Git presentation used by the
// scanner. Repository-local format.pretty, log.date, color and textconv
// settings must not change the parser's input or execute repository-controlled
// helpers. The flags are inserted before a user pathspec separator.
func forceGitOutputFlags(args []string) []string {
	args = forceNoColor(args)
	extra := []string{"--no-ext-diff", "--no-textconv", "--date=default", "--pretty=medium"}
	separator := len(args)
	for i, arg := range args {
		if arg == "--" {
			separator = i
			break
		}
	}
	out := make([]string, 0, len(args)+len(extra))
	out = append(out, args[:separator]...)
	out = append(out, extra...)
	out = append(out, args[separator:]...)
	return out
}

func forceNoColor(args []string) []string {
	separator := len(args)
	for i, arg := range args {
		if arg == "--" {
			separator = i
			break
		}
	}

	for _, arg := range args[:separator] {
		if arg == "--color" || strings.HasPrefix(arg, "--color=") {
			logging.Warn().Msg("overriding user-supplied --color option; machine-readable diff parsing requires --no-color")
			break
		}
	}

	if separator < len(args) {
		out := make([]string, 0, len(args)+1)
		out = append(out, args[:separator]...)
		out = append(out, "--no-color")
		out = append(out, args[separator:]...)
		return out
	}
	return append(args, "--no-color")
}

// GitCmd helps to work with Git's output.
type GitCmd struct {
	cmd         *exec.Cmd
	stdout      *bufio.Reader
	diffFilesCh <-chan *gitdiff.File
	errCh       <-chan error
	repoPath    string
	stderr      *gitStderr
	tracker     *gitOutputTracker
	// staged is true for `git diff --staged`; binary-marked text must then be
	// read from the index rather than skipped or read from the worktree.
	staged bool
	// contextLines is the number of unchanged context lines included in each
	// history patch. Zero is the addition-only default.
	contextLines int

	drainOnce sync.Once
	unparsed  int64
	waitOnce  sync.Once
	waitErr   error
}

// gitOutputTracker counts complete "diff --git " section headers in the raw
// output. The diff parser silently stops on some malformed or truncated input,
// so comparing raw section count with parsed file count is the only way to
// notice that the tail of the stream was never scanned.
type gitOutputTracker struct {
	mu        sync.Mutex
	prefix    [16]byte
	prefixLen int
	sections  int64
}

func (t *gitOutputTracker) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	written := len(p)
	// Count only line-leading prefixes. IndexByte is optimized for the
	// newline search, avoiding a per-byte Go loop over every diff byte.
	for len(p) > 0 {
		i := bytes.IndexByte(p, '\n')
		if i < 0 {
			t.appendPrefix(p)
			break
		}
		t.appendPrefix(p[:i])
		t.countLine()
		p = p[i+1:]
	}
	return written, nil
}

func (t *gitOutputTracker) appendPrefix(line []byte) {
	if t.prefixLen >= len(t.prefix) {
		return
	}
	n := copy(t.prefix[t.prefixLen:], line)
	t.prefixLen += n
}

func (t *gitOutputTracker) finish() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.countLine()
}

func (t *gitOutputTracker) countLine() {
	if t.prefixLen >= len("diff --git ") &&
		string(t.prefix[:len("diff --git ")]) == "diff --git " {
		t.sections++
	}
	t.prefixLen = 0
}

func (t *gitOutputTracker) sectionCount() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.sections
}

const maxGitStderrLines = 20

var (
	gitURLUserinfoPattern = regexp.MustCompile(`(?i)([a-z][a-z0-9+.-]*://)[^/@\s]+@`)
	gitSensitiveQuery     = regexp.MustCompile(`(?i)([?&](?:access_token|token|password|passwd|secret|api[_-]?key|apikey|auth|key)=)[^&\s'":,)\]]+`)
	gitCredentialToken    = regexp.MustCompile(`[A-Za-z0-9_+/=-]{16,}`)
)

// sanitizeGitDiagnostic keeps Git's diagnostic shape while removing URL
// userinfo, sensitive query values, and token-shaped runs. Git stderr is not a
// credential channel, but remotes and authentication helpers can place secrets
// there; raw stderr must not become an additional disclosure surface.
func sanitizeGitDiagnostic(line string) string {
	line = gitURLUserinfoPattern.ReplaceAllString(line, "${1}REDACTED@")
	line = gitSensitiveQuery.ReplaceAllString(line, "${1}REDACTED")
	line = gitCredentialToken.ReplaceAllString(line, "REDACTED")
	return line
}

// gitStderr records bounded stderr output and whether it contained an explicit
// Git error. Nonempty stderr alone is not a failure: Git emits warnings and
// progress messages there during successful commands.
type gitStderr struct {
	mu      sync.Mutex
	lines   []string
	total   int
	fatal   bool
	readErr error
}

func (g *gitStderr) add(line string, fatal bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.total++
	if len(g.lines) < maxGitStderrLines {
		g.lines = append(g.lines, sanitizeGitDiagnostic(line))
	}
	if fatal {
		g.fatal = true
	}
}

func (g *gitStderr) setReadErr(err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.readErr = err
	g.fatal = true
}

func (g *gitStderr) summary() string {
	if g == nil {
		return ""
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	parts := append([]string(nil), g.lines...)
	if g.total > len(g.lines) {
		parts = append(parts, fmt.Sprintf("%d additional stderr line(s) omitted", g.total-len(g.lines)))
	}
	if g.readErr != nil {
		parts = append(parts, fmt.Sprintf("could not read stderr: %v", g.readErr))
	}
	return strings.Join(parts, "; ")
}

func (g *gitStderr) err() error {
	if g == nil {
		return nil
	}
	g.mu.Lock()
	fatal := g.fatal
	readErr := g.readErr
	g.mu.Unlock()
	if !fatal && readErr == nil {
		return nil
	}
	if summary := g.summary(); summary != "" {
		return fmt.Errorf("git reported an error on stderr: %s", summary)
	}
	return errors.New("git reported an error on stderr")
}

func isFatalGitStderr(line string) bool {
	lower := strings.ToLower(strings.TrimSpace(line))
	return strings.HasPrefix(lower, "fatal:") ||
		strings.HasPrefix(lower, "error:") ||
		strings.HasPrefix(lower, "panic:")
}

// blobReader provides a ReadCloser interface git cat-file blob to fetch
// a blob from a repo
type blobReader struct {
	io.ReadCloser
	cmd *exec.Cmd
	eof bool
}

// Read records complete exhaustion so callers that intentionally stop early
// (for example after MIME classification) can terminate the producer instead
// of draining the rest of a large object.
func (br *blobReader) Read(p []byte) (int, error) {
	n, err := br.ReadCloser.Read(p)
	if err == io.EOF {
		br.eof = true
	}
	return n, err
}

// CloseAfterSkip closes a reader whose caller has deliberately finished before
// EOF. It terminates the producer rather than draining an arbitrarily large
// remaining object, then waits for the process. A successfully issued kill is
// intentional and its signal status is not a scan failure.
func (br *blobReader) CloseAfterSkip() error {
	if br.eof || br.cmd == nil || br.cmd.Process == nil {
		return br.Close()
	}

	killErr := br.cmd.Process.Kill()
	closeErr := br.ReadCloser.Close()
	waitErr := br.cmd.Wait()
	if closeErr != nil {
		return closeErr
	}
	if killErr == nil {
		return nil
	}
	return waitErr
}

// Close closes the underlying reader and then waits for the command to complete,
// releasing its resources.
func (br *blobReader) Close() error {
	// Discard the remaining data from the pipe to avoid blocking
	_, drainErr := io.Copy(io.Discard, br)
	// Close the pipe (should signal the command to stop if it hasn't already)
	closeErr := br.ReadCloser.Close()
	// Wait to prevent zombie processes.
	waitErr := br.cmd.Wait()
	// Return the first error encountered
	if drainErr != nil {
		return drainErr
	}
	if closeErr != nil {
		return closeErr
	}
	return waitErr
}

// NewGitLogCmd returns `*DiffFilesCmd` with two channels: `<-chan *gitdiff.File` and `<-chan error`.
// Caller should read everything from channels until receiving a signal about their closure and call
// the `func (*DiffFilesCmd) Wait()` error in order to release resources.
func NewGitLogCmd(source string, logOpts string) (*GitCmd, error) {
	return NewGitLogCmdContext(context.Background(), source, logOpts)
}

// NewGitLogCmdContext is the same as NewGitLogCmd but supports passing in a
// context to use for timeouts
func NewGitLogCmdContext(ctx context.Context, source string, logOpts string) (*GitCmd, error) {
	return NewGitLogCmdContextWithOptions(ctx, source, logOpts, GitLogOptions{})
}

// GitLogOptions controls optional history-selection behavior without changing
// the established default. ScanMerges asks Git for a merge commit's diff
// against its first parent, which recovers secrets introduced only by a
// conflict resolution or evil merge at the cost of reporting branch content a
// second time at the merge commit.
type GitLogOptions struct {
	ScanMerges bool
	// ContextLines adds unchanged context around each patch and includes it in
	// matching fragments. Findings are still accepted only when their line
	// span overlaps an added line. The default zero preserves addition-only
	// scanning.
	ContextLines int
}

const maxGitContextLines = 100

func (o GitLogOptions) validate() error {
	if o.ContextLines < 0 || o.ContextLines > maxGitContextLines {
		return fmt.Errorf("git context lines must be between 0 and %d", maxGitContextLines)
	}
	return nil
}

// NewGitLogCmdContextWithOptions is NewGitLogCmdContext with explicit optional
// history behavior.
func NewGitLogCmdContextWithOptions(ctx context.Context, source string, logOpts string, options GitLogOptions) (*GitCmd, error) {
	userArgs, err := parseGitLogOpts(logOpts)
	if err != nil {
		return nil, err
	}
	if err := validateGitLogArgs(userArgs); err != nil {
		return nil, err
	}
	if err := options.validate(); err != nil {
		return nil, err
	}
	sourceClean := filepath.Clean(source)
	shallow, bare, top, topOK := gitRepositoryState(ctx, sourceClean)
	if shallow {
		logging.Warn().Msg("source is a shallow Git repository; history before the shallow boundary is unavailable and cannot be treated as a complete full-history scan")
	}
	if bare {
		logging.Warn().Msg("source is a bare Git repository; automatic working-tree .gitleaks.toml/.gitleaksignore discovery is unavailable, so pass explicit --config/--gitleaks-ignore-path when needed")
	} else if topOK && !sameFilesystemPath(top, sourceClean) {
		logging.Warn().Msg("Git history is repository-wide; the selected source is a subdirectory. Use --log-opts with a pathspec to limit history. Directory scans remain bounded to the selected path")
	}
	args := []string{"-C", sourceClean, "log", "-p", fmt.Sprintf("-U%d", options.ContextLines)}
	if options.ScanMerges {
		// Placed before user arguments so an explicit --no-diff-merges (or a
		// supported --diff-merges mode) remains authoritative.
		args = append(args, "--diff-merges=first-parent")
	}
	if len(userArgs) > 0 {
		args = append(args, userArgs...)
	} else {
		args = append(args, "--full-history", "--all", "--diff-filter=tuxdb")
	}
	args = forceGitOutputFlags(args)
	cmd := exec.CommandContext(ctx, "git", args...)

	logging.Debug().Msgf("executing: %s", cmd.String())

	gitCmd, err := newGitCmd(cmd, sourceClean)
	if err != nil {
		return nil, err
	}
	gitCmd.contextLines = options.ContextLines
	return gitCmd, nil
}

// parseGitLogOpts tokenizes --log-opts without a shell. It supports single and
// double quoting, escaped whitespace/quotes, and literal Windows backslashes.
// The parsed vector is authoritative for both execution and format validation.
func parseGitLogOpts(raw string) ([]string, error) {
	if raw == "" {
		return nil, nil
	}
	var (
		args               []string
		token              strings.Builder
		inSingle, inDouble bool
		escaped            bool
		hasToken           bool
		runes              = []rune(raw)
	)
	for i, r := range runes {
		if escaped {
			token.WriteRune(r)
			escaped, hasToken = false, true
			continue
		}
		if r == '\\' && !inSingle {
			nextEscapes := i+1 < len(runes) && (runes[i+1] == '\\' || runes[i+1] == '\'' || runes[i+1] == '"' || unicode.IsSpace(runes[i+1]))
			if nextEscapes {
				escaped = true
				continue
			}
			token.WriteRune(r)
			hasToken = true
			continue
		}
		if inSingle {
			if r == '\'' {
				inSingle = false
			} else {
				token.WriteRune(r)
			}
			hasToken = true
			continue
		}
		if inDouble {
			if r == '"' {
				inDouble = false
			} else {
				token.WriteRune(r)
			}
			hasToken = true
			continue
		}
		switch {
		case unicode.IsSpace(r):
			if hasToken {
				if token.Len() == 0 {
					return nil, fmt.Errorf("empty --log-opts argument")
				}
				args = append(args, token.String())
				token.Reset()
				hasToken = false
			}
		case r == '\'':
			inSingle, hasToken = true, true
		case r == '"':
			inDouble, hasToken = true, true
		default:
			token.WriteRune(r)
			hasToken = true
		}
	}
	if escaped {
		return nil, fmt.Errorf("unmatched escape in --log-opts")
	}
	if inSingle || inDouble {
		return nil, fmt.Errorf("unmatched quote in --log-opts")
	}
	if hasToken {
		if token.Len() == 0 {
			return nil, fmt.Errorf("empty --log-opts argument")
		}
		args = append(args, token.String())
	}
	return args, nil
}

func gitTopLevel(ctx context.Context, source string) (string, bool) {
	out, err := exec.CommandContext(ctx, "git", "-C", source, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", false
	}
	top := strings.TrimSpace(string(out))
	return top, top != ""
}

func isBareRepository(ctx context.Context, source string) bool {
	out, err := exec.CommandContext(ctx, "git", "-C", source, "rev-parse", "--is-bare-repository").Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// gitRepositoryState collects the shallow/bare/top-level facts used for
// provenance warnings. The normal non-bare case is answered by one rev-parse
// call. Git rejects --show-toplevel in a bare repository, so that case falls
// back to a two-state query instead of treating the warning check as a
// failure. This replaces three sequential subprocesses on the normal path,
// which is material process-startup overhead for short scans.
func gitRepositoryState(ctx context.Context, source string) (shallow, bare bool, top string, topOK bool) {
	// The common non-bare case is answered in one subprocess. Git rejects
	// --show-toplevel in a bare repository, so fall back to the two-state
	// query there rather than treating the warning check as a failure.
	out, err := exec.CommandContext(
		ctx,
		"git", "-C", source,
		"rev-parse", "--is-shallow-repository", "--is-bare-repository", "--show-toplevel",
	).Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 3 {
			return strings.TrimSpace(lines[0]) == "true",
				strings.TrimSpace(lines[1]) == "true",
				strings.TrimSpace(lines[2]),
				true
		}
	}

	out, err = exec.CommandContext(ctx, "git", "-C", source, "rev-parse", "--is-shallow-repository", "--is-bare-repository").Output()
	if err != nil {
		return false, false, "", false
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 2 {
		return false, false, "", false
	}
	return strings.TrimSpace(lines[0]) == "true", strings.TrimSpace(lines[1]) == "true", "", false
}

func sameFilesystemPath(a, b string) bool {
	absA, errA := filepath.Abs(a)
	absB, errB := filepath.Abs(b)
	if errA != nil || errB != nil {
		return false
	}
	if filepath.Clean(absA) == filepath.Clean(absB) {
		return true
	}
	infoA, errA := os.Stat(absA)
	infoB, errB := os.Stat(absB)
	return errA == nil && errB == nil && os.SameFile(infoA, infoB)
}

// isShallowRepository reports whether Git identifies the source as shallow.
// Errors are treated as unknown rather than scanning-fatal; history selection
// and finding identity remain unchanged.
func isShallowRepository(ctx context.Context, source string) bool {
	cmd := exec.CommandContext(ctx, "git", "-C", source, "rev-parse", "--is-shallow-repository")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// validateGitLogOpts is the string boundary used by CLI validation. It parses
// first, then validates the exact argument vector that will be executed.
func validateGitLogOpts(logOpts string) error {
	args, err := parseGitLogOpts(logOpts)
	if err != nil {
		return err
	}
	return validateGitLogArgs(args)
}

// validateGitLogArgs rejects formats that suppress patch output or emit
// combined/remerge diffs the parser cannot complete. Validation operates on
// the parsed vector so quoting cannot hide a disallowed option.
func validateGitLogArgs(args []string) error {
	afterSeparator := false
	for i, option := range args {
		if option == "--" {
			afterSeparator = true
			continue
		}
		if afterSeparator {
			continue
		}
		if option == "--diff-merges" && i+1 < len(args) {
			switch args[i+1] {
			case "cc", "combined", "dense-combined", "remerge":
				return fmt.Errorf("unsupported --log-opts value %q: the scanner requires ordinary patch output (use --diff-merges=first-parent for merge-only changes)", option)
			}
		}
		switch {
		case option == "--output",
			strings.HasPrefix(option, "--output="),
			option == "-o":
			return fmt.Errorf("unsupported --log-opts value %q: Git output redirection would bypass the scanner's stdout and could produce a false-clean result", option)
		case option == "--name-only",
			option == "--name-status",
			option == "--raw",
			option == "--numstat",
			option == "--shortstat",
			option == "--no-patch",
			option == "-s",
			option == "--remerge-diff",
			option == "--diff-merges=cc",
			option == "--diff-merges=combined",
			option == "--diff-merges=dense-combined",
			option == "--diff-merges=remerge",
			option == "--cc",
			option == "-c",
			strings.HasPrefix(option, "-c") && len(option) > 2:
			return fmt.Errorf("unsupported --log-opts value %q: the scanner requires ordinary patch output (use --diff-merges=first-parent for merge-only changes)", option)
		}
	}
	return nil
}

// NewGitDiffCmd returns `*DiffFilesCmd` with two channels: `<-chan *gitdiff.File` and `<-chan error`.
// Caller should read everything from channels until receiving a signal about their closure and call
// the `func (*DiffFilesCmd) Wait()` error in order to release resources.
func NewGitDiffCmd(source string, staged bool) (*GitCmd, error) {
	return NewGitDiffCmdContext(context.Background(), source, staged)
}

// NewGitDiffCmdContext is the same as NewGitDiffCmd but supports passing in a
// context to use for timeouts
func NewGitDiffCmdContext(ctx context.Context, source string, staged bool) (*GitCmd, error) {
	sourceClean := filepath.Clean(source)
	var cmd *exec.Cmd
	cmd = exec.CommandContext(ctx, "git", "-C", sourceClean, "diff", "-U0", "--no-ext-diff", "--no-textconv", "--no-color", ".")
	if staged {
		cmd = exec.CommandContext(ctx, "git", "-C", sourceClean, "diff", "-U0", "--no-ext-diff", "--no-textconv",
			"--no-color", "--staged", ".")
	}
	logging.Debug().Msgf("executing: %s", cmd.String())

	gitCmd, err := newGitCmd(cmd, sourceClean)
	if err != nil {
		return nil, err
	}
	gitCmd.staged = staged
	return gitCmd, nil
}

// newGitCmd starts cmd, connects its output to the diff parser, and records
// stderr separately from the parsed diff stream.
func newGitCmd(cmd *exec.Cmd, repoPath string) (*GitCmd, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	stderrInfo := &gitStderr{}
	errCh := make(chan error, 1)
	go listenForStdErr(stderr, stderrInfo, errCh)

	tracker := &gitOutputTracker{}
	// A larger read buffer reduces pipe-read syscalls for large histories.
	buffered := bufio.NewReaderSize(io.TeeReader(stdout, tracker), 64*1024)
	diffFiles, err := gitdiff.Parse(buffered)
	if err != nil {
		// The command is running but its output never reached the parser.
		// Kill it, drain both pipes, and reap it before returning the error.
		_ = cmd.Process.Kill()
		_, _ = io.Copy(io.Discard, buffered)
		for range errCh {
		}
		_ = cmd.Wait()
		return nil, fmt.Errorf("could not parse git diff: %w", err)
	}

	return &GitCmd{
		cmd:         cmd,
		stdout:      buffered,
		diffFilesCh: diffFiles,
		errCh:       errCh,
		repoPath:    repoPath,
		stderr:      stderrInfo,
		tracker:     tracker,
	}, nil
}

// DiffFilesCh returns a channel with *gitdiff.File.
func (c *GitCmd) DiffFilesCh() <-chan *gitdiff.File {
	return c.diffFilesCh
}

// ErrCh returns a channel that produces an error when Git reports an explicit
// fatal/error diagnostic on stderr. Warnings and progress messages do not
// produce errors; process exit status is checked by Wait.
func (c *GitCmd) ErrCh() <-chan error {
	return c.errCh
}

// Wait drains any stdout left by the parser, waits for the command to exit, and
// returns the command's exit error. It is safe to call more than once.
func (c *GitCmd) Wait() error {
	c.drain()
	c.waitOnce.Do(func() {
		c.waitErr = c.cmd.Wait()
	})
	return c.waitErr
}

func (c *GitCmd) drain() {
	c.drainOnce.Do(func() {
		if c.stdout == nil {
			return
		}
		c.unparsed, _ = io.Copy(io.Discard, c.stdout)
		if c.tracker != nil {
			c.tracker.finish()
		}
	})
}

// unparsedBytes reports how much Git output remained after the parser closed
// its output channel. It must be called after DiffFilesCh is closed.
func (c *GitCmd) unparsedBytes() int64 {
	c.drain()
	return c.unparsed
}

func (c *GitCmd) sectionCount() int64 {
	c.drain()
	if c.tracker == nil {
		return 0
	}
	return c.tracker.sectionCount()
}

func (c *GitCmd) commandError(err error) error {
	if summary := c.stderr.summary(); summary != "" {
		return fmt.Errorf("git command failed: %w; stderr: %s", err, summary)
	}
	return fmt.Errorf("git command failed: %w", err)
}

// String displays the command used for GitCmd
func (c *GitCmd) String() string {
	return c.cmd.String()
}

// NewBlobReader returns an io.ReadCloser that can be used to read a blob
// within the git repo used to create the GitCmd.
//
// The caller is responsible for closing the reader.
func (c *GitCmd) NewBlobReader(commit, path string) (io.ReadCloser, error) {
	return c.NewBlobReaderContext(context.Background(), commit, path)
}

// NewBlobReaderContext is the same as NewBlobReader but supports passing in a
// context to use for timeouts
func (c *GitCmd) NewBlobReaderContext(ctx context.Context, commit, path string) (io.ReadCloser, error) {
	return c.newBlobReaderContext(ctx, commit+":"+path)
}

// NewIndexBlobReaderContext returns the staged index content for path. The
// explicit stage-zero prefix uses the index rather than the worktree.
func (c *GitCmd) NewIndexBlobReaderContext(ctx context.Context, path string) (io.ReadCloser, error) {
	return c.newBlobReaderContext(ctx, ":0:"+path)
}

func (c *GitCmd) newBlobReaderContext(ctx context.Context, object string) (io.ReadCloser, error) {
	gitArgs := []string{"-C", c.repoPath, "cat-file", "blob", object}
	cmd := exec.CommandContext(ctx, "git", gitArgs...)
	cmd.Stderr = io.Discard
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start git command: %w", err)
	}
	return &blobReader{
		ReadCloser: stdout,
		cmd:        cmd,
	}, nil
}

// listenForStdErr records and logs Git's stderr. Known warnings and progress
// messages are not failures; explicit fatal/error diagnostics are signalled on
// errCh and rechecked against the process exit status by Git.Fragments.
func listenForStdErr(stderr io.ReadCloser, info *gitStderr, errCh chan<- error) {
	defer close(errCh)

	var signalled bool
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		fatal := isFatalGitStderr(line)
		info.add(line, fatal)
		if fatal {
			logging.Error().Msgf("[git] %s", sanitizeGitDiagnostic(line))
			if !signalled {
				select {
				case errCh <- errors.New("git reported an error on stderr"):
					signalled = true
				default:
				}
			}
			continue
		}
		logging.Warn().Msgf("[git] %s", sanitizeGitDiagnostic(line))
	}

	if err := scanner.Err(); err != nil {
		info.setReadErr(err)
		select {
		case errCh <- fmt.Errorf("could not read git stderr: %w", err):
		default:
		}
	}
}
