package events

import "github.com/webx-top/echo/param"

type Map = param.Store

// DispatchStrategy defines strategy of delivery event to handlers
type DispatchStrategy func(Event, map[Listener]struct{}) error

// Listener defines event handler interface
type Listener interface {
	Handle(Event) error
}

type ID interface {
	ID() string
}

// Stream implements Listener interface on channel
type Stream chan Event

// Handle Listener
func (stream Stream) Handle(event Event) error {
	stream <- event
	return nil
}

// Callback implements Listener interface on function
func Callback(function func(Event) error, id ...string) Listener {
	var _id string
	if len(id) > 0 {
		_id = id[0]
	}
	return callback{function: &function, id: _id}
}

type callback struct {
	function *func(Event) error
	id       string
}

// Handle Listener
func (callback callback) Handle(event Event) error {
	return (*callback.function)(event)
}

// ID Listener ID
func (callback callback) ID() string {
	return callback.id
}
