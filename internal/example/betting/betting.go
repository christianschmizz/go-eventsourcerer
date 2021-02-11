package main


import (
	"fmt"
	"time"
)

type BookingStatus string

const (
	Unspecified BookingStatus = "unspecified"
	Booked      BookingStatus = "booked"
	Unbooked    BookingStatus = "unbooked"
)

type Command struct {
}

type BookMatchCommand struct {
	ID string
	Command
}

// Aggegrat
type Booking struct {
	ID     string
	Status BookingStatus
}




func NewMatchBookedEvent(ID string) *MatchBookedEvent {
	return &MatchBookedEvent{
		ID: ID,
	}
}

func NewMatchUnbookedEvent(ID string) *MatchUnbookedEvent {
	return &MatchUnbookedEvent{
		ID: ID,
	}
}

type MatchBookings map[string]Booking

func (b MatchBookings) Apply(ev interface{}) error {
	switch v := ev.(type) {
	case *MatchBookedEvent:
		booking, exists := b[v.ID]
		if !exists {
			booking = Booking{ID: v.ID, Status: Unspecified}
		}
		booking.Status = Booked
		b[v.ID] = booking
	case *MatchUnbookedEvent:
		booking, exists := b[v.ID]
		if !exists {
			booking = Booking{ID: v.ID, Status: Unspecified}
		}
		booking.Status = Unbooked
		b[v.ID] = booking
	}
	return nil
}

type Aggregate struct {
	Bookings MatchBookings
	Journal  []interface{}
}

func (a *Aggregate) Apply(ev interface{}) {
	err := a.Bookings.Apply(ev)
	if err != nil {
		fmt.Println("no")
	}
	a.Journal = append(a.Journal, ev)
}

func (a *Aggregate) Replay(journal []interface{}) {
	for _, ev := range journal {
		a.Apply(ev)
	}
}

func NewAggregate() *Aggregate {
	return &Aggregate{
		Bookings: make(MatchBookings, 0),
		Journal:  make([]interface{}, 0, 10),
	}
}

func CommandHandler(cmdC <-chan interface{}, evC chan<- interface{}) {
	fmt.Println("command handler")
	for {
		select {
		case cmd, more := <-cmdC:
			if !more {
				fmt.Println("commands closed, so we do")
				return
			}

			switch v := cmd.(type) {
			case *BookMatchCommand:
				evC <- &MatchBookedEvent{ID: v.ID}
			}
		default:
		}
	}
}

func EventHandler(evC <-chan interface{}, aggregate *Aggregate) {
	fmt.Println("event handler")
	for {
		ev, more := <-evC
		if !more {
			fmt.Println("events closed")
			return
		}
		// apply event
		aggregate.Apply(ev)
	}
}

func main() {
	evC := make(chan interface{})
	cmdC := make(chan interface{})

	aggregate := NewAggregate()

	go EventHandler(evC, aggregate)
	go CommandHandler(cmdC, evC)

	cmdC <- &BookMatchCommand{ID: "YAY"}

	time.Sleep(5 * time.Second)

	//	aggregate.Apply(NewMatchBookedEvent("ABC"))
	// aggregate.Apply(NewMatchBookedEvent("DEF"))
	//	fmt.Printf("%+v\n", aggregate)

	//	aggregate.Apply(NewMatchUnbookedEvent("DEF"))
	//	aggregate.Apply(NewMatchBookedEvent("DEF"))
	//	aggregate.Apply(NewMatchUnbookedEvent("DEF"))
	//	aggregate.Apply(NewMatchBookedEvent("DEF"))
	fmt.Printf("%+v\n", aggregate)

	//	aggregate2 := NewAggregate()
	//	aggregate2.Replay(aggregate.Journal)
	//	fmt.Printf("%+v\n", aggregate2)

	// fmt.Printf("len=%d cap=%d %v\n", len(journal), cap(journal), journal)

}
