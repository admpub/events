package events

import (
	"fmt"

	"github.com/admpub/events/meta"
)

const (
	//ModeAsync async
	ModeAsync = iota

	//ModeSync sync
	ModeSync

	//ModeWait async & sync.Wait
	ModeWait
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
	return Event{
		Key:     name,
		Context: meta.Map{},
	}
}

type Event struct {
	Key     string
	Context meta.Map
	aborted bool
}

func (event *Event) String() string {
	return event.Key
}

func (event *Event) Abort() *Event {
	event.aborted = true
	return event
}

func (event *Event) Aborted() bool {
	return event.aborted
}

func ToMap(key string, value interface{}, args ...interface{}) meta.Map {
	context := meta.Map{key: value}
	for i, j := 0, len(args); i < j; i += 2 {
		if i%2 == 0 {
			key = fmt.Sprint(args[i])
			continue
		}
		context[key] = args[i]
	}
	return context
}

func Map(context meta.Map, mode ...int) meta.Map {
	if len(mode) > 0 {
		switch mode[0] {
		case ModeSync:
			context["_sync"] = struct{}{}
		case ModeWait:
			context["_wait"] = struct{}{}
		}
	}
	return context
}
