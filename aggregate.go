package eventsourcerer

import "github.com/pkg/math"

type AggregateID int
type AggregateVersion int64

type Aggregate interface {
	ID() AggregateID
	Consume(event Event)
	Apply(events ...Event)
	Replay(journal Journal)
}

type BaseAggregate struct {
	Aggregate

	// transientEvents keeps collection of all events which have not been written to the event-store
	transientEvents []Event
}

func CreateBaseAggregate() *BaseAggregate {
	return &BaseAggregate{
		transientEvents: make([]Event, 0, 100),
	}
}

// Events will return all transient events as Journal
func (a *BaseAggregate) Journal() Journal {
	result := make(Journal, math.Min(len(a.transientEvents), 1))
	go func() {
		for _, e := range a.transientEvents {
			result <- e
		}
		close(result)
		a.Clear()
	}()
	return result
}

// Append appends events to the list of transient events
func (a *BaseAggregate) Append(events ...Event) {
	a.transientEvents = append(a.transientEvents, events...)
}

// Clear empties the list of transient events
func (a *BaseAggregate) Clear() {
	a.transientEvents = a.transientEvents[:0]
}
