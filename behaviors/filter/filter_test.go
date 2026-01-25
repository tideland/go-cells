// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package filter_test

import (
	"testing"
	"time"

	"tideland.dev/go/asserts/generators"
	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/behaviors/filter"
	"tideland.dev/go/cells/mesh"
)

// TestIncludingSuccess verifies the successful including filter of events.
func TestIncludingSuccess(t *testing.T) {
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
	verify.NoError(t, err)
}

// TestExcludingSuccess verifies the successful excluding filter of events.
func TestExcludingSuccess(t *testing.T) {
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
	verify.NoError(t, err)
}
