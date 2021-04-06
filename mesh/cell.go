// Tideland Go Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORT
//--------------------

import (
	"context"
	"errors"
	"sync"
)

//--------------------
// CELL SET
//--------------------

// cellSet manages a set of cells.
type cellSet struct {
	mu    sync.RWMutex
	cells map[*cell]struct{}
}

// newCellSet creates an empty cell set.
func newCellSet() *cellSet {
	return &cellSet{
		cells: make(map[*cell]struct{}),
	}
}

// add adds another cell to the set. Already added
// ones are ignored.
func (cs *cellSet) add(c *cell) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.cells[c] = struct{}{}
}

// remove deletes a cell from the set.
func (cs *cellSet) remove(c *cell) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.cells, c)
}

// do perform f for each cell of the set.
func (cs *cellSet) do(f func(c *cell) error) error {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	for c := range cs.cells {
		if err := f(c); err != nil {
			return err
		}
	}
	return nil
}

//--------------------
// CELL
//--------------------

// cell runs a behevior networked with other cells.
type cell struct {
	mu       sync.RWMutex
	ctx      context.Context
	name     string
	mesh     Mesh
	behavior Behavior
	in       *stream
	input    *cellSet
	output   *cellSet
	drop     func()
}

// newCell starts a new cell working in the background.
func newCell(ctx context.Context, name string, m Mesh, b Behavior, drop func()) *cell {
	c := &cell{
		ctx:      ctx,
		name:     name,
		mesh:     m,
		behavior: b,
		in:       newStream(),
		input:    newCellSet(),
		output:   newCellSet(),
		drop:     drop,
	}
	go c.backend()
	return c
}

// Context implements Cell.
func (c *cell) Context() context.Context {
	return c.ctx
}

// Name implements Cell.
func (c *cell) Name() string {
	return c.name
}

// Mesh implements Cell.
func (c *cell) Mesh() Mesh {
	return nil
}

// subscribeTo adds this cell to the out-streams of the
// given in-cell.
func (c *cell) subscribeTo(ic *cell) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.input.add(ic)
	ic.output.add(c)
}

// unsubscribeFrom removes this cell from the out-streams of the
// given in-cell.
func (c *cell) unsubscribeFrom(ic *cell) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.input.remove(ic)
	ic.output.remove(c)
}

// receive creates an passes an event to handle to the cell.
func (c *cell) receive(topic string, payload ...interface{}) error {
	evt, err := NewEvent(topic, payload...)
	if err != nil {
		return err
	}
	return c.receiveEvent(evt)
}

// receiveEvent passes an event to handle to the cell.
func (c *cell) receiveEvent(evt Event) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.in == nil {
		return errors.New("cell deactivated")
	}
	return c.in.EmitEvent(evt)
}

// shutdown deactivates the in-stream, unsubscribes from all cells
// and tells the mesh that it's not available anymore.
func (c *cell) shutdown() {
	c.drop()
	c.in = nil
	c.input.do(func(ic *cell) error {
		ic.output.remove(c)
		return nil
	})
}

// Pull implements Receptor.
func (c *cell) Pull() <-chan Event {
	return c.in.Pull()
}

// Emit implements Emitter.
func (c *cell) Emit(topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return c.EmitEvent(evt)
}

// EmitEvent implements Emitter.
func (c *cell) EmitEvent(evt Event) error {
	return c.output.do(func(oc *cell) error {
		if err := oc.receiveEvent(evt); err != nil {
			return err
		}
		return nil
	})
}

// backend runs as goroutine and cares for the behavior.
func (c *cell) backend() {
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

// EOF
