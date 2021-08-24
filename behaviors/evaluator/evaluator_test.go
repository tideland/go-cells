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
	// Test events in testbed.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		switch evt.Topic() {
		case evaluator.TopicEvaluationDone:
			var evaluation evaluator.Evaluation
			if err := evt.Payload(&evaluation); err != nil {
				tbe.SetFail("can not retrieve evaluation payload: %v", err)
			}
			if evaluation.Count != 10000 {
				tbe.SetFail("evaluation count is wrong: %d", evaluation.Count)
			}
			if evaluation.MinRating != 3.0 {
				tbe.SetFail("evaluation min rating is wrong: %f", evaluation.MinRating)
			}
			if evaluation.MaxRating != 8.0 {
				tbe.SetFail("evaluation max rating is wrong: %f", evaluation.MaxRating)
			}
		case evaluator.TopicResetDone:
			tbe.SetSuccess()
		}
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
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
	// Test events in testbed.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		if evt.Topic() == mesh.TopicTestbedError {
			var cellError mesh.PayloadCellError
			if err := evt.Payload(&cellError); err != nil {
				tbe.SetFail("invalid payload")
			}
			if cellError.Error != "ouch" {
				tbe.SetFail("invalid error: %s", cellError.Error)
			}
			tbe.SetSuccess()
			return nil
		}
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("ouch")
	}, time.Second)
	assert.NoError(err)
}

// EOF
