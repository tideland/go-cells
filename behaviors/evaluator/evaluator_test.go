// Tideland Go Cells - Behaviors - Evaluator - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package evaluator_test // import "tideland.dev/go/cells/behaviors/evaluator"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/evaluator"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful evaluation of events.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	evaluateFunc := func(evt *mesh.Event) (float64, error) {
		l := len(evt.Topic())
		return float64(l), nil
	}
	behavior := evaluator.New(evaluateFunc)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() == 2 }, "not yet all events emitted")
			evt, ok := tbe.Peek(0)
			tbe.Assert(ok, "first event missing")
			tbe.Assert(evt.Topic() == evaluator.TopicEvaluationDone, "topic wrong: %v", evt.Topic())
			var evaluation evaluator.Evaluation
			err := evt.Payload(&evaluation)
			tbe.Assert(err == nil, "payload error not nil: %v", err)
			tbe.Assert(evaluation.Count == 10000, "evaluation count not 10000: %v", evaluation.Count)
			tbe.Assert(evaluation.MinRating == 3.0, "evaluation min rating not 3.0: %v", evaluation.MinRating)
			tbe.Assert(evaluation.MaxRating == 8.0, "evaluation max rating not 8.0: %v", evaluation.MaxRating)
			evt, ok = tbe.Peek(1)
			tbe.Assert(ok, "second event missing")
			tbe.Assert(evt.Topic() == evaluator.TopicResetDone, "topic not equal 'reset-done': %v", evt.Topic())
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := generator.LimitedWord(3, 8)
			out.Emit(topic)
		}
		// Retrieve and reset.
		out.Emit(evaluator.TopicEvaluate)
		out.Emit(evaluator.TopicReset)
	}, time.Second)
	assert.NoError(err)
}

// TestFail verifies the wanted failing of the evaluation.
func TestFail(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evaluateFunc := func(evt *mesh.Event) (float64, error) {
		if evt.Topic() == "ouch" {
			return 0.0, errors.New("ouch")
		}
		return 1.0, nil
	}
	behavior := evaluator.New(evaluateFunc)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() == 1 }, "no event published")
			evt, ok := tbe.First()
			tbe.Assert(ok, "cannot retrieve first event")
			tbe.Assert(evt.Topic() == mesh.TopicTestbedError, "invalid topic: %v", evt.Topic())
			var cellErr mesh.PayloadCellError
			err := evt.Payload(&cellErr)
			tbe.Assert(err == nil, "retrieving payload returned an error: %v", err)
			tbe.Assert(cellErr.Error == "ouch", "invalid returned cell error: %v", cellErr.Error)
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("ouch")
	}, time.Second)
	assert.NoError(err)
}

// EOF
