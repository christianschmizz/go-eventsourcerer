package main

import (
	"testing"

	es "github.com/christianschmizz/go-eventsourcerer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(desc es.EventDescriptor) {
	m.Called(desc)
}

func TestCreateEventStore(t *testing.T) {
	publisher := &MockPublisher{}
	store := CreateEventStore(publisher)

	var aggr *Account

	t.Run("create aggregate", func(t *testing.T) {
		aggr = OpenAccount("Paul Panzer", 12345)
		assert.Equal(t, es.AggregateID(12345), aggr.ID())

		publisher.On("Publish", mock.Anything).Times(1)

		err := store.Save(aggr.ID(), aggr.Journal(), es.InitialVersion)
		assert.NoError(t, err, "failed to save aggregate: %+v", aggr)

		publisher.AssertExpectations(t)
	})

	t.Run("store aggregate's transient events", func(t *testing.T) {
		publisher.On("Publish", mock.Anything).Times(3)
		aggr.Consume(MoneyDepositedEvent{
			Amount: 100,
		})
		assert.Equal(t, 100, aggr.Balance)
		aggr.Consume(MoneyWithdrawnEvent{
			Amount: 50,
		})
		assert.Equal(t, 50, aggr.Balance)
		aggr.Consume(MoneyDepositedEvent{
			Amount: 200,
		})
		assert.Equal(t, 250, aggr.Balance)
		err := store.Save(aggr.ID(), aggr.Journal(), 0)
		assert.NoError(t, err, "failed to save aggregate: %+v", aggr)
		publisher.AssertExpectations(t)
	})

	t.Run("read journal", func(t *testing.T) {
		journal, err := store.GetJournalForAggregate(aggr.ID())
		assert.NoError(t, err)
		assert.Len(t, journal.SliceAndDrain(), 4)
	})

	t.Run("replay journal for aggregate", func(t *testing.T) {
		journal, err := store.GetJournalForAggregate(aggr.ID())
		assert.NoError(t, err)

		account := &Account{BaseAggregate: es.CreateBaseAggregate()}
		account.Replay(journal)

		assert.Equal(t, "Paul Panzer", account.AccountHolder)
		assert.Equal(t, AccountNumber(12345), account.Number)
		assert.Equal(t, 250, account.Balance)
		assert.Equal(t, 0, account.Limit)

		j := account.Journal()
		assert.Len(t, j.SliceAndDrain(), 0, "there should be no transient events as this was a replay")
	})

}
