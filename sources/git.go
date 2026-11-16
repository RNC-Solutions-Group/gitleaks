package sources

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/fatih/semgroup"
	"github.com/gitleaks/go-gitdiff/gitdiff"

	"github.com/zricethezav/gitleaks/v8/cmd/scm"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
	"github.com/zricethezav/gitleaks/v8/internal/scanerr"
	"github.com/zricethezav/gitleaks/v8/logging"
)

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
	MaxFileSize     int
	MaxArchiveDepth int
	// Informative is the optional coverage collector. Nil disables collection.
	Informative *informative.Collector
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

// isCommitAllowed reports whether the commit SHA is covered by any global
// allowlist. It returns the matched allowlist's commit value so callers can
// trace which entry applied. It stops at the first match, since one match is
// enough to skip the commit.
func isCommitAllowed(allowlists []*config.Allowlist, sha string) (bool, string) {
	for _, a := range allowlists {
		if ok, c := a.CommitAllowed(sha); ok {
			return true, c
		}
	}
	return false, ""
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

			// Check if commit is allowed
			commitSHA := ""
			var commitInfo *CommitInfo
			if gitdiffFile.PatchHeader != nil {
				commitSHA = gitdiffFile.PatchHeader.SHA
				if ok, c := isCommitAllowed(s.Config.Allowlists, commitSHA); ok {
					logging.Trace().Str("allowed-commit", c).Msg("skipping commit: global allowlist")
					continue
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
					target := s.Informative.Track(
						informative.KindGitDiff,
						gitdiffFile.NewName,
						commitSHA,
						0,
						false,
					)
					target.Exclude(informative.ReasonBinaryContent, informative.ExclusionOptions{
						Detail: "git reported binary content",
					})
					continue
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

		target := s.Informative.Track(
			informative.KindGitBlob,
			diffFile.NewName,
			commitSHA,
			0,
			false,
		)
		file := File{
			Content:         blob,
			Path:            diffFile.NewName,
			MaxFileSize:     s.MaxFileSize,
			MaxArchiveDepth: s.MaxArchiveDepth,
			Config:          s.Config,
			Informative:     s.Informative,
			Target:          target,
			Commit:          commitSHA,
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

	var target *informative.Target
	for _, textFragment := range diffFile.TextFragments {
		if textFragment == nil {
			continue
		}
		if target == nil {
			target = s.Informative.Track(
				informative.KindGitDiff,
				diffFile.NewName,
				commitSHA,
				0,
				false,
			)
		}

		fragment := Fragment{
			CommitSHA:         commitSHA,
			FilePath:          diffFile.NewName,
			Raw:               textFragment.Raw(gitdiff.OpAdd),
			StartLine:         int(textFragment.NewPosition),
			CommitInfo:        commitInfo,
			InformativeTarget: target,
		}

		if err := yield(fragment, nil); err != nil {
			return err
		}
	}

	return nil
}
