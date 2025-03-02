package events

import "errors"

var ErrAborted = errors.New(`aborted`)

// Broadcast event to all handlers
func Broadcast[V any](event V, handlers map[Listener[V]]struct{}) (err error) {
	for handler := range handlers {
		if err = handler.Handle(event); err != nil {
			if errors.Is(err, ErrAborted) {
				err = nil
			}
			return
		}
	}
	return
}
