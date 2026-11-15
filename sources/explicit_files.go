package sources

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/semgroup"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/internal/informative"
	"github.com/zricethezav/gitleaks/v8/internal/scanerr"
	"github.com/zricethezav/gitleaks/v8/logging"
)

// ExplicitFiles yields only the named files, without walking their parent
// directories. It is the bounded source used by the pre-commit explicit-file
// hook: a filename list must never silently expand into a repository scan.
type ExplicitFiles struct {
	Paths           []string
	Config          *config.Config
	FollowSymlinks  bool
	MaxFileSize     int
	MaxArchiveDepth int
	Informative     *informative.Collector
	// Sema bounds concurrent file work. Nil creates a private group with the
	// scanner's standard limit; callers normally pass the detector's group.
	Sema *semgroup.Group
}

func (s *ExplicitFiles) Fragments(ctx context.Context, yield FragmentsFunc) error {
	if s.Sema == nil {
		s.Sema = semgroup.NewGroup(ctx, 40)
	}
	collector := scanerr.New(10)
	seen := make(map[string]struct{}, len(s.Paths))
	for _, requested := range s.Paths {
		path := filepath.Clean(requested)
		if _, duplicate := seen[path]; duplicate {
			continue
		}
		seen[path] = struct{}{}
		scanPathArg := path
		s.Sema.Go(func() error {
			if err := s.scanPath(ctx, scanPathArg, yield); err != nil {
				collector.Add(err)
			}
			return nil
		})
	}
	collector.Add(s.Sema.Wait())
	return collector.Err()
}

func (s *ExplicitFiles) scanPath(ctx context.Context, path string, yield FragmentsFunc) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s: explicit file does not exist", path)
	}
	if err != nil {
		return fmt.Errorf("%s: could not inspect file: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s: explicit file path is a directory", path)
	}

	scanPath := path
	symlink := ""
	if info.Mode()&os.ModeSymlink != 0 {
		if !s.FollowSymlinks {
			logging.Debug().Str("path", path).Msg("skipping explicit symlink")
			return nil
		}
		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return fmt.Errorf("%s: could not resolve symlink: %w", path, err)
		}
		targetInfo, err := os.Stat(realPath)
		if err != nil {
			return fmt.Errorf("%s: could not inspect symlink target: %w", path, err)
		}
		if targetInfo.IsDir() {
			return fmt.Errorf("%s: explicit symlink target is a directory", path)
		}
		scanPath = realPath
		symlink = path
		info = targetInfo
	}

	target := s.Informative.Track(informative.KindFile, path, "", info.Size(), true)
	if info.Size() == 0 {
		target.Exclude(informative.ReasonEmptyFile, informative.ExclusionOptions{
			Detail: "empty file",
		})
		return nil
	}

	content, err := os.Open(scanPath)
	if err != nil {
		return fmt.Errorf("%s: could not open file: %w", path, err)
	}
	file := &File{
		Content:         content,
		Path:            scanPath,
		Symlink:         symlink,
		Config:          s.Config,
		MaxFileSize:     s.MaxFileSize,
		Size:            info.Size(),
		MaxArchiveDepth: s.MaxArchiveDepth,
		Informative:     s.Informative,
		Target:          target,
	}
	fragmentsErr := file.Fragments(ctx, yield)
	closeErr := content.Close()
	if fragmentsErr != nil {
		if closeErr != nil {
			return fmt.Errorf("%s: %w (close: %v)", path, fragmentsErr, closeErr)
		}
		return fmt.Errorf("%s: %w", path, fragmentsErr)
	}
	if closeErr != nil {
		return fmt.Errorf("%s: could not close file: %w", path, closeErr)
	}
	return nil
}
