// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tideland.dev/go/cells/mesh/internal"
)

// TestbedEvaluator allows to store events received during evaluation.
// Here it supports the interface EventSink.
//
// A success can be signaled with SignalSuccess(), a failing with
// SignalFail(reason string, vs ...any).
type TestbedEvaluator struct {
	EventSink

	tb *Testbed
}

// newTestbedEvaluator returns an initialized testbed context.
func newTestbedEvaluator(tb *Testbed) *TestbedEvaluator {
	return &TestbedEvaluator{
		EventSink: internal.NewEventSink(0),
		tb:        tb,
	}
}

// WaitFor waits until the given function returns true. It runs the function up to 5 times
// with growing pauses inbetween. If the final call returns false the test fails.
func (tbe *TestbedEvaluator) WaitFor(assertion func() bool) {
	tbe.AssertRetry(assertion, "waitong failed")
}

// AssertRetry waits until the given function returns true. It runs the function up to 5 times
// with growing pauses inbetween. If the final call returns false the test fails.
func (tbe *TestbedEvaluator) AssertRetry(assertion func() bool, reason string, vs ...any) {
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

// Assert tests if an assertion is true, otherwise it segnals a
// failing test.
func (tbe *TestbedEvaluator) Assert(assertion bool, reason string, vs ...any) {
	if assertion {
		return
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

// TestbedRunner contains the operations running in the background
// and emitting all events used by the tested behaviors as input.
type TestbedRunner func(out Emitter)

// TestbedTest defines a function signature used for checking the final
// success when the testbed runner has sent all test input events.
type TestbedTester func(tbe *TestbedEvaluator)

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
	test       TestbedTester
	cell       testbedCellRef
	succeededc chan struct{}
	failedc    chan string
	errc       chan error
}

// testbedCellRef is an internal reference type for testbed cell
type testbedCellRef interface {
	Push(evt *Event) error
}

// NewTestbed starts a test cell with the given behavior. The tester function
// will be called for each event emitted by the behavior.
func NewTestbed(behavior Behavior, tester TestbedTester) *Testbed {
	ctx, cancel := context.WithCancel(context.Background())
	tb := &Testbed{
		ctx:        ctx,
		cancel:     cancel,
		test:       tester,
		succeededc: make(chan struct{}, 1),
		failedc:    make(chan string, 1),
		errc:       make(chan error, 1),
	}
	tb.evaluator = newTestbedEvaluator(tb)
	tb.cell = internal.NewTestbedCell(ctx, behavior, func(evt *Event) {
		tb.evaluator.Push(evt)
	})
	return tb
}

// Go runs the testbed runner to emit events to the cell to test. Afterwards it
// waits timeout duration for testbed out tester and testbed end tester doing
// their tests. In case of no fail signal or error the tests succeeds.
func (tb *Testbed) Go(run TestbedRunner, timeout time.Duration) error {
	go func() {
		tbe := internal.NewTestbedEmitter(tb.cell)

		run(tbe)
		tb.test(tb.evaluator)

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
