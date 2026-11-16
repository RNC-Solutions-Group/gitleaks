package sources

import (
	"testing"

	"github.com/gitleaks/go-gitdiff/gitdiff"
	"github.com/stretchr/testify/assert"
)

func TestGitDiffFragmentWithContext(t *testing.T) {
	fragment := &gitdiff.TextFragment{
		NewPosition: 10,
		Lines: []gitdiff.Line{
			{Op: gitdiff.OpContext, Line: "before\n"},
			{Op: gitdiff.OpDelete, Line: "removed\n"},
			{Op: gitdiff.OpAdd, Line: "first\n"},
			{Op: gitdiff.OpAdd, Line: "second\n"},
			{Op: gitdiff.OpContext, Line: "after\n"},
		},
	}

	raw, added := gitDiffFragment(fragment, true)
	assert.Equal(t, "before\nfirst\nsecond\nafter\n", raw)
	assert.Equal(t, []LineRange{{Start: 11, End: 12}}, added)
}

func TestGitDiffFragmentDefaultIsAdditionsOnly(t *testing.T) {
	fragment := &gitdiff.TextFragment{
		NewPosition: 10,
		Lines: []gitdiff.Line{
			{Op: gitdiff.OpContext, Line: "before\n"},
			{Op: gitdiff.OpAdd, Line: "added\n"},
		},
	}
	raw, added := gitDiffFragment(fragment, false)
	assert.Equal(t, "added\n", raw)
	assert.Nil(t, added)
}
