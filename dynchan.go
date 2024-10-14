package dynchan

// Chan is a dynamic channel.
// It is like a Go channel but with a dynamically sized buffer.
// (Normal Go channels have a fixed-size buffer.)
// Because of this, sends never block.
// To send items to the Chan, use its Send field.
// To receive items from the Chan, use Recv.
// To signal the end of input, close the Send channel.
type Chan[T any] struct {
	Send chan<- T
	Recv <-chan T

	cancel chan struct{}
}

// New creates a new [Chan].
// This is equivalent to calling [NewWithBuffer] with a FIFO queue from [NewFifo].
// You must close the Chan's Send channel when you are done sending items.
// You must also call the Chan's Close method when you are done receiving items.
func New[T any]() Chan[T] {
	return NewWithBuffer[T](NewFifo[T]())
}

// NewWithBuffer creates a new [Chan] with a given buffer.
// This can be [NewFifo] for normal channel semantics,
// [NewHeap] for priority-queue semantics,
// or any type implementing the [Buffer] interface.
// You must close the Chan's Send channel when you are done sending items.
// You must also call the Chan's Close method when you are done receiving items.
func NewWithBuffer[T any](b Buffer[T]) Chan[T] {
	in, out, cancel := make(chan T), make(chan T), make(chan struct{})

	go func() {
		for val := range in {
			b.Enqueue(val)
		}
		b.Close()
	}()

	go func() {
		defer close(out)

		for {
			val, ok := b.Dequeue()
			if !ok {
				return
			}

			select {
			case out <- val:
			case <-cancel:
				return
			}
		}
	}()

	return Chan[T]{Send: in, Recv: out, cancel: cancel}
}

// Close closes the receiving end of the Chan.
func (dc Chan[T]) Close() {
	if dc.cancel == nil {
		// Make this method idempotent.
		return
	}
	close(dc.cancel)
	dc.cancel = nil
}

// Buffer is a dynamically sized buffer for use in a dynamic channel.
type Buffer[T any] interface {
	// Enqueue adds an item to the buffer.
	// It must be safe for concurrent use by multiple goroutines.
	Enqueue(T)

	// Dequeue removes and returns an item from the buffer, and true.
	// If the buffer is empty and closed, Dequeue returns the zero value and false.
	// If the buffer is empty but not closed,
	// Dequeue blocks until an item is added to the buffer or the buffer is closed.
	// It must be safe for concurrent use by multiple goroutines.
	Dequeue() (T, bool)

	// Close closes the buffer,
	// signaling that no more items will be added to it.
	Close()
}
