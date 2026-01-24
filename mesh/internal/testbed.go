// Tideland Go Cells - Mesh - Internal
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package internal // import "tideland.dev/go/cells/mesh/internal"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
)

//--------------------
// TESTBED MESH
//--------------------

// TestbedMesh implements the Mesh interface.
type testbedMeshImpl struct{}

// Go implements Mesh and always returns an error.
func (tbm testbedMeshImpl) Go(name string, b Behavior) error {
	return fmt.Errorf("cell name '%s' already used", name)
}

// Subscribe implements Mesh and always returns an error.
func (tbm testbedMeshImpl) Subscribe(emitterName, receptorName string) error {
	return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
}

// Unsubscribe implements Mesh and always returns an error.
func (tbm testbedMeshImpl) Unsubscribe(emitterName, receptorName string) error {
	return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
}

// Emit implements Mesh and always returns an error.
func (tbm testbedMeshImpl) Emit(name, topic string, payloads ...any) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return tbm.EmitEvent(name, evt)
}

// EmitEvent implements Mesh and always returns an error.
func (tbm testbedMeshImpl) EmitEvent(name string, evt *Event) error {
	return fmt.Errorf("cell '%s' does not exist", name)
}

// Emitter implements Mesh and always returns an error.
func (tbm testbedMeshImpl) Emitter(name string) (Emitter, error) {
	return nil, fmt.Errorf("cell '%s' does not exist", name)
}

//--------------------
// TESTBED CELL
//--------------------

// TestbedCell runs the behavior and provides the needed interfaces.
type testbedCellImpl struct {
	ctx      context.Context
	behavior Behavior
	inc      chan *Event
	pusher   func(evt *Event)
}

// NewTestbedCell initializes the testbed cell and spawns the goroutine.
func NewTestbedCell(ctx context.Context, behavior Behavior, pusher func(evt *Event)) interface{Push(evt *Event) error} {
	tbc := &testbedCellImpl{
		ctx:      ctx,
		behavior: behavior,
		inc:      make(chan *Event),
		pusher:   pusher,
	}
	go tbc.backend()
	return tbc
}

// Context implements mesh.Cell.
func (tbc *testbedCellImpl) Context() context.Context {
	return tbc.ctx
}

// Name implements mesh.Cell and returns a static name.
func (tbc *testbedCellImpl) Name() string {
	return "testbed"
}

// Mesh implements Cell.
func (tbc *testbedCellImpl) Mesh() Mesh {
	return testbedMeshImpl{}
}

// Pull implements mesh.Receptor.
func (tbc *testbedCellImpl) Pull() <-chan *Event {
	return tbc.inc
}

// Emit implements Emitter.
func (tbc *testbedCellImpl) Emit(topic string, payloads ...any) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return tbc.EmitEvent(evt)
}

// EmitEvent implements Emitter and evaluates the event.
func (tbc *testbedCellImpl) EmitEvent(evt *Event) error {
	evt.appendEmitter(tbc.Name())
	tbc.pusher(evt)
	return nil
}

// Push writes an event into the input channel.
func (tbc *testbedCellImpl) Push(evt *Event) error {
	select {
	case <-tbc.ctx.Done():
		// Ignore as test result.
		return nil
	case tbc.inc <- evt:
		return nil
	}
}

// backend runs the behavior to test.
func (tbc *testbedCellImpl) backend() {
	// Execute the behavior.
	err := tbc.behavior.Go(tbc, tbc, tbc)
	if err != nil {
		// Notify subscribers about error.
		tbc.Emit(TopicTestbedError, PayloadCellError{
			CellName: tbc.Name(),
			Error:    err.Error(),
		})
	} else {
		// Notify subscribers about termination.
		tbc.Emit(TopicTestbedTerminated, PayloadTermination{
			CellName: tbc.Name(),
		})
	}
}

//--------------------
// TESTBED EMITTER
//--------------------

// TestbedEmitter allows the testbed runner to emit events to the testbed.
type testbedEmitterImpl struct {
	cell *testbedCellImpl
}

// NewTestbedEmitter initializes the testbed emitter.
func NewTestbedEmitter(cell interface{}) *testbedEmitterImpl {
	return &testbedEmitterImpl{
		cell: cell.(*testbedCellImpl),
	}
}

// Emit creates an event and sends it to the behavior.
func (tbe *testbedEmitterImpl) Emit(topic string, payloads ...any) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return tbe.EmitEvent(evt)
}

// EmitEvent sends an event to the behavior.
func (tbe *testbedEmitterImpl) EmitEvent(evt *Event) error {
	evt.initEmitters()
	return tbe.cell.Push(evt)
}

// EOF
