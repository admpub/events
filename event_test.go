package events_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/admpub/events"
)

var _ = Describe("Event", func() {
	It("should create event object", func() {
		Expect(events.New("test")).To(BeEquivalentTo(events.Event{Key: "test", Context: events.Map{}}))
	})

	It("should subscribe callback to event", func() {
		emitter := events.NewEmitter()
		var i int
		emitter.On("test", events.Callback(func(event events.Event) error {
			i++
			return nil
		}))
		err := emitter.Fire("test")
		Expect(err).To(BeNil())
		Expect(i).To(BeEquivalentTo(1))

		emitter.On("test", events.Callback(func(event events.Event) error {
			i++
			return nil
		}))
		err = emitter.Fire("test")
		Expect(err).To(BeNil())
		Expect(i).To(BeEquivalentTo(3))

		emitter.On("test", events.Callback(func(event events.Event) error {
			return errors.New(`testError`)
		}))
		err = emitter.Fire("test")
		Expect(err).To(MatchError(`testError`))
	})
})
