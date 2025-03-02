package events_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/admpub/events"
)

var _ = Describe("Event", func() {
	It("should create event object", func() {
		Expect(events.New("test")).To(BeEquivalentTo(&events.Event{Key: "test", Context: events.Map{}}))
	})

	It("should subscribe callback to event", func() {
		emitter := events.NewEmitter[any]()
		var i int
		emitter.On("test", events.Callback(func(event any) error {
			i++
			return nil
		}))
		err := emitter.Fire("test", nil)
		Expect(err).To(BeNil())
		Expect(i).To(BeEquivalentTo(1))

		emitter.On("test", events.Callback(func(event any) error {
			i++
			return nil
		}))
		err = emitter.Fire("test", nil)
		Expect(err).To(BeNil())
		Expect(i).To(BeEquivalentTo(3))

		var j int
		emitter.On("test", events.Callback(func(event any) error {
			j++
			return errors.New(`testError`)
		}, `testError`))
		err = emitter.Fire("test", nil)
		Expect(err).To(MatchError(`testError`))
		Expect(j).To(BeEquivalentTo(1))

		emitter.On("test", events.Callback(func(event any) error {
			j++
			return errors.New(`testError2`)
		}, `testError`))
		err = emitter.Fire("test", nil)
		Expect(err).To(MatchError(`testError2`))
		Expect(j).To(BeEquivalentTo(2))
	})
})
