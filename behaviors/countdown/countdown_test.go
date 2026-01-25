// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package countdown_test

import (
	"testing"
	"time"

	"tideland.dev/go/asserts/generators"
	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/behaviors/countdown"
	"tideland.dev/go/cells/mesh"
)

// TestSuccess verifies the countdown of events.
func TestSuccess(t *testing.T) {
	generator := generators.New(generators.FixedRand())
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	count := 10
	zeroer := func(r mesh.EventSinkReader) (*mesh.Event, error) {
		if r.Len() == count {
			return mesh.NewEvent("length-ok")
		}
		return nil, nil
	}
	behavior := countdown.New(count, zeroer)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() == 1 }, "not yet all events emitted")
			evt, ok := tbe.First()
			tbe.Assert(ok, "no first event")
			tbe.Assert(evt.Topic() == "length-ok", "topic is not 'length-ok': %v", evt.Topic())
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < count; i++ {
			topic := generator.OneStringOf(topics...)
			out.Emit(topic)
		}
	}, time.Second)
	verify.NoError(t, err)
}
