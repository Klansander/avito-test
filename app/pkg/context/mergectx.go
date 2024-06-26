package context

import (
	"context"
	"sync"
	"time"
)

type cx struct {
	sync.Mutex
	c0, c1 context.Context
	cq     chan struct{}
	err    error
	dlFunc func() (time.Time, bool)
}

func newCtx(c0, c1 context.Context) *cx {
	return &cx{c0: c0, c1: c1, cq: make(chan struct{})}
}

func Link(c0, c1 context.Context) context.Context {
	c := newCtx(c0, c1)
	c.dlFunc = c.first
	go c.link()
	return c
}

func (c *cx) Deadline() (deadline time.Time, ok bool) { return c.dlFunc() }

func (c *cx) first() (deadline time.Time, ok bool) {
	if d1, ok1 := c.c0.Deadline(); !ok1 {
		deadline, ok = c.c1.Deadline()
	} else if d2, ok2 := c.c1.Deadline(); !ok2 {
		deadline, ok = d1, true
	} else if d2.Before(d1) {
		deadline, ok = d2, true
	} else {
		deadline, ok = d1, true
	}

	return
}

func (c *cx) Done() <-chan struct{} { return c.cq }

func (c *cx) Err() error {
	c.Lock()
	defer c.Unlock()
	return c.err
}

func (c *cx) Value(key interface{}) (v interface{}) {
	if v = c.c0.Value(key); v == nil {
		v = c.c1.Value(key)
	}
	return
}

func (c *cx) link() {
	var dc context.Context
	select {
	case <-c.c0.Done():
		dc = c.c0
	case <-c.c1.Done():
		dc = c.c1
	case <-c.cq:
		return
	}

	c.Lock()
	if c.err == nil {
		c.err = dc.Err()
		close(c.cq)
	}
	c.Unlock()
}
