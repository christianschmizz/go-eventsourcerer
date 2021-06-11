package main

import (
	es "github.com/christianschmizz/go-eventsourcerer"
	"github.com/rs/zerolog/log"
)

type SuspiciousOrIllegalActivityDetectedEvent struct {
	es.BaseEvent
	Number AccountNumber
}

type UnpaidDebtsThroughCreditorsReportedEvent struct {
	es.BaseEvent
	Number AccountNumber
}

type ExternalEventHandler struct{}

func (h *ExternalEventHandler) Handle(ev es.Event, c chan<- es.Command) {
	log.Debug().Msgf("E: %#v", ev)
	switch e := ev.(type) {
	case SuspiciousOrIllegalActivityDetectedEvent:
	case UnpaidDebtsThroughCreditorsReportedEvent:
		// send FreezeAccountCommand
		c <- FreezeAccountCommand{
			AccountNumber:   e.Number,
			OriginalVersion: es.IgnoreVersion,
		}
	}
}
