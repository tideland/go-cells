// Tideland Go Cells - Behaviors - Pairer - Unit Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package pairer_test // import "tideland.dev/go/cells/behaviors/pairer"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/cells/behaviors/pairer"
	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSuccess verifies the successful finding of at least two matching.
func TestSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	pairerFunc := func(fstEvt, sndEvt *mesh.Event) bool {
		// Simply try to find the second event with the same
		// topic as the first.
		if fstEvt == nil {
			return true
		}
		return fstEvt.Topic() == sndEvt.Topic()
	}
	behavior := pairer.New(pairerFunc, time.Second)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		tbe.Push(evt)
		count := 0
		tbe.Do(func(i int, eevt *mesh.Event) error {
			switch evt.Topic() {
			case pairer.TopicPairerMatch:
				count++
				if count == 1 {
					// Skip waiting for second.
					return nil
				}
				var pair pairer.Pair
				if err := evt.Payload(&pair); err != nil {
					tbe.SetFail("event payload is no pair: %v", evt)
				}
				if pair.First.Topic() != pair.Second.Topic() {
					tbe.SetFail("pair event topic missmatch: %q <> %q", pair.First.Topic(), pair.Second.Topic())
				}
				tbe.SetSuccess()
			case pairer.TopicPairerTimeout:
				tbe.SetFail("no pairing received")
			}
			return nil
		})
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := generator.LimitedWord(4, 5)
			out.Emit(topic)
		}
	}, 5*time.Second)
	assert.NoError(err)
}

// TestFailOneHit verifies the failing of finding a pair after already
// one has been found.
func TestFailOneHit(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	pairerFunc := func(fstEvt, sndEvt *mesh.Event) bool {
		// Only find first one.
		if fstEvt == nil {
			return true
		}
		return false
	}
	behavior := pairer.New(pairerFunc, 10*time.Millisecond)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		switch evt.Topic() {
		case pairer.TopicPairerTimeout:
			tbe.SetSuccess()
		case pairer.TopicPairerMatch:
			tbe.SetFail("pairing found")
		}
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 1000; i++ {
			topic := generator.LimitedWord(4, 5)
			out.Emit(topic)
		}
	}, time.Second)
	assert.NoError(err)
}

// TestFailNoHit verifies the failing of finding any pair.
func TestFailNoHit(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	pairerFunc := func(fstEvt, sndEvt *mesh.Event) bool {
		// Never find any one.
		return false
	}
	behavior := pairer.New(pairerFunc, 10*time.Millisecond)
	// Testing.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		switch evt.Topic() {
		case pairer.TopicPairerMatch, pairer.TopicPairerTimeout:
			tbe.SetFail("any pairer output wrong here")
		}
		return nil
	}
	// Run test.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		for i := 0; i < 10000; i++ {
			topic := generator.LimitedWord(4, 5)
			out.Emit(topic)
		}
	}, time.Second)
	assert.ErrorContains(err, "timeout")
}

// EOF
