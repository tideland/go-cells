// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package internal

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// cellSet manages a set of cells.
type cellSet struct {
	mu    sync.RWMutex
	cells map[*CellImpl]struct{}
}

// newCellSet creates an empty cell set.
func newCellSet() *cellSet {
	return &cellSet{
		cells: make(map[*CellImpl]struct{}),
	}
}

// add adds another cell to the set. Already added
// ones are ignored.
func (cs *cellSet) add(c *CellImpl) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.cells[c] = struct{}{}
}

// remove deletes a cell from the set.
func (cs *cellSet) remove(c *CellImpl) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.cells, c)
}

// do perform f for each cell of the set.
func (cs *cellSet) do(f func(c *CellImpl) error) error {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	for c := range cs.cells {
		if err := f(c); err != nil {
			return err
		}
	}
	return nil
}

// Cell runs a behevior networked with other cells.
type CellImpl struct {
	mu       sync.RWMutex
	active   atomic.Value
	ctx      context.Context
	name     string
	meshRef  Mesh
	behavior Behavior
	in       *streamImpl
	input    *cellSet
	output   *cellSet
	drop     func()
}

// NewCell starts a new cell working in the background.
func NewCell(ctx context.Context, name string, m Mesh, b Behavior, drop func()) *CellImpl {
	c := &CellImpl{
		ctx:      ctx,
		name:     name,
		meshRef:  m,
		behavior: b,
		in:       NewStream(),
		input:    newCellSet(),
		output:   newCellSet(),
		drop:     drop,
	}
	c.active.Store(true)
	go c.backend()
	return c
}

// Context implements mesh.Cell.
func (c *CellImpl) Context() context.Context {
	return c.ctx
}

// Name implements mesh.Cell.
func (c *CellImpl) Name() string {
	return c.name
}

// Mesh implements mesh.Cell.
func (c *CellImpl) Mesh() Mesh {
	return nil
}

// SubscribeTo adds this cell to the out-streams of the
// given in-cell.
func (c *CellImpl) SubscribeTo(ic *CellImpl) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.input.add(ic)
	ic.output.add(c)
}

// UnsubscribeFrom removes this cell from the out-streams of the
// given in-cell.
func (c *CellImpl) UnsubscribeFrom(ic *CellImpl) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.input.remove(ic)
	ic.output.remove(c)
}

// Receive creates an passes an event to handle to the cell.
func (c *CellImpl) Receive(topic string, payload ...any) error {
	evt, err := NewEvent(topic, payload...)
	if err != nil {
		return err
	}
	return c.ReceiveEvent(evt)
}

// ReceiveEvent passes an event to handle to the cell.
func (c *CellImpl) ReceiveEvent(evt *Event) error {
	if !c.active.Load().(bool) {
		return errors.New("cell deactivated")
	}
	return c.in.EmitEvent(evt)
}

// shutdown deactivates the in-stream, unsubscribes from all cells
// and tells the mesh that it's not available anymore.
func (c *CellImpl) shutdown() {
	c.active.Store(false)
	c.drop()
	c.input.do(func(ic *CellImpl) error {
		ic.output.remove(c)
		return nil
	})
}

// Pull implements mesh.Receptor.
func (c *CellImpl) Pull() <-chan *Event {
	return c.in.Pull()
}

// Emit implements mesh.Emitter.
func (c *CellImpl) Emit(topic string, payloads ...any) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return c.EmitEvent(evt)
}

// EmitEvent implements mesh.Emitter.
func (c *CellImpl) EmitEvent(evt *Event) error {
	evt.appendEmitter(c.name)
	return c.output.do(func(oc *CellImpl) error {
		if err := oc.ReceiveEvent(evt); err != nil {
			return err
		}
		return nil
	})
}

// backend runs as goroutine and cares for the behavior.
func (c *CellImpl) backend() {
	defer c.shutdown()
	if err := c.behavior.Go(c, c, c); err != nil {
		// Notify subscribers about error.
		c.Emit(TopicError, PayloadCellError{
			CellName: c.name,
			Error:    err.Error(),
		})
	} else {
		// Notify subscribers about termination.
		c.Emit(TopicTerminated, PayloadTermination{
			CellName: c.name,
		})
	}
}
