package events

import (
	"sort"
	"sync"
)

var Default = NewEmitter[Event]()

// EmitterOption defines option for Emitter
type EmitterOption[V any] struct {
	apply func(*Emitter[V])
}

// WithEventStategy sets delivery strategy for provided event
func WithEventStategy[V any](event string, strategy DispatchStrategy[V]) EmitterOption[V] {
	return EmitterOption[V]{func(emitter *Emitter[V]) {
		if dispatcher, exists := emitter.dispatchers[event]; exists {
			dispatcher.strategy = strategy
			return
		}

		emitter.dispatchers[event] = NewDispatcher(strategy)
	}}
}

// WithDefaultStrategy sets default delivery strategy for event emitter
func WithDefaultStrategy[V any](strategy DispatchStrategy[V]) EmitterOption[V] {
	return EmitterOption[V]{func(emitter *Emitter[V]) {
		emitter.strategy = strategy
	}}
}

// NewEmitter creates new event emitter
func NewEmitter[V any](options ...EmitterOption[V]) *Emitter[V] {
	emitter := new(Emitter[V])
	emitter.strategy = Broadcast[V]
	emitter.dispatchers = make(map[string]*Dispatcher[V])

	for _, option := range options {
		option.apply(emitter)
	}

	return emitter
}

// Emitter
type Emitter[V any] struct {
	guard       sync.Mutex
	strategy    DispatchStrategy[V]
	dispatchers map[string]*Dispatcher[V]
}

type Emitterer[V any] interface {
	On(string, ...Listener[V]) Emitterer[V]
	AddEventListener(handler Listener[V], events ...string)
	Off(string) Emitterer[V]
	RemoveEventListener(handler Listener[V])
	Fire(string, V) error
	EventNames() []string
	HasEvent(string) bool
}

// On subscribes listeners to provided event and return emitter
// usefull for chain subscriptions
func (emitter *Emitter[V]) On(event string, handlers ...Listener[V]) Emitterer[V] {
	emitter.AddEventListeners(event, handlers...)
	return emitter
}

// AddEventListeners subscribes listeners to provided event
func (emitter *Emitter[V]) AddEventListeners(event string, handlers ...Listener[V]) {
	emitter.guard.Lock()

	if _, exists := emitter.dispatchers[event]; !exists {
		emitter.dispatchers[event] = NewDispatcher(emitter.strategy)
	}
	emitter.dispatchers[event].AddSubscribers(handlers)

	emitter.guard.Unlock()
}

// AddEventListener subscribes listeners to provided events
func (emitter *Emitter[V]) AddEventListener(handler Listener[V], events ...string) {
	emitter.guard.Lock()
	for _, event := range events {
		if _, exists := emitter.dispatchers[event]; !exists {
			emitter.dispatchers[event] = NewDispatcher(emitter.strategy)
		}

		emitter.dispatchers[event].AddSubscriber(handler)
	}
	emitter.guard.Unlock()
}

// Off unsubscribe all listeners from provided event
func (emitter *Emitter[V]) Off(event string) Emitterer[V] {
	emitter.RemoveEventListeners(event)
	return emitter
}

// RemoveEventListeners unsubscribe all listeners from provided event
func (emitter *Emitter[V]) RemoveEventListeners(event string) {
	emitter.guard.Lock()
	delete(emitter.dispatchers, event)
	emitter.guard.Unlock()
}

// RemoveEventListener unsubscribe provided listener from all events
func (emitter *Emitter[V]) RemoveEventListener(handler Listener[V]) {
	emitter.guard.Lock()
	for _, dispatcher := range emitter.dispatchers {
		dispatcher.RemoveSubscriber(handler)
	}
	emitter.guard.Unlock()
}

// Fire start delivering event to listeners
func (emitter *Emitter[V]) Fire(name string, data V) (err error) {
	if dispatcher, ok := emitter.dispatchers[name]; ok {
		err = dispatcher.Dispatch(data)
	}
	return
}

// EventNames ...
func (emitter *Emitter[V]) EventNames() []string {
	names := make([]string, len(emitter.dispatchers))
	var i int
	for name := range emitter.dispatchers {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return names
}

// HasEvent ...
func (emitter *Emitter[V]) HasEvent(event string) bool {
	_, ok := emitter.dispatchers[event]
	return ok
}
