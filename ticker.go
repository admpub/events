package events

import (
	"reflect"
	"sync"
	"time"
)

// NewTicker creates new PeriodicEmitter
func NewTicker[V any](emitter *Emitter[V]) *PeriodicEmitter[V] {
	actions := make(chan func())

	ticker := &PeriodicEmitter[V]{
		Emitter: emitter,
		actions: actions,
		events:  make(map[string]*time.Ticker),
		timers:  []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(actions)}},
		mapping: make(map[int]string),
	}

	go ticker.run()

	return ticker
}

// PeriodicEmitter is a source of periodic events
type PeriodicEmitter[V any] struct {
	*Emitter[V]

	stopOnce sync.Once
	actions  chan func()
	mapping  map[int]string
	events   map[string]*time.Ticker
	timers   []reflect.SelectCase
}

func (emitter *PeriodicEmitter[V]) Stop() {
	emitter.stopOnce.Do(func() {
		close(emitter.actions)
	})
}

// RegisterEvent registers new periodic event
func (emitter *PeriodicEmitter[V]) RegisterEvent(event string, period interface{}, handlers ...Listener[V]) {
	var timer *time.Ticker
	switch value := period.(type) {
	case time.Duration:
		timer = time.NewTicker(value)
	case *time.Ticker:
		timer = value
	default:
		return
	}

	emitter.actions <- func() {
		emitter.mapping[len(emitter.timers)] = event
		emitter.timers = append(emitter.timers, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timer.C)})
		emitter.events[event] = timer

		if len(handlers) > 0 {
			emitter.AddEventListeners(event, handlers...)
		}

		emitter.refresh()
	}
}

// RemoveEvent removes provided event
func (emitter *PeriodicEmitter[V]) RemoveEvent(event string) {
	emitter.actions <- func() {
		if emitter.remove(event) {
			emitter.refresh()
		}
	}
}

func (emitter *PeriodicEmitter[V]) remove(event string) bool {
	timer, exists := emitter.events[event]
	if !exists {
		return false
	}

	timer.Stop()
	delete(emitter.events, event)
	emitter.RemoveEventListeners(event)

	return true
}

func (emitter *PeriodicEmitter[V]) stop() {
	for event := range emitter.events {
		emitter.remove(event)
	}
}

func (emitter *PeriodicEmitter[V]) refresh() {
	emitter.timers = []reflect.SelectCase{emitter.timers[0]}

	for event, timer := range emitter.events {
		emitter.mapping[len(emitter.timers)] = event
		emitter.timers = append(emitter.timers, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(timer.C),
		})
	}
}

func (emitter *PeriodicEmitter[V]) run() {
	for {
		index, value, opened := reflect.Select(emitter.timers)

		switch index {
		case 0:
			if opened {
				value.Call(nil)
			} else {
				emitter.stop()
				return
			}
		default:
			if event, exists := emitter.mapping[index]; exists {
				if opened {
					emitter.Fire(event, value.Interface().(V))
				} else {
					delete(emitter.events, event)
					emitter.refresh()
				}
			}
		}
	}
}
