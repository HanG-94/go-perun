// Copyright (c) 2019 The Perun Authors. All rights reserved.
// This file is part of go-perun. Use of this source code is governed by a
// MIT-style license that can be found in the LICENSE file.

package sync

import (
	"sync"

	"github.com/pkg/errors"

	"perun.network/go-perun/pkg/sync/atomic"
)

// Closer is a utility type for implementing an "onclose" event.
// It supports registering handlers, waiting for the event, and status checking.
// A default-initialised Closer is a valid value.
type Closer struct {
	once     sync.Once     // Initializes the closed channel.
	isClosed atomic.Bool   // Whether the Closer is currently closed.
	closed   chan struct{} // Closed when the Closer is closed.

	onClosedMtx sync.Mutex // Protects callbacks.
	onClosed    []func()   // Executed when Close() is called.
}

func (c *Closer) initOnce() {
	c.once.Do(func() { c.closed = make(chan struct{}) })
}

// Closed returns a channel to be used in a select statement.
func (c *Closer) Closed() <-chan struct{} {
	c.initOnce()
	return c.closed
}

// Close starts all registered callbacks in goroutines.
// If Close was already called before, returns an AlreadyClosedError, otherwise,
// returns nil.
func (c *Closer) Close() error {
	c.initOnce()

	if !c.isClosed.TrySet() {
		return newAlreadyClosedError()
	}

	close(c.closed)

	c.onClosedMtx.Lock()
	defer c.onClosedMtx.Unlock()
	for _, fn := range c.onClosed {
		go fn()
	}

	return nil
}

// IsClosed returns whether the Closer is currently closed.
func (c *Closer) IsClosed() bool {
	return c.isClosed.IsSet()
}

// OnClose registers the passed callback to be called when the Closer is closed.
// If the Closer is already closed, immediately executes the callback in a
// goroutine.
func (c *Closer) OnClose(handler func()) {
	c.onClosedMtx.Lock()
	defer c.onClosedMtx.Unlock()
	// Check again, because Close might have been called before the lock was
	// acquired.
	if c.IsClosed() {
		go handler()
	} else {
		c.onClosed = append(c.onClosed, handler)
	}
}

var _ error = alreadyClosedError{}

type alreadyClosedError struct{}

const alreadyClosedMsg = "Closer already closed"

func (alreadyClosedError) Error() string {
	return alreadyClosedMsg
}

func newAlreadyClosedError() error {
	return errors.WithStack(alreadyClosedError{})
}

// IsAlreadyClosedError checks whether an error is an AlreadyClosedError.
func IsAlreadyClosedError(err error) bool {
	_, ok := errors.Cause(err).(alreadyClosedError)
	return ok
}
