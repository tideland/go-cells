// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package condition_test

import (
	"testing"
	"time"

	"tideland.dev/go/asserts/generators"
	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/behaviors/condition"
	"tideland.dev/go/cells/mesh"
)

// TestSuccess verifies the successful scanning for conditions.
func TestSuccess(t *testing.T) {
	generator := generators.New(generators.FixedRand())
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}
	tester := func(evt *mesh.Event) bool {
		return evt.Topic() == "now"
	}
	processor := func(cell mesh.Cell, evt *mesh.Event, out mesh.Emitter) error {
		topic := "found-" + evt.Topic()
		return out.Emit(topic)
	}
	behavior := condition.New(tester, processor)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() > 0 }, "collected events have to be at least 1: %d", tbe.Len())

			tbe.Do(func(i int, evt *mesh.Event) error {
				tbe.Assert(evt.Topic() == "found-now", "collected event topic has to be 'found-now': %v", evt)
				return nil
			})
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 50; i++ {
			topic := generator.OneStringOf(topics...)
			out.Emit(topic)
		}
	}, time.Second)
	verify.NoError(t, err)
}
