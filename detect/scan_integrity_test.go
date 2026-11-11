package detect

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
	"github.com/zricethezav/gitleaks/v8/sources"
)

func scanStatusDetector() *Detector {
	cfg := config.Config{
		Rules: map[string]config.Rule{
			"test-rule": {
				RuleID:   "test-rule",
				Regex:    regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
				Keywords: []string{"akia"},
			},
		},
		Keywords: map[string]struct{}{"akia": {}},
	}
	return NewDetector(cfg)
}

type fragmentErrorSource struct {
	fragments []sources.Fragment
	errs      []error
	finalErr  error
}

func (s *fragmentErrorSource) Fragments(_ context.Context, yield sources.FragmentsFunc) error {
	for i, fragment := range s.fragments {
		var err error
		if i < len(s.errs) {
			err = s.errs[i]
		}
		if yieldErr := yield(fragment, err); yieldErr != nil {
			return yieldErr
		}
	}
	return s.finalErr
}

func TestDetectSourceReturnsFragmentErrorsAndPreservesFindings(t *testing.T) {
	readErr := errors.New("synthetic fragment read failure")
	detector := scanStatusDetector()
	findings, err := detector.DetectSource(context.Background(), &fragmentErrorSource{
		fragments: []sources.Fragment{
			{},
			{Raw: `aws_access_key_id = "AKIAABCDEFGHIJKLMNOP"`, FilePath: "secret.txt"},
		},
		errs: []error{readErr, nil},
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, readErr)
	require.Len(t, findings, 1)
	assert.Equal(t, "secret.txt", findings[0].File)
}

func TestDetectSourceReturnsSourceErrorAfterFindings(t *testing.T) {
	sourceErr := errors.New("synthetic source failure")
	detector := scanStatusDetector()
	findings, err := detector.DetectSource(context.Background(), &fragmentErrorSource{
		fragments: []sources.Fragment{
			{Raw: `aws_access_key_id = "AKIAABCDEFGHIJKLMNOP"`, FilePath: "secret.txt"},
		},
		finalErr: sourceErr,
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, sourceErr)
	require.Len(t, findings, 1)
}

type workerErrorSource struct {
	detector *Detector
	err      error
}

func (s *workerErrorSource) Fragments(_ context.Context, _ sources.FragmentsFunc) error {
	s.detector.Sema.Go(func() error { return s.err })
	return nil
}

func TestDetectSourceReturnsWorkerErrors(t *testing.T) {
	workerErr := errors.New("synthetic worker failure")
	detector := scanStatusDetector()
	_, err := detector.DetectSource(context.Background(), &workerErrorSource{
		detector: detector,
		err:      workerErr,
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, workerErr)
}

type cancellationSource struct{}

func (cancellationSource) Fragments(ctx context.Context, _ sources.FragmentsFunc) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestDetectSourceReturnsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := scanStatusDetector().DetectSource(ctx, cancellationSource{})
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestDetectSourceCleanAndFindingControls(t *testing.T) {
	detector := scanStatusDetector()
	findings, err := detector.DetectSource(context.Background(), &fragmentErrorSource{
		fragments: []sources.Fragment{
			{Raw: "ordinary content", FilePath: "clean.txt"},
		},
	})
	require.NoError(t, err)
	assert.Empty(t, findings)

	detector = scanStatusDetector()
	findings, err = detector.DetectSource(context.Background(), &fragmentErrorSource{
		fragments: []sources.Fragment{
			{Raw: `aws_access_key_id = "AKIAABCDEFGHIJKLMNOP"`, FilePath: "secret.txt"},
		},
	})
	require.NoError(t, err)
	require.Len(t, findings, 1)
	assert.Equal(t, "test-rule", findings[0].RuleID)
}
