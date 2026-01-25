// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh_test

import (
	"testing"
	"time"

	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/mesh"
)

// TestTestbedSuccess verifies the successful working of the testbed
// for behavior tests.
func TestTestbedSuccess(t *testing.T) {
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
		func(tbe *mesh.TestbedEvaluator) {
			tbe.AssertRetry(func() bool { return tbe.Len() == 3 }, "collected events not 3: %d", tbe.Len())
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
	}, time.Second)
	verify.NoError(t, err)
}

// TestTestbedFail verifies the failing working of the testbed
// for behavior tests.
func TestTestbedFail(t *testing.T) {
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
		func(tbe *mesh.TestbedEvaluator) {
			tbe.Assert(false, "must fail")
		},
	)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
	}, time.Second)
	verify.ErrorContains(t, err, "test failed: must fail")
}

// TestTestbedMesh verifies the Mesh stubbing of the testbed.
func TestTestbedMesh(t *testing.T) {
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
					verify.ErrorContains(t, err, "cell name 'anything' already used")
				case "subscribe":
					out.EmitEvent(evt)
					err := cell.Mesh().Subscribe("anything", "anything-else")
					verify.ErrorContains(t, err, "emitter cell 'anything' does not exist")
				case "unsubscribe":
					out.EmitEvent(evt)
					err := cell.Mesh().Unsubscribe("anything", "anything-else")
					verify.ErrorContains(t, err, "emitter cell 'anything' does not exist")
				case "emit":
					out.EmitEvent(evt)
					err := cell.Mesh().EmitEvent("anything", evt)
					verify.ErrorContains(t, err, "cell 'anything' does not exist")
				case "emitter":
					out.EmitEvent(evt)
					emtr, err := cell.Mesh().Emitter("anything")
					verify.ErrorContains(t, err, "cell 'anything' does not exist")
					if emtr != nil {
						t.Fatalf("expected nil emitter, got %v", emtr)
					}
				}
			}
		}
	}
	behavior := mesh.BehaviorFunc(mesher)
	// Run tests.
	tb := mesh.NewTestbed(
		behavior,
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
	verify.NoError(t, err)
}
