package eventsourcerer

import (
	"context"

	"github.com/rs/zerolog/log"
)

type boundedContext struct {
	commandCh      chan Command
	commandHandler CommandHandler
	eventCh        chan Event
	eventHandler   EventHandler
}

func CreateBoundedContext(commandHandler CommandHandler, eventHandler EventHandler) *boundedContext {
	return &boundedContext{
		commandCh:      make(chan Command, 10),
		commandHandler: commandHandler,
		eventCh:        make(chan Event, 10),
		eventHandler:   eventHandler,
	}
}

func (c *boundedContext) Listen(ctx context.Context) error {
	defer close(c.eventCh)
	defer close(c.commandCh)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case cmd, more := <-c.commandCh:
			if !more {
				log.Debug().Msg("commands closed, so we do")
				return nil
			}
			log.Debug().Msgf("received command: %#v\n", cmd)
			if err := c.commandHandler.Handle(cmd); err != nil {
				log.Error().Err(err).Msg("command failed")
			}
		case ev, more := <-c.eventCh:
			if !more {
				log.Debug().Msg("events closed, so we do")
				return nil
			}
			log.Debug().Msgf("received event: %#v\n", ev)
			c.eventHandler.Handle(ev, c.commandCh)
		}
	}
	return nil
}

func (c *boundedContext) Commands() chan<- Command {
	return c.commandCh
}

func (c *boundedContext) Events() chan<- Event {
	return c.eventCh
}