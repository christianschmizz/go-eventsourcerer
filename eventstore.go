package eventsourcerer

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type EventStore interface {
	// GetJournalForAggregate retrieves all events of an Aggregate from a store with the given ID
	GetJournalForAggregate(id AggregateID) (Journal, error)

	// Save puts all journal-items for an Aggregate at the store as long as it matches the expected version
	Save(id AggregateID, journal Journal, expectedVersion AggregateVersion) error
}

// A ConcurrencyError signalizes a version mismatch
type ConcurrencyError struct {
	ExpectedVersion AggregateVersion
	CurrentVersion  AggregateVersion
}

func NewConcurrencyError(expectedVersion, currentVersion AggregateVersion) ConcurrencyError {
	return ConcurrencyError{
		ExpectedVersion: expectedVersion,
		CurrentVersion:  currentVersion,
	}
}

func (e ConcurrencyError) Error() string {
	return fmt.Sprintf("version did not match. expected version: %d got: %d", e.ExpectedVersion, e.CurrentVersion)
}

type MemoryEventStore struct {
	publisher EventPublisher
	events    map[AggregateID][]EventDescriptor
}

func NewMemoryEventStore(publisher EventPublisher) *MemoryEventStore {
	return &MemoryEventStore{
		publisher: publisher,
		events:    map[AggregateID][]EventDescriptor{},
	}
}

func (s *MemoryEventStore) GetJournalForAggregate(id AggregateID) (Journal, error) {
	descriptors, exists := s.events[id]
	if !exists {
		return nil, fmt.Errorf("aggregate not found: %d", id)
	}

	j := make(Journal, len(descriptors))
	go func() {
		for _, descriptor := range descriptors {
			j <- descriptor.Event
		}
		close(j)
	}()
	return j, nil
}

const InitialVersion AggregateVersion = -1
const IgnoreVersion AggregateVersion = -2

func (s *MemoryEventStore) Save(id AggregateID, journal Journal, expectedVersion AggregateVersion) error {
	descriptors, exists := s.events[id]
	if !exists {
		log.Debug().Msgf("aggregate not found: %d", id)
		descriptors = make([]EventDescriptor, 0, 1)
		s.events[id] = descriptors
	} else {
		lastDesc := descriptors[len(descriptors)-1]
		if expectedVersion == IgnoreVersion {
			expectedVersion = lastDesc.Version
		} else if expectedVersion != InitialVersion && lastDesc.Version != expectedVersion {
			return NewConcurrencyError(expectedVersion, lastDesc.Version)
		}
	}

	i := expectedVersion
	for event := range journal {
		i++

		desc := EventDescriptor{
			Event:   event,
			ID:      id,
			Version: i,
		}

		// Storing the event and its metadata in memory
		s.events[id] = append(s.events[id], desc)
		log.Debug().Msgf("event saved: %#v", event)

		// we cannot update the version of an event as it is a value object
		// and interface type. Publishing the descriptor for further
		// processing by subscribers is the only option for now.
		s.publisher.Publish(desc)
	}

	return nil
}
