package events

// Broadcast event to all handlers
func Broadcast(event IEvent, handlers map[Listener]struct{}) (err error) {
	for handler := range handlers {
		if err = handler.Handle(event); err != nil || event.Aborted() {
			return
		}
	}
	return
}
