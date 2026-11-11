package sources

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fatih/semgroup"
	"github.com/gitleaks/go-gitdiff/gitdiff"

	"github.com/zricethezav/gitleaks/v8/cmd/scm"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/internal/scanerr"
	"github.com/zricethezav/gitleaks/v8/logging"
)

var quotedOptPattern = regexp.MustCompile(`^(?:"[^"]+"|'[^']+')$`)

// forceNoColor makes --no-color the effective color setting without treating
// the option as a pathspec. A later --no-color wins over user-provided
// --color/--color=... options; when the user supplied a "--" separator, it is
// inserted immediately before the separator rather than after it.
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
	for _, b := range p {
		if b == '\n' {
			t.countLine()
			continue
		}
		if t.prefixLen < len(t.prefix) {
			t.prefix[t.prefixLen] = b
			t.prefixLen++
		}
	}
	return len(p), nil
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
	sourceClean := filepath.Clean(source)
	var cmd *exec.Cmd
	if logOpts != "" {
		args := []string{"-C", sourceClean, "log", "-p", "-U0", "--no-color"}

		// Ensure that the user-provided |logOpts| aren't wrapped in quotes.
		// https://github.com/gitleaks/gitleaks/issues/1153
		userArgs := strings.Split(logOpts, " ")
		var quotedOpts []string
		for _, element := range userArgs {
			if quotedOptPattern.MatchString(element) {
				quotedOpts = append(quotedOpts, element)
			}
		}
		if len(quotedOpts) > 0 {
			logging.Warn().Msgf("the following `--log-opts` values may not work as expected: %v\n\tsee https://github.com/gitleaks/gitleaks/issues/1153 for more information", quotedOpts)
		}

		args = append(args, userArgs...)
		args = forceNoColor(args)
		cmd = exec.CommandContext(ctx, "git", args...)
	} else {
		cmd = exec.CommandContext(ctx, "git", "-C", sourceClean, "log", "-p", "-U0", "--no-color",
			"--full-history", "--all", "--diff-filter=tuxdb")
	}

	logging.Debug().Msgf("executing: %s", cmd.String())

	return newGitCmd(cmd, sourceClean)
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
	cmd = exec.CommandContext(ctx, "git", "-C", sourceClean, "diff", "-U0", "--no-ext-diff", "--no-color", ".")
	if staged {
		cmd = exec.CommandContext(ctx, "git", "-C", sourceClean, "diff", "-U0", "--no-ext-diff",
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
	buffered := bufio.NewReader(io.TeeReader(stdout, tracker))
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

// RemoteInfo provides the info needed for reconstructing links from findings
type RemoteInfo struct {
	Platform scm.Platform
	Url      string
}

// Git is a source for yielding fragments from a git repo
type Git struct {
	Cmd             *GitCmd
	Config          *config.Config
	Remote          *RemoteInfo
	Sema            *semgroup.Group
	MaxArchiveDepth int
}

// CommitInfo captures metadata about the commit
type CommitInfo struct {
	AuthorEmail string
	AuthorName  string
	Date        string
	Message     string
	Remote      *RemoteInfo
	SHA         string
}

// Fragments yields fragments from a git repo
func (s *Git) Fragments(ctx context.Context, yield FragmentsFunc) error {
	if s.Cmd == nil {
		return errors.New("git source has no command")
	}
	if s.Sema == nil {
		s.Sema = semgroup.NewGroup(ctx, 40)
	}

	collector := scanerr.New(10)
	diffFilesCh := s.Cmd.DiffFilesCh()
	errCh := s.Cmd.ErrCh()
	ctxDone := ctx.Done()

	var (
		ctxErr      error
		parseErr    error
		parsedDiffs int64
	)

	// Drain both Git channels before waiting for the process. Returning early
	// would leave the parser and stderr goroutines blocked on their pipes.
	for diffFilesCh != nil || errCh != nil {
		select {
		case <-ctxDone:
			ctxErr = ctx.Err()
			ctxDone = nil
		case gitdiffFile, open := <-diffFilesCh:
			if !open {
				diffFilesCh = nil
				if unparsed := s.Cmd.unparsedBytes(); unparsed > 0 {
					parseErr = fmt.Errorf(
						"diff parsing stopped early, leaving %d bytes of git output unscanned: results are incomplete",
						unparsed,
					)
				}
				if expected := s.Cmd.sectionCount(); expected > parsedDiffs {
					sectionErr := fmt.Errorf(
						"diff parsing stopped early: parsed %d of %d diff sections; results are incomplete",
						parsedDiffs, expected,
					)
					parseErr = errors.Join(parseErr, sectionErr)
				}
				continue
			}
			if gitdiffFile == nil {
				continue
			}
			parsedDiffs++

			if ctxErr != nil {
				continue
			}
			if gitdiffFile.IsDelete {
				continue
			}

			// Git marks staged text as binary when attributes or
			// core.bigFileThreshold prevent a textual diff. Staged scans must
			// fall back to the index blob; history scans retain the established
			// skip for non-archive binary files. File.Fragments still applies
			// MIME classification to the blob content. The staged-index approach
			// follows gitleaks/gitleaks#2338 (@VISHNU0906); this implementation
			// uses the CP-V3 error and lifecycle path.
			yieldFromBlob := false
			if gitdiffFile.IsBinary {
				if isArchive(ctx, gitdiffFile.NewName) || s.Cmd.staged {
					yieldFromBlob = true
				} else {
					continue
				}
			}

			// Check if commit is allowed
			commitSHA := ""
			var commitInfo *CommitInfo
			if gitdiffFile.PatchHeader != nil {
				commitSHA = gitdiffFile.PatchHeader.SHA
				for _, a := range s.Config.Allowlists {
					if ok, c := a.CommitAllowed(gitdiffFile.PatchHeader.SHA); ok {
						logging.Trace().Str("allowed-commit", c).Msg("skipping commit: global allowlist")
						continue
					}
				}

				commitInfo = &CommitInfo{
					Date:    gitdiffFile.PatchHeader.AuthorDate.UTC().Format(time.RFC3339),
					Message: gitdiffFile.PatchHeader.Message(),
					Remote:  s.Remote,
					SHA:     commitSHA,
				}

				if gitdiffFile.PatchHeader.Author != nil {
					commitInfo.AuthorName = gitdiffFile.PatchHeader.Author.Name
					commitInfo.AuthorEmail = gitdiffFile.PatchHeader.Author.Email
				}
			}

			diffFile := gitdiffFile
			fromBlob := yieldFromBlob
			sha := commitSHA
			info := commitInfo
			s.Sema.Go(func() error {
				if err := s.processDiffFile(ctx, diffFile, fromBlob, sha, info, yield); err != nil {
					collector.Add(fmt.Errorf("%s: %w", diffFile.NewName, err))
				}
				return nil
			})
		case _, open := <-errCh:
			if !open {
				errCh = nil
			}
		}
	}

	// Join workers before checking process completion so findings from readable
	// input are complete and no worker can run after this function returns.
	collector.Add(s.Sema.Wait())

	waitErr := s.Cmd.Wait()
	if waitErr != nil {
		collector.Add(s.Cmd.commandError(waitErr))
	} else {
		// A fatal stderr diagnostic with a zero exit status is still a failure.
		collector.Add(s.Cmd.stderr.err())
	}
	collector.Add(parseErr)
	collector.Add(ctxErr)

	return collector.Err()
}

func (s *Git) processDiffFile(
	ctx context.Context,
	diffFile *gitdiff.File,
	fromBlob bool,
	commitSHA string,
	commitInfo *CommitInfo,
	yield FragmentsFunc,
) error {
	if fromBlob {
		var (
			blob io.ReadCloser
			err  error
		)
		if s.Cmd.staged {
			blob, err = s.Cmd.NewIndexBlobReaderContext(ctx, diffFile.NewName)
		} else {
			blob, err = s.Cmd.NewBlobReaderContext(ctx, commitSHA, diffFile.NewName)
		}
		if err != nil {
			return fmt.Errorf("could not read blob: %w", err)
		}

		file := File{
			Content:         blob,
			Path:            diffFile.NewName,
			MaxArchiveDepth: s.MaxArchiveDepth,
			Config:          s.Config,
		}

		// enrich and yield fragments
		fragmentsErr := file.Fragments(ctx, func(fragment Fragment, err error) error {
			fragment.CommitSHA = commitSHA
			fragment.CommitInfo = commitInfo
			return yield(fragment, err)
		})

		var closeErr error
		if br, ok := blob.(*blobReader); ok && fragmentsErr == nil {
			// File.Fragments intentionally stops after MIME exclusion (or an
			// archive-depth skip) without consuming the rest of the object.
			// Terminate the producer instead of reading every remaining byte.
			closeErr = br.CloseAfterSkip()
		} else {
			closeErr = blob.Close()
		}
		if fragmentsErr != nil || closeErr != nil {
			return fmt.Errorf("could not read blob content: %w", errors.Join(fragmentsErr, closeErr))
		}
		return nil
	}

	for _, textFragment := range diffFile.TextFragments {
		if textFragment == nil {
			continue
		}

		fragment := Fragment{
			CommitSHA:  commitSHA,
			FilePath:   diffFile.NewName,
			Raw:        textFragment.Raw(gitdiff.OpAdd),
			StartLine:  int(textFragment.NewPosition),
			CommitInfo: commitInfo,
		}

		if err := yield(fragment, nil); err != nil {
			return err
		}
	}

	return nil
}

// NewRemoteInfo builds a new RemoteInfo for generating finding links
func NewRemoteInfo(platform scm.Platform, source string) *RemoteInfo {
	return NewRemoteInfoContext(context.Background(), platform, source)
}

// NewRemoteInfoContext is the same as NewRemoteInfo but supports passing in a
// context to use for timeouts
func NewRemoteInfoContext(ctx context.Context, platform scm.Platform, source string) *RemoteInfo {
	if platform == scm.NoPlatform {
		return &RemoteInfo{Platform: platform}
	}

	remoteUrl, err := getRemoteUrl(ctx, source)
	if err != nil {
		if strings.Contains(err.Error(), "No remote configured") {
			logging.Debug().Msg("skipping finding links: repository has no configured remote.")
			platform = scm.NoPlatform
		} else {
			logging.Error().Err(err).Msg("skipping finding links: unable to parse remote URL")
		}
		goto End
	}

	if platform == scm.UnknownPlatform {
		platform = platformFromHost(remoteUrl)
		if platform == scm.UnknownPlatform {
			logging.Info().
				Str("host", remoteUrl.Hostname()).
				Msg("Unknown SCM platform. Use --platform to include links in findings.")
		} else {
			logging.Debug().
				Str("host", remoteUrl.Hostname()).
				Str("platform", platform.String()).
				Msg("SCM platform parsed from host")
		}
	}

End:
	var rUrl string
	if remoteUrl != nil {
		rUrl = remoteUrl.String()
	}
	return &RemoteInfo{
		Platform: platform,
		Url:      rUrl,
	}
}

var sshUrlpat = regexp.MustCompile(`^git@([a-zA-Z0-9.-]+):(?:\d{1,5}/)?([\w/.-]+?)(?:\.git)?$`)

func getRemoteUrl(ctx context.Context, source string) (*url.URL, error) {
	// This will return the first remote — typically, "origin".
	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--quiet", "--get-url")
	if source != "." {
		cmd.Dir = source
	}

	stdout, err := cmd.Output()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return nil, fmt.Errorf("command failed (%d): %w, stderr: %s", exitError.ExitCode(), err,
				sanitizeGitDiagnostic(string(bytes.TrimSpace(exitError.Stderr))))
		}
		return nil, err
	}

	remoteUrl := string(bytes.TrimSpace(stdout))
	if matches := sshUrlpat.FindStringSubmatch(remoteUrl); matches != nil {
		remoteUrl = fmt.Sprintf("https://%s/%s", matches[1], matches[2])
	}
	remoteUrl = strings.TrimSuffix(remoteUrl, ".git")

	parsedUrl, err := url.Parse(remoteUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse remote URL: %s", sanitizeGitDiagnostic(err.Error()))
	}

	// Remove any user info.
	parsedUrl.User = nil
	return parsedUrl, nil
}

func platformFromHost(u *url.URL) scm.Platform {
	switch strings.ToLower(u.Hostname()) {
	case "github.com":
		return scm.GitHubPlatform
	case "gitlab.com":
		return scm.GitLabPlatform
	case "dev.azure.com", "visualstudio.com":
		return scm.AzureDevOpsPlatform
	case "gitea.com", "code.forgejo.org", "codeberg.org":
		return scm.GiteaPlatform
	case "bitbucket.org":
		return scm.BitbucketPlatform
	default:
		return scm.UnknownPlatform
	}
}
