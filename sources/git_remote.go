package sources

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"

	"github.com/zricethezav/gitleaks/v8/cmd/scm"
	"github.com/zricethezav/gitleaks/v8/logging"
)

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

	return normalizeGitRemoteURL(string(bytes.TrimSpace(stdout)))
}

// normalizeGitRemoteURL converts supported SCP-style remotes to their HTTPS
// form and removes credentials. Azure DevOps v3 SSH paths are mapped to the
// canonical dev.azure.com organization/project/repository URL.
func normalizeGitRemoteURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if matches := sshUrlpat.FindStringSubmatch(raw); matches != nil {
		host, path := matches[1], matches[2]
		if host == "ssh.dev.azure.com" && strings.HasPrefix(path, "v3/") {
			parts := strings.Split(strings.TrimPrefix(path, "v3/"), "/")
			if len(parts) >= 3 {
				org, project := parts[0], parts[1]
				repository := strings.Join(parts[2:], "/")
				raw = fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", org, project, repository)
			} else {
				raw = fmt.Sprintf("https://%s/%s", host, path)
			}
		} else {
			raw = fmt.Sprintf("https://%s/%s", host, path)
		}
	}
	raw = strings.TrimSuffix(raw, ".git")
	parsedUrl, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("unable to parse remote URL: %s", sanitizeGitDiagnostic(err.Error()))
	}
	// Remove any user info, including tokens embedded in HTTPS remotes.
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
