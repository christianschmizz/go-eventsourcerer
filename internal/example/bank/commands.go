package main

import es "github.com/christianschmizz/go-eventsourcerer"

type OpenAccountCommand struct {
	es.BaseCommand
	Username string
}

type CloseAccountCommand struct {
	es.BaseCommand
	AccountNumber   AccountNumber
	OriginalVersion es.AggregateVersion
}

type FreezeAccountCommand struct {
	es.BaseCommand
	AccountNumber   AccountNumber
	OriginalVersion es.AggregateVersion
}

type UnfreezeAccountCommand struct {
	es.BaseCommand
	AccountNumber   AccountNumber
	OriginalVersion es.AggregateVersion
}

type DepositMoneyCommand struct {
	es.BaseCommand
	Amount          int
	AccountNumber   AccountNumber
	OriginalVersion es.AggregateVersion
}

type WithdrawMoneyCommand struct {
	es.BaseCommand
	Amount          int
	AccountNumber   AccountNumber
	OriginalVersion es.AggregateVersion
}
