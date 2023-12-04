package cp

import "sync"

type Subscriber interface{}

type Event interface {
	// Dispatch(name string, data any)
	// Subscribe(name string, fn func())
}

type event struct {
}

type eventBus struct {
	subscribers map[string][]Subscriber
	m           sync.RWMutex
}

// func (b *eventBus) Dispatch(name string, data any) {
// 	subscribers, ok := b.subscribers[name]
// 	if !ok {
// 		return
// 	}
// 	for _, s := range subscribers {
//
// 	}
// }
//
// func (b *eventBus) Subscribe(name string, fn func(data any)) {
//
// }
