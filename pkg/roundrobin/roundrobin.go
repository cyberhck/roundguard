package roundrobin

import (
	"errors"
	"sync"
	"sync/atomic"
)

type RoundRobin[T any] struct {
	items []T
	index int32
	// we're using RW mutex so that we only block reads while the reset is going on, this helps because a normal lock would throttle requests (no concurrent lock would be acquired)
	mutex *sync.RWMutex
}

func (r *RoundRobin[T]) ResetWithNewItems(items []T) {
	r.mutex.Lock()
	// once we have this lock, no other goroutines will be able to read the value till this function has finished executing.
	defer r.mutex.Unlock()
	r.items = items
}

func (r *RoundRobin[T]) GetAllItems() []T {
	return r.items
}

func (r *RoundRobin[T]) GetNext() (*T, error) {
	if len(r.items) == 0 {
		return nil, errors.New("no remaining items")
	}
	newIndex := (r.index + 1) % int32(len(r.items))
	atomic.StoreInt32(&r.index, newIndex)
	r.mutex.RLock()
	defer r.mutex.RUnlock() // I'm using Read lock/unlock so that it's independent of reads/writes.

	return &r.items[r.index], nil
}

func New[T any](items []T) *RoundRobin[T] {
	return &RoundRobin[T]{
		items: items,
		mutex: &sync.RWMutex{},
		index: -1, // starting with -1 so that I can keep the GetNext method receiver simple
	}
}
