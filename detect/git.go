package detect

import (
	"context"

	"github.com/zricethezav/gitleaks/v8/cmd/scm"
	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"
)

// RemoteInfo is an alias for sources.RemoteInfo for backwards compatibility
//
// Deprecated: This will be replaced with sources.RemoteInfo in v9
type RemoteInfo sources.RemoteInfo

// DetectGit runs detections against a GitCmd with its remote info
//
// Deprecated: Use sources.Git and detector.DetectSource instead
func (d *Detector) DetectGit(cmd *sources.GitCmd, remote *RemoteInfo) ([]report.Finding, error) {
	return d.DetectGitContext(context.Background(), cmd, remote)
}

// DetectGitContext runs detections against a GitCmd with its remote info and
// propagates the supplied context to the source and detector workers.
func (d *Detector) DetectGitContext(ctx context.Context, cmd *sources.GitCmd, remote *RemoteInfo) ([]report.Finding, error) {
	return d.DetectSource(
		ctx,
		&sources.Git{
			Cmd:             cmd,
			Config:          &d.Config,
			Remote:          (*sources.RemoteInfo)(remote),
			Sema:            d.Sema,
			MaxFileSize:     d.MaxTargetMegaBytes * 1_000_000,
			MaxArchiveDepth: d.MaxArchiveDepth,
			Informative:     d.Informative,
		},
	)
}

// Deprecated: use sources.NewRemoteInfo instead
func NewRemoteInfo(platform scm.Platform, source string) *RemoteInfo {
	return (*RemoteInfo)(sources.NewRemoteInfo(platform, source))
}
