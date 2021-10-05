package eventsourcerer

// The EventDescriptor adds some required information to an Event as themselves are not mutable.
type EventDescriptor struct {
	Event   Event            `json:"event"`
	ID      AggregateID      `json:"aggregate_id"`
	Version AggregateVersion `json:"aggregate_version"`
}
