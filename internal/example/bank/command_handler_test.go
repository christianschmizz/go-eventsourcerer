package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	es "github.com/christianschmizz/go-eventsourcerer"
)

func TestAccountCommandHandler_Handle(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	publisher := &EventPublisher{C: make(chan es.EventDescriptor, 1)}
	accountCommandHandler := NewAccountCommandHandler(NewAccountRepository(CreateEventStore(publisher)))
	boundedContext := es.CreateBoundedContext(accountCommandHandler, &ExternalEventHandler{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := boundedContext.Listen(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				panic(err)
			}
		}
	}()

	for x := 1; x <= 1000; x++ {
		var (
			accountOpenedEvent  AccountOpenedEvent
			accountNumber       AccountNumber
			desc es.EventDescriptor
		)

		{
			boundedContext.Commands() <- OpenAccountCommand{
				Username: fmt.Sprintf("Paul %d", x),
			}
			desc := <-publisher.C
			accountOpenedEvent = desc.Event.(AccountOpenedEvent)
			accountNumber = accountOpenedEvent.Number
		}
		{
			boundedContext.Commands() <- DepositMoneyCommand{
				AccountNumber:   accountNumber,
				Amount:          50,
				OriginalVersion: desc.Version,
			}
			desc = <-publisher.C
		}
		{
			boundedContext.Commands() <- WithdrawMoneyCommand{
				AccountNumber:   accountNumber,
				Amount:          50,
				OriginalVersion: desc.Version,
			}
			desc = <-publisher.C
		}
		{
			boundedContext.Commands() <- DepositMoneyCommand{
				AccountNumber:   accountNumber,
				Amount:          150,
				OriginalVersion: desc.Version,
			}
			desc = <-publisher.C
		}
	}
	cancel()
	wg.Wait()
}