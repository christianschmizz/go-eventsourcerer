package main

import (
	es "github.com/christianschmizz/go-eventsourcerer"
)

type EventStore struct {
	*es.MemoryEventStore
}

func CreateEventStore(publisher es.EventPublisher) *EventStore {
	return &EventStore{es.NewMemoryEventStore(publisher)}
}
