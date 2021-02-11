package eventsourcerer

type EventPublisher interface {
	Publish(EventDescriptor)
}
