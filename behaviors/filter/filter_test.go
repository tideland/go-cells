// Tideland Go Cells - Behaviors - Filter - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package filter_test // import "tideland.dev/go/cells/behaviors/filter"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/filter"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestIncludingSuccess verifies the successful including filter of events.
func TestIncludingSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	filterFunc := func(evt *mesh.Event) (bool, error) {
		// Filter wants to include short names.
		return len(evt.Topic()) < 6, nil
	}
	behavior := filter.NewIncluding(filterFunc)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		if evt.Topic() == "!!!" {
			tbe.SignalSuccess()
		}
		if len(evt.Topic()) > 5 {
			tbe.SignalFail("topic length of %q too long (> 5) for filter", evt.Topic())
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
		out.Emit("!!!")
	}, time.Second)
	assert.NoError(err)
}

// TestExcludingSuccess verifies the successful excluding filter of events.
func TestExcludingSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	filterFunc := func(evt *mesh.Event) (bool, error) {
		// Filter wants to exclude short names.
		return len(evt.Topic()) < 6, nil
	}
	behavior := filter.NewExcluding(filterFunc)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		if evt.Topic() == "!!!!!!" {
			tbe.SignalSuccess()
		}
		if len(evt.Topic()) < 6 {
			tbe.SignalFail("topic length of %q too short (< 6) for filter", evt.Topic())
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
		out.Emit("!!!!!!")
	}, time.Second)
	assert.NoError(err)
}

// EOF
