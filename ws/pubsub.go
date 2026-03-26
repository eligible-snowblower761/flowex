package ws

import "sync"

// PubSub is a generic, type-safe fan-out broadcaster.
// Publishers send values; each subscriber gets its own buffered channel.
// Slow consumers are dropped (non-blocking publish).
type PubSub[T any] struct {
	mu     sync.RWMutex
	next   uint64
	subs   map[uint64]chan T
	closed bool
}

// NewPubSub creates a new PubSub broadcaster.
func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{subs: make(map[uint64]chan T)}
}

// Subscribe returns a channel that receives published values and a cancel
// function to unsubscribe. The buffer size controls backpressure tolerance.
func (ps *PubSub[T]) Subscribe(buf int) (<-chan T, func()) {
	ch := make(chan T, buf)

	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.closed {
		close(ch)
		return ch, func() {}
	}

	id := ps.next
	ps.next++
	ps.subs[id] = ch

	cancel := func() {
		ps.mu.Lock()
		c, ok := ps.subs[id]
		if ok {
			delete(ps.subs, id)
			close(c)
		}
		ps.mu.Unlock()
	}

	return ch, cancel
}

// Publish sends a value to all subscribers. Non-blocking: if a subscriber's
// channel is full the message is silently dropped for that subscriber.
func (ps *PubSub[T]) Publish(v T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if ps.closed || len(ps.subs) == 0 {
		return
	}
	for _, ch := range ps.subs {
		select {
		case ch <- v:
		default:
		}
	}
}

// Close shuts down the PubSub, closing all subscriber channels.
func (ps *PubSub[T]) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.closed {
		return
	}
	ps.closed = true
	for id, ch := range ps.subs {
		delete(ps.subs, id)
		close(ch)
	}
}
