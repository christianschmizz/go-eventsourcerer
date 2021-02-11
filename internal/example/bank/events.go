package main

import (
	"time"

	es "github.com/christianschmizz/go-eventsourcerer"
)

type AccountOpenedEvent struct {
	es.BaseEvent
	AccountHolder string
	Number        AccountNumber
	Balance       int
	OpenedAt      time.Time
}

type AccountClosedEvent struct {
	es.BaseEvent
	ClosedAt time.Time
}

type MoneyDepositedEvent struct {
	es.BaseEvent
	Amount int
}

type MoneyWithdrawnEvent struct {
	es.BaseEvent
	Amount int
}

type AccountFrozenEvent struct {
	es.BaseEvent
	FrozenAt time.Time
}

type AccountUnfrozenEvent struct {
	es.BaseEvent
	UnfrozenAt time.Time
}
