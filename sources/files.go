package sources

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fatih/semgroup"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
	"github.com/zricethezav/gitleaks/v8/internal/scanerr"
	"github.com/zricethezav/gitleaks/v8/logging"
)

// TODO: remove this in v9 and have scanTargets yield file sources
type ScanTarget struct {
	Path    string
	Symlink string
	// Root is the directory scan root, when the target came from one. It is
	// empty for a single-file target.
	Root string
	// Size is the target file's size when known.
	Size int64
	// Informative is the optional coverage-record handle for this target.
	Informative *informative.Target
}

// Deprecated: Use Files and detector.DetectSource instead
func DirectoryTargets(sourcePath string, s *semgroup.Group, followSymlinks bool, allowlists []*config.Allowlist) (<-chan ScanTarget, error) {
	paths := make(chan ScanTarget)

	// create a Files source
	files := Files{
		FollowSymlinks: followSymlinks,
		Path:           sourcePath,
		Sema:           s,
		Config: &config.Config{
			Allowlists: allowlists,
		},
	}

	s.Go(func() error {
		ctx := context.Background()
		err := files.scanTargets(ctx, func(scanTarget ScanTarget, err error) error {
			if err != nil {
				return err
			}
			paths <- scanTarget
			return nil
		})
		close(paths)
		return err
	})

	return paths, nil
}

// Files is a source for yielding fragments from a collection of files
type Files struct {
	Config          *config.Config
	FollowSymlinks  bool
	MaxFileSize     int
	Path            string
	Sema            *semgroup.Group
	MaxArchiveDepth int
	// Informative is the optional coverage collector. Nil disables collection.
	Informative *informative.Collector

	// root is Path when it names a directory. It lets path matching evaluate
	// both the spelling the user supplied and the path relative to the scan
	// root. It stays empty for a single-file target.
	root string
}

// scanTargets yields scan targets to a callback func
func (s *Files) scanTargets(ctx context.Context, yield func(ScanTarget, error) error) error {
	s.root = ""
	if info, err := os.Stat(s.Path); err == nil && info.IsDir() {
		s.root = s.Path
	}

	return filepath.WalkDir(s.Path, func(path string, d fs.DirEntry, err error) error {
		scanTarget := ScanTarget{Path: path}
		logger := logging.With().Str("path", path).Logger()

		if err != nil {
			return fmt.Errorf("could not walk %q: %w", path, err)
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("could not inspect %q: %w", path, err)
		}

		var target *informative.Target
		if !d.IsDir() {
			target = s.Informative.Track(informative.KindFile, path, "", info.Size(), true)

			// Empty; nothing to do here.
			if info.Size() == 0 {
				target.Exclude(informative.ReasonEmptyFile, informative.ExclusionOptions{
					Detail: "empty file",
				})
				logger.Debug().Msg("skipping empty file")
				return nil
			}

			// Size enforcement happens in File after archive verification so
			// an archive's members, not merely its compressed outer size, are
			// bounded.
			scanTarget.Size = info.Size()
		}

		// set the initial scan target values
		if d.Type() == fs.ModeSymlink {
			if !s.FollowSymlinks {
				target.Exclude(informative.ReasonSymlinkNotFollowed, informative.ExclusionOptions{
					OmitSize: true,
					Detail:   "follow symlinks disabled",
				})
				logger.Debug().Msg("skipping symlink: follow symlinks disabled")
				return nil
			}
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return fmt.Errorf("could not resolve symlink %q: %w", path, err)
			}
			realPathFileInfo, err := os.Stat(realPath)
			if err != nil {
				return fmt.Errorf("could not inspect symlink target %q: %w", realPath, err)
			}
			if realPathFileInfo.IsDir() {
				target.Exclude(informative.ReasonSymlinkNotFollowed, informative.ExclusionOptions{
					OmitSize: true,
					Detail:   "symlink target is a directory",
				})
				logger.Debug().Str("target", realPath).Msgf("skipping symlink: target is directory")
				return nil
			}
			if realPathFileInfo.Size() == 0 {
				// Preserve the historical traversal path: the empty target is
				// still opened and yields no fragments, while the informative
				// report records why no content was examined.
				target.Exclude(informative.ReasonEmptyFile, informative.ExclusionOptions{
					Detail: "empty symlink target",
				})
			}
			target.SetSize(realPathFileInfo.Size(), true)
			scanTarget = ScanTarget{
				Path:        realPath,
				Symlink:     path,
				Size:        realPathFileInfo.Size(),
				Informative: target,
			}
		}

		// handle dir cases (mainly just see if it should be skipped
		if info.IsDir() {
			if shouldSkipPathAtRoot(s.Config, s.root, path) {
				dirTarget := s.Informative.Track(informative.KindDirectory, path, "", 0, false)
				dirTarget.Exclude(informative.ReasonPathAllowlist, informative.ExclusionOptions{
					Detail: "global allowlist",
				})
				logger.Debug().Msg("skipping directory: global allowlist")
				return filepath.SkipDir
			}
			return nil
		}

		if shouldSkipPathAtRoot(s.Config, s.root, path) {
			target.Exclude(informative.ReasonPathAllowlist, informative.ExclusionOptions{
				Detail: "global allowlist",
			})
			logger.Debug().Msg("skipping file: global allowlist")
			return nil
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		scanTarget.Root = s.root
		scanTarget.Informative = target
		return yield(scanTarget, nil)
	})
}

// Fragments yields fragments from files discovered under the path
func (s *Files) Fragments(ctx context.Context, yield FragmentsFunc) error {
	if s.Sema == nil {
		s.Sema = semgroup.NewGroup(ctx, 40)
	}
	collector := scanerr.New(10)

	walkErr := s.scanTargets(ctx, func(scanTarget ScanTarget, err error) error {
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}

		s.Sema.Go(func() error {
			if err := s.scanFile(ctx, scanTarget, yield); err != nil {
				collector.Add(fmt.Errorf("%s: %w", scanTarget.Path, err))
			}
			return nil
		})
		return nil
	})
	if walkErr != nil {
		collector.Add(walkErr)
	}

	// Always join workers before returning, including after cancellation.
	collector.Add(s.Sema.Wait())
	if walkErr == nil {
		collector.Add(ctx.Err())
	}

	return collector.Err()
}

func (s *Files) scanFile(ctx context.Context, scanTarget ScanTarget, yield FragmentsFunc) error {
	f, err := os.Open(scanTarget.Path)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	file := File{
		Content:         f,
		Path:            scanTarget.Path,
		Symlink:         scanTarget.Symlink,
		Root:            scanTarget.Root,
		Config:          s.Config,
		MaxFileSize:     s.MaxFileSize,
		Size:            scanTarget.Size,
		MaxArchiveDepth: s.MaxArchiveDepth,
		Informative:     s.Informative,
		Target:          scanTarget.Informative,
	}

	fragmentsErr := file.Fragments(ctx, yield)
	closeErr := f.Close()
	if fragmentsErr != nil {
		return fragmentsErr
	}
	if closeErr != nil {
		return fmt.Errorf("could not close file: %w", closeErr)
	}
	return nil
}
