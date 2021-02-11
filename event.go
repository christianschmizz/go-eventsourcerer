package eventsourcerer

type Event interface {
	//	Version() AggregateVersion
}

type BaseEvent struct {
	Event
}

// An EventHandler listens to (external) events in order to issue commands.
type EventHandler interface {
	Handle(Event, chan<- Command)
}

type EventDescriptor struct {
	Event   Event
	ID      AggregateID
	Version AggregateVersion
}
