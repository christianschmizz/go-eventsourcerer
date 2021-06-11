package eventsourcerer

type Repository interface {
	// GetByID loads the Root Aggregate from the store
	GetByID(id AggregateID) (Aggregate, error)

	// Save stores the Aggregate and ensures that it matches the expected version before
	Save(aggr Aggregate, expectedVersion AggregateVersion)
}
