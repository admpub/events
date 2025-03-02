package events_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/admpub/events"
)

func TestTicker(t *testing.T) {
	tick := events.NewTicker(events.NewEmitter[any]())
	defer tick.Stop()
	tick.RegisterEvent("periodic.event.1", 2*time.Second, events.Callback(func(i any) error {
		fmt.Println(time.Now())
		return nil
	}))
	time.Sleep(10 * time.Second)
}
