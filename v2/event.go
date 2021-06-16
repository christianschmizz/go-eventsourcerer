package eventsourcerer

type Event interface {
	//	Version() AggregateVersion
}

type EventBase struct {
	Event
}
