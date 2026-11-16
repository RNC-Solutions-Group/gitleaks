package detect

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zricethezav/gitleaks/v8/config"
)

func TestTotalBytesCountsFragmentContentOnce(t *testing.T) {
	d := NewDetector(config.Config{})
	content := []byte("synthetic content")
	d.Detect(Fragment{Raw: string(content), Bytes: content})
	assert.Equal(t, uint64(len(content)), d.TotalBytes.Load())
}
