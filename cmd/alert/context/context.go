// Package context contains type information for an alert's context, which is
// used to keep track of data an alert has accumulated through its life-cycle.
// This gets its own package to prevent import cycles.
package context

import (
	"time"
)

// Context describes information about an alert it accumulates through its
// life-cycle
type Context struct {
	Name      string
	StartedTS uint64
	time.Time `luautil:"-"`
}
