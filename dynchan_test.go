package dynchan

import (
	"testing"
	"time"
)

func TestFIFO(t *testing.T) {
	dc := New[int]()
	defer dc.Close()

	dc.Send <- 1
	dc.Send <- 2

	if got, ok := <-dc.Recv; !ok || got != 1 {
		t.Errorf("got %v, %v, want 1, true", got, ok)
	}
	if got, ok := <-dc.Recv; !ok || got != 2 {
		t.Errorf("got %v, %v, want 2, true", got, ok)
	}

	select {
	case _, ok := <-dc.Recv:
		if ok {
			t.Errorf("got _, true, want _, false")
		}
	default:
	}

	close(dc.Send)

	if _, ok := <-dc.Recv; ok {
		t.Errorf("got _, true, want _, false")
	}
}

func TestHeap(t *testing.T) {
	dc := NewWithBuffer(NewHeap[int]())
	defer dc.Close()

	dc.Send <- 2
	time.Sleep(time.Second) // A little race-y: give the sending goroutine time to spin up and block on sending this value before the other values arrive.

	dc.Send <- 1
	dc.Send <- 0

	if got, ok := <-dc.Recv; !ok || got != 2 {
		t.Errorf("got %v, %v, want 2, true", got, ok)
	}
	if got, ok := <-dc.Recv; !ok || got != 0 {
		t.Errorf("got %v, %v, want 0, true", got, ok)
	}
	if got, ok := <-dc.Recv; !ok || got != 1 {
		t.Errorf("got %v, %v, want 1, true", got, ok)
	}

	select {
	case _, ok := <-dc.Recv:
		if ok {
			t.Errorf("got _, true, want _, false")
		}
	default:
	}

	close(dc.Send)

	if _, ok := <-dc.Recv; ok {
		t.Errorf("got _, true, want _, false")
	}
}
