// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mapper_test

import (
	"strings"
	"testing"
	"time"

	"tideland.dev/go/asserts/generators"
	"tideland.dev/go/asserts/verify"
	"tideland.dev/go/cells/behaviors/mapper"
	"tideland.dev/go/cells/mesh"
)

// TestSuccess verifies mapping of events by upper-casing their payload.
func TestSuccess(t *testing.T) {
	generator := generators.New(generators.FixedRand())
	mapperFunc := func(evt *mesh.Event) (*mesh.Event, error) {
		if evt.Topic() != "map" {
			return evt, nil
		}
		var in []string
		if err := evt.Payload(&in); err != nil {
			return nil, err
		}
		out := []string{in[0], strings.ToUpper(in[1])}
		return mesh.NewEvent(evt.Topic(), out)
	}
	behavior := mapper.New(mapperFunc)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator) {
			tbe.WaitFor(func() bool { return tbe.Len() == 1000 })
			tbe.Do(func(i int, evt *mesh.Event) error {
				var mapped []string
				err := evt.Payload(&mapped)
				tbe.Assert(err == nil, "error accessing payload: %v", err)
				tbe.Assert(strings.ToUpper(mapped[0]) == mapped[1], "payload %d not mapped: %v", i, mapped)
				return nil
			})
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 1000; i++ {
			orig := generator.Word()
			out.Emit("map", []string{orig, orig})
		}
	}, time.Second)
	verify.NoError(t, err)
}
