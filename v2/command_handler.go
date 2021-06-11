package eventsourcerer

// The CommandHandler processes the Commands. After checking the validity of the
// state-transition it applies the Command to the Aggregate, which itself creates a
// transient event from it.
type CommandHandler interface {
	Handle(Command) error
}
