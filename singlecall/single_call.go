package singlecall

import "sync"

type key string

type callable func() (interface{}, error)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type SingleCall struct {
	wg    sync.WaitGroup
	l     sync.Mutex
	calls map[key]*call
}

func New() *SingleCall {
	return &SingleCall{
		l:     sync.Mutex{},
		calls: make(map[key]*call),
	}
}

func (sc *SingleCall) Run(k key, fn callable) (interface{}, error) {
	sc.l.Lock()
	if c, ok := sc.calls[k]; ok {
		sc.l.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	sc.calls[k] = c
	sc.l.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	sc.l.Lock()
	delete(sc.calls, k)
	sc.l.Unlock()

	return c.val, c.err
}
