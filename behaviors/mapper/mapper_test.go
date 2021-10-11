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
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		switch evt.Topic() {
		case "!!!":
			tbe.SignalSuccess()
		case "map":
			var mapped []string
			if err := evt.Payload(&mapped); err != nil {
				tbe.SignalFail("error accessing payload: %v", err)
			}
			if strings.ToUpper(mapped[0]) != mapped[1] {
				tbe.SignalFail("payload not mapped: %v", mapped)
			}
		}
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 1000; i++ {
			orig := generator.Word()
			out.Emit("map", []string{orig, orig})
		}
		out.Emit("!!!")
	}, time.Second)
	assert.NoError(err)
}

// EOF
