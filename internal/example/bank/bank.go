package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	es "github.com/christianschmizz/go-eventsourcerer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)

	ctx, cancel := context.WithCancel(context.Background())

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	publisher := &EventPublisher{C: make(chan es.EventDescriptor, 1)}
	eventStore := CreateEventStore(publisher)
	accountRepository := NewAccountRepository(eventStore)
	accountCommandHandler := NewAccountCommandHandler(accountRepository)
	banking := es.CreateBoundedContext(accountCommandHandler, &ExternalEventHandler{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := banking.Listen(ctx); err != nil {
			panic(err)
		}
	}()

	for x := 1; x <= 10; x++ {
		var (
			accountOpenedEvent AccountOpenedEvent
			accountNumber      AccountNumber
			desc               es.EventDescriptor
		)

		{
			// We issue our initial command for opening an account
			banking.Commands() <- OpenAccountCommand{
				Username: fmt.Sprintf("Paul %d", x),
			}
			// As a consequence an AccountOpenedEvent is pusblished as soon as it
			// was persisted at the storage
			desc := <-publisher.C
			accountOpenedEvent = desc.Event.(AccountOpenedEvent)
			accountNumber = accountOpenedEvent.Number
		}
		{
			banking.Commands() <- DepositMoneyCommand{
				AccountNumber:   accountNumber,
				Amount:          50,
				OriginalVersion: desc.Version,
			}
			desc = <-publisher.C
		}
		{
			banking.Events() <- UnpaidDebtsThroughCreditorsReportedEvent{
				Number: accountNumber,
			}
			desc = <-publisher.C
		}
		{
			banking.Commands() <- WithdrawMoneyCommand{
				AccountNumber:   accountNumber,
				Amount:          50,
				OriginalVersion: desc.Version,
			}
			// desc = <-publisher.C
		}
		{
			banking.Commands() <- UnfreezeAccountCommand{
				AccountNumber:   accountNumber,
				OriginalVersion: desc.Version,
			}
			desc = <-publisher.C
		}
		{
			banking.Commands() <- DepositMoneyCommand{
				AccountNumber:   accountNumber,
				Amount:          150,
				OriginalVersion: desc.Version,
			}
			desc = <-publisher.C
		}
		{
			banking.Commands() <- CloseAccountCommand{
				AccountNumber:   accountNumber,
				OriginalVersion: desc.Version,
			}
			// desc = <-publisher.C
		}
		{
			banking.Commands() <- WithdrawMoneyCommand{
				AccountNumber:   accountNumber,
				Amount:          200,
				OriginalVersion: desc.Version,
			}
			desc = <-publisher.C
		}
		{
			banking.Commands() <- CloseAccountCommand{
				AccountNumber:   accountNumber,
				OriginalVersion: desc.Version,
			}
			desc = <-publisher.C
		}
		time.Sleep(time.Second)
		fmt.Println(".")
	}

	<-sigs
	cancel()
	wg.Wait()
}
