package events

import (
	"github.com/admpub/events/meta"
)

type Emitter interface {
	On(string, ...Listener) Emitter //AddEventListener
	Off(string) Emitter             //RemoveEventListeners
	Fire(interface{}, ...meta.Map)
}

type Dispatcher interface {
	AddSubscribers(...Listener)
	Dispatch(Event)
}

type DispatcherFactory func() Dispatcher

type Listener interface {
	Handle(Event)
}

type Stream chan Event

func (stream Stream) Handle(event Event) {
	stream <- event
}

type Callback func(Event)

func (callback Callback) Handle(event Event) {
	callback(event)
}

func New(name string) Event {
	return Event{name, meta.Map{}}
}

type Event struct {
	Key     string
	Context meta.Map
}

func (event *Event) String() string {
	return event.Key
}
