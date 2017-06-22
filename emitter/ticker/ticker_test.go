package ticker

import (
	"testing"
	"time"

	"github.com/admpub/events"
	"github.com/admpub/events/emitter"
)

type TestDispatcher struct {
	Count  int
	Target bool
}

func (handler *TestDispatcher) AddSubscribers(t ...events.Listener) {
}

func (handler *TestDispatcher) Dispatch(event events.Event) error {
	handler.Count++
	handler.Target = true
	return nil
}

func TestNewPeriodicEmitter(t *testing.T) {
	ticker := New(emitter.New())

	if len(ticker.timers) != 1 {
		t.Log("Fail ticker create - no timers storage")
		t.Fail()
	}
}

func TestRegisterEvent(t *testing.T) {
	ticker := New(emitter.New())
	dispatcher := new(TestDispatcher)

	if dispatcher.Target != false {
		t.Fail()
	}

	ticker.Emitter.(*emitter.Emitter).Dispatchers["test"] = dispatcher
	ticker.RegisterEvent("test", 1*time.Millisecond)
	time.Sleep(4 * time.Millisecond)

	if dispatcher.Target != true {
		t.Log("fail event fire")
		t.Log(ticker.events, ticker.timers)
		t.Fail()
	}

	if dispatcher.Count < 2 {
		t.Log("fail event fire - not enough events")
		t.Fail()
	}
}

func TestRemoveEvent(t *testing.T) {
	ticker := New(emitter.New())
	dispatcher := new(TestDispatcher)

	if dispatcher.Target != false {
		t.Fail()
	}

	ticker.Emitter.(*emitter.Emitter).Dispatchers["test"] = dispatcher
	ticker.RegisterEvent("test", 5*time.Millisecond)
	time.Sleep(11 * time.Millisecond)
	ticker.RemoveEvent("test")
	time.Sleep(11 * time.Millisecond)

	if !dispatcher.Target {
		t.Log("fail event fire")
		t.Fail()
	}

	if dispatcher.Count < 2 {
		t.Log("fail event fire", dispatcher.Count)
		t.Fail()
	}
}
