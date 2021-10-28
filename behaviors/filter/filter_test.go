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
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.Do(func(i int, evt *mesh.Event) error {
				tbe.Assert(len(evt.Topic()) < 6, "topic length of event %d too long: %v", i, evt)
				return nil
			})
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := generator.LimitedWord(3, 8)
			out.Emit(topic)
		}
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
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.Do(func(i int, evt *mesh.Event) error {
				tbe.Assert(len(evt.Topic()) >= 6, "topic length of event %d too short: %v", i, evt)
				return nil
			})
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := generator.LimitedWord(3, 8)
			out.Emit(topic)
		}
	}, time.Second)
	assert.NoError(err)
}

// EOF
