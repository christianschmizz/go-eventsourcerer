package eventsourcerer

type EventDescriptor struct {
	Event   Event            `json:"event"`
	ID      AggregateID      `json:"aggregate_id"`
	Version AggregateVersion `json:"aggregate_version"`
}
