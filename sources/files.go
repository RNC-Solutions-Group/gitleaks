package sources

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fatih/semgroup"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/internal/scanerr"
	"github.com/zricethezav/gitleaks/v8/logging"
)

// TODO: remove this in v9 and have scanTargets yield file sources
type ScanTarget struct {
	Path    string
	Symlink string
	// Size is the target file's size when known.
	Size int64
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
}

// scanTargets yields scan targets to a callback func
func (s *Files) scanTargets(ctx context.Context, yield func(ScanTarget, error) error) error {
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

		if !d.IsDir() {
			// Empty; nothing to do here.
			if info.Size() == 0 {
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
				logger.Debug().Str("target", realPath).Msgf("skipping symlink: target is directory")
				return nil
			}
			scanTarget = ScanTarget{
				Path:    realPath,
				Symlink: path,
				Size:    realPathFileInfo.Size(),
			}
		}

		// handle dir cases (mainly just see if it should be skipped
		if info.IsDir() {
			if shouldSkipPath(s.Config, path) {
				logger.Debug().Msg("skipping directory: global allowlist")
				return filepath.SkipDir
			}
			return nil
		}

		if shouldSkipPath(s.Config, path) {
			logger.Debug().Msg("skipping file: global allowlist")
			return nil
		}

		if err := ctx.Err(); err != nil {
			return err
		}

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
		Config:          s.Config,
		MaxFileSize:     s.MaxFileSize,
		Size:            scanTarget.Size,
		MaxArchiveDepth: s.MaxArchiveDepth,
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
