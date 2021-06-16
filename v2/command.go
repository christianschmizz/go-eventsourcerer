package eventsourcerer

type Command interface {
	Reject(err error) error
}

// The CommandBase is the foundation for every command
type CommandBase struct {
}

func (bc CommandBase) Reject(err error) error {
	return NewCommandRejectedError(bc, err)
}
