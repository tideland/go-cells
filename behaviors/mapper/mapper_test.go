// Tideland Go Cells - Behaviors - Mapper - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mapper_test // import "tideland.dev/go/cells/behaviors/mapper"

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/cells/behaviors/mapper"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies mapping of events by upper-casing their payload.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
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
	assert.NoError(err)
}

// EOF
