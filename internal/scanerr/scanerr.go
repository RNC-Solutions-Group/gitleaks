// Package scanerr provides bounded error accumulation for scan operations.
package scanerr

import (
	"errors"
	"fmt"
	"sync"
)

// DefaultLimit is the number of individual errors retained by default.
const DefaultLimit = 10

// Collector accumulates operational errors without allowing a failed scan to
// produce an unbounded error value. It is safe for concurrent use.
type Collector struct {
	mu      sync.Mutex
	limit   int
	errs    []error
	total   int
	dropped int
}

// New returns a collector that retains at most limit errors. A non-positive
// limit uses DefaultLimit.
func New(limit int) *Collector {
	if limit <= 0 {
		limit = DefaultLimit
	}
	return &Collector{limit: limit}
}

// Add records err if it is non-nil.
func (c *Collector) Add(err error) {
	if err == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.total++
	if len(c.errs) < c.limit {
		c.errs = append(c.errs, err)
		return
	}
	c.dropped++
}

// Err returns nil when no errors were added. Otherwise it returns a joined
// error containing the retained errors and, when applicable, a summary of the
// omitted errors.
func (c *Collector) Err() error {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.total == 0 {
		return nil
	}

	errs := make([]error, 0, len(c.errs)+1)
	errs = append(errs, c.errs...)
	if c.dropped > 0 {
		errs = append(errs, fmt.Errorf("%d additional scan error(s) omitted", c.dropped))
	}
	return errors.Join(errs...)
}
