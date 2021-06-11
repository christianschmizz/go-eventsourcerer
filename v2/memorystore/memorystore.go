package memorystore

import (
	es2 "github.com/christianschmizz/go-eventsourcerer/v2"
	"github.com/rs/zerolog"
)

type memoryEventStore struct {
	es2.EventStore
	publisher es2.Publisher
	events    map[es2.AggregateID][]es2.EventDescriptor
	logger    zerolog.Logger
}

func NewStore(publisher es2.Publisher, logger zerolog.Logger) *memoryEventStore {
	return &memoryEventStore{
		publisher: publisher,
		events:    map[es2.AggregateID][]es2.EventDescriptor{},
		logger:    logger,
	}
}

func (s *memoryEventStore) LoadJournal(id es2.AggregateID) (es2.Journal, error) {
	descriptors, exists := s.events[id]
	if !exists {
		return nil, es2.NewAggregateDoesNotExistError(id)
	}

	j := make(es2.Journal, len(descriptors))
	go func() {
		for _, descriptor := range descriptors {
			j <- descriptor.Event
		}
		close(j)
	}()
	return j, nil
}

func (s *memoryEventStore) SaveJournal(id es2.AggregateID, journal es2.Journal, expectedVersion es2.AggregateVersion) error {
	descriptors, exists := s.events[id]
	if !exists {
		s.logger.Debug().Msgf("aggregate not found: %d", id)
		descriptors = make([]es2.EventDescriptor, 0, 1)
		s.events[id] = descriptors
	} else {
		lastDesc := descriptors[len(descriptors)-1]
		if expectedVersion == es2.IgnoreAggregateVersion {
			expectedVersion = lastDesc.Version
		} else if expectedVersion != es2.InitialAggregateVersion && lastDesc.Version != expectedVersion {
			return es2.NewConcurrencyError(expectedVersion, lastDesc.Version)
		}
	}

	version := expectedVersion
	for item := range journal {
		version++

		// We cannot update the version of an Event as it is a value object
		// and interface type. Publishing the Descriptor for further
		// processing by subscribers is the only option for now.
		desc := es2.EventDescriptor{
			Event:   item,
			ID:      id,
			Version: version,
		}

		// Storing the event and its metadata in memory
		s.events[id] = append(s.events[id], desc)

		s.logger.Debug().Int64("aggregate_id", int64(id)).Int64("version", int64(version)).Msgf("event saved: %#v", item)

		s.publisher.Publish(desc)
	}

	return nil
}
