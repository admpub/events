package events

import "github.com/webx-top/echo/param"

type Map = param.Store

// DispatchStrategy defines strategy of delivery event to handlers
type DispatchStrategy func(IEvent, map[Listener]struct{}) error

// Listener defines event handler interface
type Listener interface {
	Handle(IEvent) error
}

type ID interface {
	ID() string
}

// Stream implements Listener interface on channel
type Stream[V any] chan V

// Handle Listener
func (stream Stream[V]) Handle(event IEvent) error {
	stream <- event.(V)
	return nil
}

type Streamer[V any] interface {
	Listener
	ID
	Chan() <-chan V
}

func StreamWithID[V any](ch chan V, id string) Streamer[V] {
	return &stream[V]{
		ch: ch,
		id: id,
	}
}

// Stream implements Listener interface on channel
type stream[V any] struct {
	ch chan V
	id string
}

// Handle Listener
func (s *stream[V]) Handle(event IEvent) error {
	s.ch <- event.(V)
	return nil
}

func (s *stream[V]) ID() string {
	return s.id
}

func (s *stream[V]) Chan() <-chan V {
	return s.ch
}

// Callback implements Listener interface on function
func Callback[V any](function func(V) error, id ...string) Listener {
	var _id string
	if len(id) > 0 {
		_id = id[0]
	}
	return callback[V]{function: &function, id: _id}
}

type callback[V any] struct {
	function *func(V) error
	id       string
}

// Handle Listener
func (c callback[V]) Handle(event IEvent) error {
	return (*c.function)(event.(V))
}

// ID Listener ID
func (c callback[V]) ID() string {
	return c.id
}

func WithID(l Listener, id string) Listener {
	return &listenerWithID{Listener: l, id: id}
}

type listenerWithID struct {
	Listener
	id string
}

func (l listenerWithID) ID() string {
	return l.id
}
