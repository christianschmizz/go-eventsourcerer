package main

import (
	es "github.com/christianschmizz/go-eventsourcerer"
	"github.com/rs/zerolog/log"
)

type EventPublisher struct {
	C chan es.EventDescriptor
}

func (p *EventPublisher) Publish(event es.EventDescriptor) {
	log.Debug().Msgf("publishing %#v", event)
	p.C <- event
}
