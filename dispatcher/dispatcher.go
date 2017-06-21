package dispatcher

import (
	"github.com/admpub/events"
)

type Type int

const (
	Sync Type = iota
	Cond
	Async
)

func New(typ Type) events.Dispatcher {
	switch typ {
	case Sync:
		return BroadcastFactory()
	case Cond:
		return ConditionalParallelBroadcastFactory()
	default:
		return ParallelBroadcastFactory()
	}
}

func Factory(typ Type) events.DispatcherFactory {
	switch typ {
	case Sync:
		return BroadcastFactory
	case Cond:
		return ConditionalParallelBroadcastFactory
	default:
		return ParallelBroadcastFactory
	}
}
