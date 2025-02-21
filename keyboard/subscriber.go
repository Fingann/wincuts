package keyboard 

import (
	"sync"
)

// SubscriberManager manages subscribers and broadcasting events to them.
type SubscriberManager[T any] struct {
	subscribers map[chan T]struct{}
	mu          sync.Mutex
}

// NewSubscriberManager creates a new SubscriberManager.
func NewSubscriberManager[T any]() *SubscriberManager[T] {
	return &SubscriberManager[T]{
		subscribers: make(map[chan T]struct{}),
		mu:          sync.Mutex{},
	}
}

// AddSubscriber adds a new subscriber and returns the channel.
func (sm *SubscriberManager[T]) AddSubscriber() chan T {
	ch := make(chan T, 100)
	sm.mu.Lock()
	sm.subscribers[ch] = struct{}{}
	sm.mu.Unlock()
	return ch
}

// RemoveSubscriber removes a subscriber and closes its channel.
func (sm *SubscriberManager[T]) RemoveSubscriber(ch chan T) {
	sm.mu.Lock()
	if _, ok := sm.subscribers[ch]; ok {
		close(ch)
		delete(sm.subscribers, ch)
	}
	sm.mu.Unlock()
}

// Broadcast sends a message to all subscribers, removing those that are closed.
func (sm *SubscriberManager[T]) Broadcast(msg T) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for ch := range sm.subscribers {
		select {
		case ch <- msg:
		default:
			// If the channel is closed or blocked, remove it
			close(ch)
			delete(sm.subscribers, ch)
		}
	}
}

// CloseAll closes all subscriber channels and clears the subscribers map.
func (sm *SubscriberManager[T]) CloseAll() {
	sm.mu.Lock()
	for ch := range sm.subscribers {
		close(ch)
		delete(sm.subscribers, ch)
	}
	sm.mu.Unlock()
}
