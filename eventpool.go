package xmuslogger

import "sync"

var eventPool = sync.Pool{
	New: func() interface{} {
		return &Event{buf: make([]byte, 0, 100)}
	},
}

func getEvent() *Event {
	e := eventPool.Get().(*Event)
	e.buf = e.buf[:0]
	return e
}

func putEvent(e *Event) {
	if cap(e.buf) <= 1<<16 { // Don't pool oversized buffers
		eventPool.Put(e)
	}
}
