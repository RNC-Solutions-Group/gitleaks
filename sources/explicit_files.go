package sources

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
}

func (s *ExplicitFiles) Fragments(ctx context.Context, yield FragmentsFunc) error {
	collector := scanerr.New(10)
	seen := make(map[string]struct{}, len(s.Paths))
	for _, requested := range s.Paths {
		path := filepath.Clean(requested)
		if _, duplicate := seen[path]; duplicate {
			continue
		}
		seen[path] = struct{}{}
		info, err := os.Lstat(path)
		if errors.Is(err, os.ErrNotExist) {
			collector.Add(fmt.Errorf("%s: explicit file does not exist", path))
			continue
		}
		if err != nil {
			collector.Add(fmt.Errorf("%s: could not inspect file: %w", path, err))
			continue
		}
		if info.IsDir() {
			collector.Add(fmt.Errorf("%s: explicit file path is a directory", path))
			continue
		}

		scanPath := path
		symlink := ""
		if info.Mode()&os.ModeSymlink != 0 {
			if !s.FollowSymlinks {
				logging.Debug().Str("path", path).Msg("skipping explicit symlink")
				continue
			}
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				collector.Add(fmt.Errorf("%s: could not resolve symlink: %w", path, err))
				continue
			}
			targetInfo, err := os.Stat(realPath)
			if err != nil {
				collector.Add(fmt.Errorf("%s: could not inspect symlink target: %w", path, err))
				continue
			}
			if targetInfo.IsDir() {
				collector.Add(fmt.Errorf("%s: explicit symlink target is a directory", path))
				continue
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
			continue
		}

		content, err := os.Open(scanPath)
		if err != nil {
			collector.Add(fmt.Errorf("%s: could not open file: %w", path, err))
			continue
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
				collector.Add(fmt.Errorf("%s: %w (close: %v)", path, fragmentsErr, closeErr))
			} else {
				collector.Add(fmt.Errorf("%s: %w", path, fragmentsErr))
			}
			continue
		}
		if closeErr != nil {
			collector.Add(fmt.Errorf("%s: could not close file: %w", path, closeErr))
			continue
		}
	}
	return collector.Err()
}
