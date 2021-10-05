package eventsourcerer

const (
	TransientEventBufferSize = 100
)

type AggregateID int64
type AggregateVersion int64

// The Aggregate (Aggregate Root) is to control and encapsulate access
// to its members in such a way as to protect it’s invariants.
type Aggregate interface {
	// ID returns the unique identifier of the aggregate
	ID() AggregateID

	// Consume applies the given Event to the Aggregate and append it to the transient
	// events list for persisting it later.
	Consume(event Event)

	// Replay applies all events of the Journal on top the aggregate
	Replay(journal Journal)
}
