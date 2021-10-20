// Tideland Go Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh_test // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestTestbedSuccess verifies the successful working of the testbed
// for behavior tests.
func TestTestbedSuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	forwarder := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				out.EmitEvent(evt)
			}
		}
	}
	behavior := mesh.BehaviorFunc(forwarder)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) {
			tbe.Push(evt)
		},
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() == 3 }, "collected events not 3: %d", tbe.Len())
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
	}, time.Second)
	assert.NoError(err)
}

// TestTestbedFail verifies the failing working of the testbed
// for behavior tests.
func TestTestbedFail(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	forwarder := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				out.EmitEvent(evt)
			}
		}
	}
	behavior := mesh.BehaviorFunc(forwarder)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) {},
		func(tbe *mesh.TestbedEvaluator) {
			tbe.Assert(false, "must fail")
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
	}, time.Second)
	assert.ErrorContains(err, "test failed: must fail")
}

// TestTestbedMesh verifies the Mesh stubbing of the testbed.
func TestTestbedMesh(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	mesher := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				switch evt.Topic() {
				case "go":
					out.EmitEvent(evt)
					err := cell.Mesh().Go("anything", nil)
					assert.ErrorContains(err, "cell name 'anything' already used")
				case "subscribe":
					out.EmitEvent(evt)
					err := cell.Mesh().Subscribe("anything", "anything-else")
					assert.ErrorContains(err, "emitter cell 'anything' does not exist")
				case "unsubscribe":
					out.EmitEvent(evt)
					err := cell.Mesh().Unsubscribe("anything", "anything-else")
					assert.ErrorContains(err, "emitter cell 'anything' does not exist")
				case "emit":
					out.EmitEvent(evt)
					err := cell.Mesh().EmitEvent("anything", evt)
					assert.ErrorContains(err, "cell 'anything' does not exist")
				case "emitter":
					out.EmitEvent(evt)
					emtr, err := cell.Mesh().Emitter("anything")
					assert.ErrorContains(err, "cell 'anything' does not exist")
					assert.Nil(emtr)
				}
			}
		}
	}
	behavior := mesh.BehaviorFunc(mesher)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) {
			tbe.Push(evt)
		},
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() == 5 }, "not all events processed: %v", tbe)
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("go")
		out.Emit("subscribe")
		out.Emit("unsubscribe")
		out.Emit("emit")
		out.Emit("emitter")
	}, time.Second)
	assert.NoError(err)
}

// TestTestbedError verifies the error handling of the Testbed.
func TestTestbedError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	failer := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				switch evt.Topic() {
				case "go":
					out.Emit("ok")
				default:
					// Fail when topic is unknown.
					out.Emit("fail")
				}
			}
		}
	}
	behavior := mesh.BehaviorFunc(failer)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
		func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) {
			if evt.Topic() == "fail" {
				tbe.SignalError(errors.New("ouch"))
			}
		},
		nil,
	)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("go")
		out.Emit("go")
		out.Emit("go")
		out.Emit("dunno!")
		out.Emit("go")
		out.Emit("go")
	}, time.Second)
	assert.ErrorContains(err, "test error: ouch")
}

// EOF
