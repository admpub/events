package events

import (
	"fmt"

	"github.com/webx-top/echo/param"
)

// EventOption for event
type EventOption struct {
	apply func(*Event)
}

// WithContext sets event metadata
func WithContext(context Map) EventOption {
	return EventOption{func(event *Event) {
		for key, value := range context {
			event.Context.Set(key, value)
		}
	}}
}

// New create new event with provided name and options
func New(data interface{}, options ...EventOption) *Event {
	var event Event

	switch value := data.(type) {
	case string:
		event = Event{Key: value, Context: Map{}}
	case Event:
		event = value
	}

	for _, option := range options {
		option.apply(&event)
	}

	return &event
}

func NewEvent[V any](key string, ctx V) *XEvent[V] {
	return &XEvent[V]{
		Key:     key,
		Context: ctx,
	}
}

type IEvent interface {
	String() string
	Abort() IEvent
	Aborted() bool
}

// Event
type XEvent[V any] struct {
	Key     string
	Context V
	aborted bool
}

func (event *XEvent[V]) String() string {
	return event.Key
}

func (event *XEvent[V]) Abort() IEvent {
	event.aborted = true
	return event
}

func (event *XEvent[V]) Aborted() bool {
	return event.aborted
}

// Event
type Event struct {
	Key     string
	Context param.Store
	aborted bool
}

func (event *Event) String() string {
	return event.Key
}

func (event *Event) Abort() IEvent {
	event.aborted = true
	return event
}

func (event *Event) Aborted() bool {
	return event.aborted
}

func ToMap(key string, value interface{}, args ...interface{}) Map {
	context := Map{key: value}
	for i, j := 0, len(args); i < j; i++ {
		if i%2 == 0 {
			key = fmt.Sprint(args[i])
			continue
		}
		context[key] = args[i]
	}
	return context
}
