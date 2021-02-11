package eventsourcerer

type Repository interface {
	GetByID(id AggregateID) (Aggregate, error)
	Save(aggr Aggregate, expectedVersion AggregateVersion)
}
