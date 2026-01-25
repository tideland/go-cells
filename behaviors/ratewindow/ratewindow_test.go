// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package ratewindow_test

import (
	"testing"
	"time"

	"tideland.dev/go/asserts/generators"
	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/behaviors/ratewindow"
	"tideland.dev/go/cells/mesh"
)

// TestSuccess verifies the successful finding and processing of matching events.
func TestSuccess(t *testing.T) {
	generator := generators.New(generators.FixedRand())
	matcher := func(evt *mesh.Event) (bool, error) {
		var payload int
		err := evt.Payload(&payload)
		verify.NoError(t, err)
		return payload > 5, nil
	}
	processor := func(reader mesh.EventSinkReader) (any, error) {
		var count int
		var sum int
		doer := func(i int, evt *mesh.Event) error {
			var payload int
			err := evt.Payload(&payload)
			verify.NoError(t, err)
			count += 1
			sum += payload
			return nil
		}
		err := reader.Do(doer)
		return sum / count, err
	}
	behavior := ratewindow.New(matcher, 50, time.Second, processor)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.WaitFor(func() bool { return tbe.Len() > 0 })
			tbe.Do(func(i int, evt *mesh.Event) error {
				var payload int
				err := evt.Payload(&payload)
				tbe.Assert(err == nil, "invalid payload: payload is no int")
				tbe.Assert(payload > 5, "invalid payload: payload is <= 5")
				return nil
			})
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := "int"
			payload := generator.Int(0, 10)
			out.Emit(topic, payload)
		}
	}, 10*time.Second)
	verify.NoError(t, err)
}
