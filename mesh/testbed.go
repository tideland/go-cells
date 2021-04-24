// Tideland Go Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

//--------------------
// TESTBED EVALUATOR
//--------------------

// TestbedContext allows to store events received with the testbed
// evaluator. It's for tests where multiple events are needed
// for a positive evaluation.
//
// A success can be signaled with SetSuccess(), a failing with
// SetFail(string).
type TestbedContext struct {
	mu      sync.Mutex
	sink    *EventSink
	done    bool
	success bool
	reason  string
}

// newTestbedContext returns an initialized testbed context.
func newTestbedContext() *TestbedContext {
	return &TestbedContext{
		sink:    NewEventSink(0),
		done:    false,
		success: false,
		reason:  "",
	}
}

// EventSink returns the event sink of the context.
func (tbctx *TestbedContext) EventSink() *EventSink {
	return tbctx.sink
}

// SetSuccess signals a successful testing.
func (tbctx *TestbedContext) SetSuccess() {
	tbctx.done = true
	tbctx.success = true
}

// SetFail signals a failing testing together with a reason.
func (tbctx *TestbedContext) SetFail(reason string) {
	tbctx.done = true
	tbctx.success = false
	tbctx.reason = reason
}

// isDone returns true if the testing is done.
func (tbctx *TestbedContext) isDone() bool {
	return tbctx.done
}

// isDone returns true in case of a successful test, otherwise
// false and the reason.
func (tbctx *TestbedContext) isSuccesful() (bool, string) {
	return tbctx.success, tbctx.reason
}

// TestbedEvaluator defines a function signature used for evaluating
// the events emitted by the tested behavior. Success or failing can
// be sugnalled via the given testbed context.
type TestbedEvaluator func(tbctx *TestbedContext, evt *Event) error

//--------------------
// TESTBED RUNNER
//--------------------

// TestbedRunner contains the operations running in the background
// and emitting all events used by the tested behaviors as input.
type TestbedRunner func(out Emitter)

//--------------------
// TESTBED MESH
//--------------------

// testbedMesh implements the Mesh interface.
type testbedMesh struct{}

// Go implements Mesh and always returns an error.
func (tbm testbedMesh) Go(name string, b Behavior) error {
	return fmt.Errorf("cell name '%s' already used", name)
}

// Subscribe implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) Subscribe(emitterName, receptorName string) error {
	return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
}

// Unsubscribe implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) Unsubscribe(emitterName, receptorName string) error {
	return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
}

// Emit implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) Emit(name, topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return tbm.EmitEvent(name, evt)
}

// EmitEvent implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) EmitEvent(name string, evt *Event) error {
	return fmt.Errorf("cell '%s' does not exist", name)
}

// Emitter implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) Emitter(name string) (Emitter, error) {
	return nil, fmt.Errorf("cell '%s' does not exist", name)
}

//--------------------
// TESTBED CELL
//--------------------

// testbedCell runs the behavior and provides the needed interfaces.
type testbedCell struct {
	ctx      context.Context
	testbed  *Testbed
	behavior Behavior
	inc      chan *Event
}

// newTestbedCell initializes the testbed cell and spawns the goroutine.
func newTestbedCell(ctx context.Context, tb *Testbed, behavior Behavior) *testbedCell {
	tbc := &testbedCell{
		ctx:      ctx,
		testbed:  tb,
		behavior: behavior,
		inc:      make(chan *Event),
	}
	go tbc.backend()
	return tbc
}

// Context imepelements mesh.Cell.
func (tbc *testbedCell) Context() context.Context {
	return tbc.ctx
}

// Name imepelements mesh.Cell and returns a static name.
func (tbc *testbedCell) Name() string {
	return "testbed"
}

// Mesh imepelements mesh.Cell.
func (tbc *testbedCell) Mesh() Mesh {
	return testbedMesh{}
}

// Pull implements mesh.Receptor.
func (tbc *testbedCell) Pull() <-chan *Event {
	return tbc.inc
}

// Emit implements mesh.Emitter.
func (tbc *testbedCell) Emit(topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return tbc.EmitEvent(evt)
}

