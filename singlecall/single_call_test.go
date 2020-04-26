package singlecall

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingleCall_Run_ReturnsValue(t *testing.T) {
	var c int32 = 0
	var val interface{}
	var err error
	sc := New()
	wg := new(sync.WaitGroup)
	fn := func() (interface{}, error) {
		time.Sleep(1 * time.Second) // let next calls to be enqueued
		atomic.AddInt32(&c, 1)
		return c, nil
	}
	calls := 10
	wg.Add(calls)
	for calls > 0 {
		go func() {
			val, err = sc.Run("counter", fn)
			wg.Done()
		}()
		calls--
	}
	wg.Wait()
	if err != nil {
		t.Errorf("unexpected error. want nil, got error %s", err.Error())
		return
	}
	if val != int32(1) {
		t.Errorf("unexpected value. want %d, got %d", 1, val)
	}
}

func TestSingleCall_Run_ReturnsError(t *testing.T) {
	var val interface{}
	var err error
	sc := New()
	wg := new(sync.WaitGroup)
	fn := func() (interface{}, error) {
		time.Sleep(1 * time.Second) // let next calls to be enqueued
		return nil, errors.New("sample error for testing")
	}
	calls := 10
	wg.Add(calls)
	for calls > 0 {
		go func() {
			val, err = sc.Run("counter", fn)
			wg.Done()
		}()
		calls--
	}
	wg.Wait()
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if val != nil {
		t.Errorf("expected value to be nil. got %v", val)
	}
}
