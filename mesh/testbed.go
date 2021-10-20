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
	"time"
)

//--------------------
// TESTBED TESTER
//--------------------

// TestbedEvaluator allows to store events received during evaluation.
// Here it supports the interface EventSink.
//
// A success can be signaled with SignalSuccess(), a failing with
// SignalFail(reason string, vs ...interface{}).
type TestbedEvaluator struct {
	EventSink

	tb *Testbed
}

// newTestbedEvaluator returns an initialized testbed context.
func newTestbedEvaluator(tb *Testbed) *TestbedEvaluator {
	return &TestbedEvaluator{
		EventSink: NewEventSink(0),
		tb:        tb,
	}
}

// Assert tests if an assertion is true, otherwise it segnals a
// failing test.
func (tbe *TestbedEvaluator) Assert(assertion bool, reason string, vs ...interface{}) {
	if assertion {
		return
	}
	tbe.tb.failedc <- fmt.Sprintf(reason, vs...)
}

// AssertRetry tests 5 times with pauses if an assertion is true, otherwise it segnals a
// failing test.
func (tbe *TestbedEvaluator) AssertRetry(assertion func() bool, reason string, vs ...interface{}) {
	duration := 2 * time.Millisecond
	for i := 0; i < 5; i++ {
		if assertion() {
			return
		}
		time.Sleep(duration)

		duration = duration * 2
	}
	tbe.tb.failedc <- fmt.Sprintf(reason, vs...)
}

// SignalError signals an error during testing.
func (tbe *TestbedEvaluator) SignalError(err error) {
	tbe.tb.errc <- err
}

// String returns a testbed evaluator representation containing all
// the topics of the sink.
func (tbe *TestbedEvaluator) String() string {
	return "TestbedEvaluator{" + tbe.EventSink.String() + "}"
}

//--------------------
// TESTBED FUNCTIONS
//--------------------

// TestbedRunner contains the operations running in the background
// and emitting all events used by the tested behaviors as input.
type TestbedRunner func(out Emitter)

// TestbedOutTester defines a function signature used for evaluating the
// events emitted by the tested behavior. Success or failing can be
// sugnalled via the given testbed evaluator.
type TestbedOutTester func(tbe *TestbedEvaluator, evt *Event)

// TestbedEndTest defines a function signature used for checking the final
// success when the testbed runner has sent all test input events.
type TestbedEndTester func(tbe *TestbedEvaluator)

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
	tb       *Testbed
	behavior Behavior
	inc      chan *Event
}

// newTestbedCell initializes the testbed cell and spawns the goroutine.
func newTestbedCell(ctx context.Context, tb *Testbed, behavior Behavior) *testbedCell {
	tbc := &testbedCell{
		ctx:      ctx,
		tb:       tb,
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
	tbc.tb.testOut(tbc.tb.evaluator, evt)
	return nil
}

// push writes an event into the input channel.
func (tbc *testbedCell) push(evt *Event) error {
	select {
	case <-tbc.ctx.Done():
		// Ignore as test result has been defined otherwhere.
		return nil
	case tbc.inc <- evt:
		return nil
	}
}

// backend runs the behavior to test.
func (tbc *testbedCell) backend() {
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

// testbedEmitter allows the testbed runner to emit events to the testbed.
type testbedEmitter struct {
	tb *Testbed
}

// newTesbedEmitter initializes the testbed emitter.
func newTestbedEmitter(tb *Testbed) *testbedEmitter {
	return &testbedEmitter{
		tb: tb,
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
	return tbe.tb.cell.push(evt)
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
	ctx        context.Context
	cancel     func()
	evaluator  *TestbedEvaluator
	testOut    TestbedOutTester
	testEnd    TestbedEndTester
	cell       *testbedCell
	succeededc chan struct{}
	failedc    chan string
	errc       chan error
}

// NewTestbed starts a test cell with the given behavior. The tester function
// will be called for each event emitted by the behavior.
func NewTestbed(
	behavior Behavior,
	testOut TestbedOutTester,
	testEnd TestbedEndTester,
) *Testbed {
	ctx, cancel := context.WithCancel(context.Background())
	tb := &Testbed{
		ctx:        ctx,
		cancel:     cancel,
		testOut:    testOut,
		testEnd:    testEnd,
		succeededc: make(chan struct{}, 1),
		failedc:    make(chan string, 1),
		errc:       make(chan error, 1),
	}
	tb.evaluator = newTestbedEvaluator(tb)
	tb.cell = newTestbedCell(ctx, tb, behavior)
	return tb
}

// Go runs the testbed runner to emit events to the cell to test. Afterwards it
// waits timeout duration for testbed out tester and testbed end tester doing
// their tests. In case of no fail signal or error the tests succeeds.
func (tb *Testbed) Go(run TestbedRunner, timeout time.Duration) error {
	go func() {
		tbe := newTestbedEmitter(tb)

		run(tbe)

		if tb.testEnd != nil {
			tb.testEnd(tb.evaluator)
		}

		tb.succeededc <- struct{}{}
	}()

	return tb.wait(timeout)
}

// wait waits until a test end or error has been signalled or a
// timeout happened.
func (tb *Testbed) wait(timeout time.Duration) error {
	defer tb.cancel()
	now := time.Now()
	for {
		select {
		case <-tb.succeededc:
			return nil
		case reason := <-tb.failedc:
			return errors.New("test failed: " + reason)
		case err := <-tb.errc:
			return fmt.Errorf("test error: %v", err)
		case to := <-time.After(timeout):
			waited := to.Sub(now)
			return errors.New("test failed: timeout after " + waited.String())
		}
	}
}

// EOF
