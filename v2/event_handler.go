package eventsourcerer

// An EventHandler listens to (external) events in order to issue commands from it.
type EventHandler interface {
	Handle(Event, chan<- Command)
}