// EmitEvent implements mesh.Emitter and evaluates the event.
func (tbc *testbedCell) EmitEvent(evt *Event) error {
	evt.appendEmitter(tbc.Name())
	if err := tbc.testbed.eval(tbc.testbed.tbctx, evt); err != nil {
		tbc.testbed.errc <- err
		return err
	}
	if tbc.testbed.tbctx.isDone() {
		tbc.testbed.donec <- struct{}{}
	}
	return nil
}

// push writers an event into the input channel.
func (tbc *testbedCell) push(evt *Event) error {
	select {
	case <-tbc.ctx.Done():
		return errors.New("cell already terminated")
	case tbc.inc <- evt:
		return nil
	}
}

// backend runs the behavior to test.
func (tbc *testbedCell) backend() {
	if err := tbc.behavior.Go(tbc, tbc, tbc); err != nil {
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

// testbedEmitter allows the testbed runner to emit events to the testbed.
type testbedEmitter struct {
	testbed *Testbed
}

// newTesbedEmitter initializes the testbed emitter.
func newTestbedEmitter(tb *Testbed) *testbedEmitter {
	return &testbedEmitter{
		testbed: tb,
	}
}

// Emit creates an event and sends it to the behavior.
func (tbe *testbedEmitter) Emit(topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return tbe.EmitEvent(evt)
}

// Emit sends an event to the behavior.
func (tbe *testbedEmitter) EmitEvent(evt *Event) error {
	evt.initEmitters()
	return tbe.testbed.cell.push(evt)
}

//--------------------
// TESTBED
//--------------------

// Testbed provides a simple environment for the testing of individual behaviors.
// So retrieving the Mesh by the Cell is possible, but using its methods leads to
// errors. Integration tests have to be done differently.
//
// A tester function given when the testbed is started allows to evaluate the
// events emitted by the behavior. As long as the tests aren't done the function
// has to return false. Once returning true for the final tested event
// Testbed.Wait() gets a signal. Otherwise a timeout will be returned to show
// an internal error.
type Testbed struct {
	ctx      context.Context
	tbctx    *TestbedContext
	cancel   func()
	eval     TestbedEvaluator
	cell     *testbedCell
	donec    chan struct{}
	stoppedc chan struct{}
	errc     chan error
}

// NewTestbed starts a test cell with the given behavior. The tester function
// will be called for each event emitted by the behavior.
func NewTestbed(behavior Behavior, evaluator TestbedEvaluator) *Testbed {
	ctx, cancel := context.WithCancel(context.Background())
	tb := &Testbed{
		ctx:      ctx,
		tbctx:    newTestbedContext(),
		cancel:   cancel,
		eval:     evaluator,
		donec:    make(chan struct{}),
		stoppedc: make(chan struct{}),
		errc:     make(chan error),
	}
	tb.cell = newTestbedCell(ctx, tb, behavior)
	return tb
}

// Go runs the testbed rzbber and waits until test ends or a timeout.
func (tb *Testbed) Go(runner TestbedRunner, timeout time.Duration) error {
	go tb.run(runner)
	return tb.wait(timeout)
}

// run runs the testbed runner.
func (tb *Testbed) run(runner TestbedRunner) {
	runner(newTestbedEmitter(tb))
	tb.stoppedc <- struct{}{}
}

// wait waits until a test end or error has been signalled or a
// timeout happened.
func (tb *Testbed) wait(timeout time.Duration) error {
	defer tb.cancel()
	running := true
	waiting := true
	now := time.Now()
	for running || waiting {
		select {
		case <-tb.stoppedc:
			running = false
		case <-tb.donec:
			waiting = false
		case err := <-tb.errc:
			return err
		case to := <-time.After(timeout):
			waited := to.Sub(now)
			return errors.New("test failed: timeout after " + waited.String())
		}
	}
	// The async testbed runner and the tests are done.
	success, reason := tb.tbctx.isSuccesful()
	if !success {
		return errors.New("test failed: " + reason)
	}
	return nil
}

// EOF
