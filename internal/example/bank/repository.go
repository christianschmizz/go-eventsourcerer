package main

import (
	es "github.com/christianschmizz/go-eventsourcerer"
	"github.com/rs/zerolog/log"
)

type AccountRepository struct {
	store es.EventStore
}

func NewAccountRepository(eventStore es.EventStore) *AccountRepository {
	return &AccountRepository{
		store: eventStore,
	}
}

func (r *AccountRepository) GetByID(aggregateID es.AggregateID) (es.Aggregate, error) {
	journal, err := r.store.GetJournalForAggregate(aggregateID)
	if err != nil {
		return nil, err
	}

	var account = NewAccount()
	account.Replay(journal)
	return account, nil
}

func (r *AccountRepository) Save(accMaybe es.Aggregate, expectedVersion es.AggregateVersion) {
	acc := accMaybe.(*Account)
	slog := log.With().Int("aggregate_id", int(acc.ID())).Logger()
	slog.Debug().Msg("trying to save aggregate")
	if err := r.store.Save(acc.ID(), acc.Journal(), expectedVersion); err != nil {
		slog.Error().Err(err).Msg("failed to save")
		panic(err)
	}
	slog.Debug().Msg("saved")
}
