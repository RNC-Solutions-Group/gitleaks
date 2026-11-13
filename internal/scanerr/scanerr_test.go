package scanerr

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectorEmpty(t *testing.T) {
	c := New(3)
	assert.NoError(t, c.Err())

	var nilCollector *Collector
	assert.NoError(t, nilCollector.Err())
}

func TestCollectorRetainsBoundedErrors(t *testing.T) {
	first := errors.New("first failure")
	second := errors.New("second failure")
	third := errors.New("third failure")

	c := New(2)
	c.Add(nil)
	c.Add(first)
	c.Add(second)
	c.Add(third)

	err := c.Err()
	require.Error(t, err)
	assert.ErrorIs(t, err, first)
	assert.ErrorIs(t, err, second)
	assert.NotErrorIs(t, err, third)
	assert.Contains(t, err.Error(), "1 additional scan error(s) omitted")
}

func TestCollectorConcurrentUse(t *testing.T) {
	c := New(5)
	const workers = 32
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Add(errors.New("worker failure"))
		}()
	}
	wg.Wait()

	err := c.Err()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "27 additional scan error(s) omitted")
	assert.Equal(t, 5, strings.Count(err.Error(), "worker failure"))
}
