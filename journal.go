package eventsourcerer

type Journal chan Event

func (j Journal) SliceAndDrain() []Event {
	s := make([]Event, 0)
	for i := range j {
		s = append(s, i)
	}
	return s
}
