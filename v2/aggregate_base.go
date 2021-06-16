package eventsourcerer

type EventApplicator func(Event) error

// The AggregateBase implements the basic structures for handling events.
type AggregateBase struct {
	// Changelog keeps collection of all events which have not yet been
	// written (transient) to the event-store.
	*Changelog

	applicator EventApplicator
}

func CreateAggregateBase() *AggregateBase {
	return &AggregateBase{
		Changelog: NewChangelog(),
	}
}

func (a *AggregateBase) SetEventApplicator(applicator EventApplicator) {
	a.applicator = applicator
}

func (a *AggregateBase) applyEvents(events ...Event) error {
	for _, event := range events {
		if err := a.applicator(event); err != nil {
			return err
		}
	}
	return nil
}

func (a *AggregateBase) Replay(journal Journal) error {
	for event := range journal {
		if err := a.applyEvents(event); err != nil {
			return err
		}
	}
	return nil
}

func (a *AggregateBase) Apply(event Event) error {
	if err := a.applyEvents(event); err != nil {
		return err
	}
	a.TrackChanges(event)
	return nil
}
