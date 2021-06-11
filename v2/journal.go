package eventsourcerer

// A Journal is used as a sequence of domain events
type Journal chan Event

// SliceAndDrain empties all events at the journal and returns them as a slice
func (j Journal) SliceAndDrain() []Event {
	s := make([]Event, 0)
	for i := range j {
		s = append(s, i)
	}
	return s
}
