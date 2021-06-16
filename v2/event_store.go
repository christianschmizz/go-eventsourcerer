package eventsourcerer

const (
	InitialAggregateVersion AggregateVersion = -1
	IgnoreAggregateVersion  AggregateVersion = -2
)

type EventStore interface {
	// LoadJournal retrieves all events of an Aggregate from a store with the given AggregateID
	LoadJournal(id AggregateID) (Journal, error)

	// Save puts all journal-items for an Aggregate at the store as long as it matches the expected version
	SaveJournal(id AggregateID, journal Journal, expectedVersion AggregateVersion) error
}
