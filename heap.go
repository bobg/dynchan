package dynchan

import (
	"cmp"
	"sync"
)

type heap[T any] struct {
	mu     sync.Mutex
	items  []T
	closed bool
	cond   *sync.Cond
	less   func(T, T) bool
}

// NewHeap creates a new heap, a.k.a. a priority queue,
// for use as a buffer in a dynamic channel.
// It is equivalent to calling [NewHeapFunc] with [cmp.Less].
func NewHeap[T cmp.Ordered]() Buffer[T] {
	return NewHeapFunc[T](cmp.Less)
}

// NewHeapFunc creates a new heap, a.k.a. a priority queue,
// for use as a buffer in a dynamic channel.
// The less function is used to compare elements in the heap,
// and the top of the heap is always the "least" element.
func NewHeapFunc[T any](less func(T, T) bool) Buffer[T] {
	h := heap[T]{less: less}
	h.cond = sync.NewCond(&h.mu)
	return &h
}

func (h *heap[T]) Enqueue(val T) {
	h.mu.Lock()
	h.items = append(h.items, val)
	h.bubbleUp()
	h.cond.Broadcast()
	h.mu.Unlock()
}

func (h *heap[T]) Dequeue() (T, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for len(h.items) == 0 && !h.closed {
		h.cond.Wait()
	}

	var (
		ok    = len(h.items) > 0
		item0 T
	)
	if ok {
		item0 = h.items[0]
		h.items[0] = h.items[len(h.items)-1]
		h.items = h.items[:len(h.items)-1]
		h.bubbleDown()
	}
	return item0, ok
}

func (h *heap[T]) Close() {
	h.mu.Lock()
	h.closed = true
	h.cond.Broadcast()
	h.mu.Unlock()
}

func (h *heap[T]) bubbleUp() {
	i := len(h.items) - 1
	for i > 0 {
		parent := (i - 1) / 2
		if h.less(h.items[parent], h.items[i]) {
			break
		}
		h.items[parent], h.items[i] = h.items[i], h.items[parent]
		i = parent
	}
}

func (h *heap[T]) bubbleDown() {
	i := 0
	for {
		left := 2*i + 1
		if left >= len(h.items) {
			break
		}

		j, right := left, left+1
		if right < len(h.items) && h.less(h.items[right], h.items[left]) {
			j = right
		}

		if h.less(h.items[i], h.items[j]) {
			break
		}

		h.items[i], h.items[j] = h.items[j], h.items[i]
		i = j
	}
}
