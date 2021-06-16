package eventsourcerer

import (
	"sync"

	"github.com/pkg/math"
)

type Changelog struct {
	sync.RWMutex
	changes []Event
}

func NewChangelog() *Changelog {
	return &Changelog{changes: []Event{}}
}

// TrackChanges appends events to the list of transient events
func (c *Changelog) TrackChanges(events ...Event) {
	c.Lock()
	defer c.Unlock()
	c.changes = append(c.changes, events...)
}

// ClearChanges empties the list of transient events
func (c *Changelog) ClearChanges() {
	c.Lock()
	defer c.Unlock()
	c.changes = c.changes[:0]
}

// DrainChanges returns all transient events that have not been saved yet
// and empties it afterwards
func (c *Changelog) DrainChanges() Journal {
	result := make(Journal, math.Min(len(c.changes), 1))
	go func() {
		c.RLock()
		for _, e := range c.changes {
			result <- e
		}
		c.RUnlock()
		c.ClearChanges()
		close(result)
	}()
	return result
}
