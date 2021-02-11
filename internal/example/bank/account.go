package main

import (
	"time"

	es "github.com/christianschmizz/go-eventsourcerer"
	"github.com/rs/zerolog/log"
)

type AccountNumber int

// The Account represents the Aggregate Root (Domain objects) of our example
type Account struct {
	*es.BaseAggregate
	Number        AccountNumber
	AccountHolder string
	Balance       int
	Limit         int
	OpenedAt      time.Time
	ClosedAt      time.Time
	FrozenAt      time.Time
	UnfrozenAt    time.Time
}

func NewAccount() *Account {
	return &Account{
		BaseAggregate: es.CreateBaseAggregate(),
	}
}

func OpenAccount(accountHolder string, number int) *Account {
	account := NewAccount()
	account.Consume(AccountOpenedEvent{
		AccountHolder: accountHolder,
		Number:        AccountNumber(number),
		Balance:       0,
		OpenedAt:      time.Now(),
	})
	return account
}

func (a *Account) ID() es.AggregateID {
	return es.AggregateID(a.Number)
}

func (a *Account) Consume(event es.Event) {
	a.Apply(event)
	a.Append(event)
}

func (a *Account) Replay(journal es.Journal) {
	for event := range journal {
		a.Apply(event)
	}
}

func (a *Account) Apply(events ...es.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case AccountOpenedEvent:
			a.AccountHolder = e.AccountHolder
			a.Number = e.Number
			a.Balance = e.Balance
		case AccountClosedEvent:
			a.ClosedAt = e.ClosedAt
		case MoneyDepositedEvent:
			a.Balance += e.Amount
		case MoneyWithdrawnEvent:
			a.Balance -= e.Amount
		case AccountFrozenEvent:
			a.FrozenAt = e.FrozenAt
		case AccountUnfrozenEvent:
			a.UnfrozenAt = e.UnfrozenAt
			a.FrozenAt = time.Time{}
		default:
			log.Warn().Msgf("unknown event: %#v", e)
		}
	}
}

func (a *Account) IsFrozen() bool {
	return !a.FrozenAt.IsZero()
}

func (a *Account) IsClosed() bool {
	return !a.ClosedAt.IsZero()
}

func (a *Account) Deposit(amount int) error {
	a.Consume(MoneyDepositedEvent{
		Amount: amount,
	})
	return nil
}

func (a *Account) Withdraw(amount int) error {
	a.Consume(MoneyWithdrawnEvent{
		Amount: amount,
	})
	return nil
}

func (a *Account) Close() error {
	a.Consume(AccountClosedEvent{
		ClosedAt: time.Now(),
	})
	return nil
}

func (a *Account) Freeze() error {
	a.Consume(AccountFrozenEvent{
		FrozenAt: time.Now(),
	})
	return nil
}

func (a *Account) Unfreeze() error {
	a.Consume(AccountUnfrozenEvent{
		UnfrozenAt: time.Now(),
	})
	return nil
}
