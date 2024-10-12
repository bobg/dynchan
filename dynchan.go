package dynchan

// Chan is like a Go channel but with a dynamically sized buffer.
// (Normal Go channels have a fixed-size buffer.)
// Because of this, sends never block.
// To send items to the Chan, use its Send field.
// To receive items from the Chan, use Recv.
// To signal the end of input, close the Send channel.
type Chan[T any] struct {
	Send chan<- T
	Recv <-chan T
}

// New creates a new Chan.
// This is equivalent to calling [NewWithBuffer] with a FIFO queue from [NewFifo].
func New[T any]() Chan[T] {
	return NewWithBuffer[T](NewFifo[T]())
}

// NewWithBuffer creates a new Chan with a given buffer.
// This can be [NewFifo] for normal channel semantics,
// [NewHeap] for priority-queue semantics,
// or any type implementing the [Buffer] interface.
func NewWithBuffer[T any](b Buffer[T]) Chan[T] {
	in, out := make(chan T), make(chan T)

	go func() {
		for val := range in {
			b.Enqueue(val)
		}
		b.Close()
	}()

	go func() {
		for {
			val, ok := b.Dequeue()
			if !ok {
				close(out)
				return
			}

			out <- val
		}
	}()

	return Chan[T]{Send: in, Recv: out}
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
