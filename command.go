package eventsourcerer

import (
	"fmt"
	"reflect"
)

type Command interface {
	Reject(err error) error
}

// The CommandHandler processes the Command. After checking the validity of the state
// transition it applies the the Command to the Aggregate, which itself creates a
// transient event from it.
type CommandHandler interface {
	Handle(Command) error
}

type CommandRejectedError struct {
	Command Command
	Err     error
}

func (e CommandRejectedError) Error() string {
	return fmt.Sprintf("%s: ", reflect.TypeOf(e.Command).String()) + e.Err.Error()
}

func NewCommandRejectedError(cmd Command, err error) CommandRejectedError {
	return CommandRejectedError{
		Command: cmd,
		Err:     err,
	}
}

// The BaseCommand is the foundation for every command
type BaseCommand struct {
}

// Reject issues a CommandRejectedError
func (bc BaseCommand) Reject(err error) error {
	return NewCommandRejectedError(bc, err)
}

// The BaseCommandHandler is the foundation for a command handler
type BaseCommandHandler struct {
}

func (h *BaseCommandHandler) Reject(cmd Command, err error) {
	fmt.Printf("REJECT: %+v %s\n", cmd, err)
}
