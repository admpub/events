package events

var value = struct{}{}

type Dispatcherer[V any] interface {
	AddSubscriber(handler Listener[V])
	AddSubscribers([]Listener[V])
	RemoveSubscriber(handler Listener[V])
	Dispatch(Event) error
}

// NewDispatcher creates new dispatcher
func NewDispatcher[V any](strategy DispatchStrategy[V]) *Dispatcher[V] {
	dispatcher := new(Dispatcher[V])
	dispatcher.strategy = strategy
	dispatcher.idm = make(map[string]Listener[V])
	dispatcher.subscribers = make(map[Listener[V]]struct{})

	return dispatcher
}

// Dispatcher stores event listeners of concrete event
type Dispatcher[V any] struct {
	strategy    DispatchStrategy[V]
	idm         map[string]Listener[V]
	subscribers map[Listener[V]]struct{}
}

// AddSubscriber adds one listener
func (dispatcher *Dispatcher[V]) AddSubscriber(handler Listener[V]) {
	dispatcher.add(handler)
}

func (dispatcher *Dispatcher[V]) add(handler Listener[V]) {
	if idi, ok := handler.(ID); ok {
		id := idi.ID()
		if len(id) > 0 {
			if h, y := dispatcher.idm[id]; y {
				dispatcher.RemoveSubscriber(h)
			}
			dispatcher.idm[id] = handler
		}
	}
	dispatcher.subscribers[handler] = value
}

// AddSubscribers adds slice of listeners
func (dispatcher *Dispatcher[V]) AddSubscribers(handlers []Listener[V]) {
	for _, handler := range handlers {
		dispatcher.add(handler)
	}
}

// RemoveSubscriber removes listener
func (dispatcher *Dispatcher[V]) RemoveSubscriber(handler Listener[V]) {
	delete(dispatcher.subscribers, handler)
}

// Dispatch deliver event to listeners using strategy
func (dispatcher *Dispatcher[V]) Dispatch(event V) error {
	return dispatcher.strategy(event, dispatcher.subscribers)
}
