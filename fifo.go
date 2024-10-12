package dynchan

import "sync"

type fifo[T any] struct {
	mu     sync.Mutex
	items  []T
	closed bool
	cond   *sync.Cond
}

// NewFifo creates a new FIFO queue for use as a buffer in a dynamic channel.
func NewFifo[T any]() Buffer[T] {
	var f fifo[T]
	f.cond = sync.NewCond(&f.mu)
	return &f
}

func (f *fifo[T]) Enqueue(val T) {
	f.mu.Lock()
	f.items = append(f.items, val)
	f.cond.Broadcast()
	f.mu.Unlock()
}

func (f *fifo[T]) Dequeue() (T, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for len(f.items) == 0 && !f.closed {
		f.cond.Wait()
	}

	var (
		ok    = len(f.items) > 0
		item0 T
	)
	if ok {
		item0, f.items = f.items[0], f.items[1:]
	}

	return item0, ok
}

func (f *fifo[T]) Close() {
	f.mu.Lock()
	f.closed = true
	f.cond.Broadcast()
	f.mu.Unlock()
}
