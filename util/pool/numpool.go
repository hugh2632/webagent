package pool

import (
	"sync"
	"time"
)

type mypool struct {
	capacity int
	busy     int
	new      func() interface{}
	lock     sync.RWMutex
}

func NewNumPool(cap int, f func() interface{}) *mypool {
	return &mypool{
		capacity: cap,
		new:      f,
	}
}

func (m *mypool) Get() interface{} {
	for {
		m.lock.Lock()
		if m.busy < m.capacity {
			m.busy++
			m.lock.Unlock()
			return m.new()
		} else {
			m.lock.Unlock()
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (m *mypool) Free(v interface{}) {
	m.lock.Lock()
	v = nil
	m.busy--
	m.lock.Unlock()
}
