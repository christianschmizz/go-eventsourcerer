package main

import (
	"testing"
	"time"

	es "github.com/christianschmizz/go-eventsourcerer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetJournalForAggregate(id es.AggregateID) (es.Journal, error) {
	args := m.Called(id)
	return args.Get(0).(es.Journal), args.Error(1)
}

func (m *MockStore) Save(id es.AggregateID, journal es.Journal, expectedVersion es.AggregateVersion) error {
	args := m.Called(id, journal, expectedVersion)
	return args.Error(0)
}

func TestAccountRepository_GetByID(t *testing.T) {
	store := &MockStore{}
	repo := NewAccountRepository(store)

	// Prepare a journal
	j := make(es.Journal, 10)
	j <- AccountOpenedEvent{
		BaseEvent:     es.BaseEvent{},
		AccountHolder: "Jane Doe",
		Number:        123,
		Balance:       500,
	}
	j <- MoneyDepositedEvent{
		BaseEvent: es.BaseEvent{},
		Amount:    100,
	}
	j <- MoneyWithdrawnEvent{
		BaseEvent: es.BaseEvent{},
		Amount:    50,
	}
	j <- AccountClosedEvent{
		BaseEvent: es.BaseEvent{},
		ClosedAt:  time.Now(),
	}
	close(j)

	// Apply the prepared journal to the mock
	store.On("GetJournalForAggregate", es.AggregateID(1)).Return(j, nil).Times(1)

	// Retrieve the account
	aggr, err := repo.GetByID(es.AggregateID(1))
	assert.NoError(t, err)
	assert.NotNil(t, aggr)

	store.AssertExpectations(t)

	// Check account for correct details
	account := aggr.(*Account)
	assert.Equal(t, 550, account.Balance)
	assert.Equal(t, "Jane Doe", account.AccountHolder)
	assert.Equal(t, 123, account.Number)
}

func TestAccountRepository_Save(t *testing.T) {
	store := &MockStore{}
	repo := NewAccountRepository(store)

	acc := OpenAccount("Jane Doe", 12345)

	store.On("Save", acc.ID(), mock.Anything, es.AggregateVersion(-1)).Return(nil).Times(1)

	repo.Save(acc, -1)
	store.AssertExpectations(t)

	assert.Len(t, acc.Journal().SliceAndDrain(), 0)
}