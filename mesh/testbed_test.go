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
	// Test evaluation.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		if tbe.Done(evt) {
			if tbe.Len() == 3 {
				tbe.SignalSuccess()
			}
		}
		tbe.Push(evt)
		return nil
	}
	// Run tests.
	tb := mesh.NewTestbed(behavior, test)
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
	// Test evaluation.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		if tbe.Done(evt) {
			tbe.SignalFail("signal fail after done")
		}
		return nil
	}
	// Run tests.
	tb := mesh.NewTestbed(behavior, test)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
	}, time.Second)
	assert.ErrorContains(err, "test failed: signal fail after done")
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
	// Test evaluation.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		if tbe.Done(evt) {
			tbe.SignalSuccess()
		}
		tbe.Push(evt)
		return nil
	}
	// Run tests.
	tb := mesh.NewTestbed(behavior, test)
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
	// Test evaluation.
	test := func(tbe *mesh.TestbedEvaluator, evt *mesh.Event) error {
		if tbe.Done(evt) {
			tbe.SignalSuccess()
		}
		switch evt.Topic() {
		case "fail":
			return errors.New("ouch")
		default:
			return nil
		}
	}
	// Run tests.
	tb := mesh.NewTestbed(behavior, test)
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
