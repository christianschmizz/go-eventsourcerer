package eventsourcerer

import (
	"fmt"
	"reflect"
)

type AggregateDoesNotExistError struct {
	ID AggregateID
}

func (e AggregateDoesNotExistError) Error() string {
	return fmt.Sprintf("aggregate does not exist: %d", int64(e.ID))
}

func NewAggregateDoesNotExistError(id AggregateID) *AggregateDoesNotExistError {
	return &AggregateDoesNotExistError{ID: id}
}

// A ConcurrencyError signalizes a version mismatch
type ConcurrencyError struct {
	ExpectedVersion AggregateVersion
	CurrentVersion  AggregateVersion
}

func NewConcurrencyError(expectedVersion, currentVersion AggregateVersion) ConcurrencyError {
	return ConcurrencyError{
		ExpectedVersion: expectedVersion,
		CurrentVersion:  currentVersion,
	}
}

func (e ConcurrencyError) Error() string {
	return fmt.Sprintf("version did not match. expected version: %d got: %d", e.ExpectedVersion, e.CurrentVersion)
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

type UnknownEventError struct {
	Event Event
}

func (e UnknownEventError) Error() string {
	return fmt.Sprintf("%s: ", reflect.TypeOf(e.Event).String())
}

func NewUnknownEventError(ev Event) UnknownEventError {
	return UnknownEventError{
		Event: ev,
	}
}
